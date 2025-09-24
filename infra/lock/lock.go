package lock

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cenk/backoff"
	"github.com/chaos-io/chaos/infra/redis"
	"github.com/chaos-io/chaos/pkg/logs"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=mocks/lock.go -package=mocks . ILocker
type ILocker interface {
	WithHolder(holder string) ILocker
	Lock(ctx context.Context, key string, expiresIn time.Duration) (bool, error)
	Unlock(key string) (bool, error)
	ExpireLockIn(key string, expiresIn time.Duration) (bool, error)
	LockBackoff(ctx context.Context, key string, expiresIn, maxWait time.Duration) (bool, error)
	// LockBackoffWithRenew 获取锁并异步保持定时续期，每次锁保持时间为 ttl，到达 maxHold 时间或被 cancel
	// 后退出续期。调用方做写操作前应检查 ctx.Done 以确认仍持有锁，发生错误时应调用 cancel 以主动释放锁。
	LockBackoffWithRenew(parent context.Context, key string, ttl, maxHold time.Duration) (locked bool, ctx context.Context, cancel func(), err error)
	LockWithRenew(parent context.Context, key string, ttl, maxHold time.Duration) (locked bool, ctx context.Context, cancel func(), err error)
}

type redisLocker struct {
	c      redis.Cmdable
	holder string
}

func NewRedisLocker(c redis.Cmdable) ILocker {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown_hostname"
	}
	return &redisLocker{
		c:      c,
		holder: fmt.Sprintf("%s-%s", hostname, uuid.New().String()),
	}
}

func NewRedisLockerWithHolder(c redis.Cmdable, holder string) ILocker {
	return &redisLocker{
		c:      c,
		holder: holder,
	}
}

func (l *redisLocker) WithHolder(holder string) ILocker {
	l.holder = holder
	return l
}

func (l *redisLocker) Lock(ctx context.Context, key string, expiresIn time.Duration) (bool, error) {
	if expiresIn < time.Second {
		return false, fmt.Errorf("lock ttl is too short")
	}
	return l.c.SetNX(ctx, key, l.holder, expiresIn).Result()
}

func (l *redisLocker) Unlock(key string) (bool, error) {
	const script = `if redis.call('GET', KEYS[1]) == ARGV[1] then redis.call('DEL', KEYS[1]); return 1; end; return 0;`
	result, err := l.c.Eval(context.Background(), script, []string{key}, l.holder).Result()
	if err != nil {
		return false, fmt.Errorf("unlock with lua script, error: %v", err)
	}

	rt, ok := result.(int64)
	if !ok {
		return false, fmt.Errorf("unknown result type %T", result)
	}

	return rt == 1, nil
}

func (l *redisLocker) ExpireLockIn(key string, expiresIn time.Duration) (bool, error) {
	const script = `if redis.call('GET', KEYS[1]) == ARGV[1] then redis.call('PEXPIRE', KEYS[1], ARGV[2]); return 1; end; return 0;`
	result, err := l.c.Eval(context.Background(), script, []string{key}, l.holder, int64(expiresIn/time.Millisecond)).Result()
	if err != nil {
		return false, fmt.Errorf("extend lock error: %v", err)
	}

	rt, ok := result.(int64)
	if !ok {
		return false, fmt.Errorf("unknown result type")
	}

	return rt == 1, nil
}

func (l *redisLocker) LockBackoff(ctx context.Context, key string, expiresIn, maxWait time.Duration) (bool, error) {
	var locked bool
	bf := backoff.NewExponentialBackOff()
	bf.InitialInterval = 50 * time.Millisecond
	bf.MaxInterval = 300 * time.Millisecond
	bf.MaxElapsedTime = maxWait

	errNotLocked := errors.New("lock hold by other locker")
	err := backoff.Retry(func() error {
		var err error
		locked, err = l.Lock(ctx, key, expiresIn)
		if err != nil {
			return err
		}
		if !locked {
			return errNotLocked
		}
		return nil
	}, bf)
	if err != nil {
		if errors.Is(err, errNotLocked) {
			return false, nil
		}
		return false, err
	}

	return locked, nil
}

func (l *redisLocker) LockWithRenew(parent context.Context, key string, ttl, maxHold time.Duration) (locked bool, ctx context.Context, cancel func(), err error) {
	nop := func() {}
	locked, err = l.Lock(parent, key, ttl)
	if err != nil || !locked {
		return locked, parent, nop, err
	}

	ctx, cancel = context.WithCancel(parent)
	go func() {
		defer cancel()
		l.renewLock(ctx, key, ttl, maxHold)
	}()

	return locked, ctx, cancel, nil
}

func (l *redisLocker) LockBackoffWithRenew(parent context.Context, key string, ttl, maxWait time.Duration) (locked bool, ctx context.Context, cancel func(), err error) {
	nop := func() {}
	locked, err = l.LockBackoff(parent, key, ttl, ttl+time.Second)
	if err != nil || !locked {
		return locked, parent, nop, err
	}

	ctx, cancel = context.WithCancel(parent)
	go func() {
		defer cancel()
		l.renewLock(ctx, key, ttl, maxWait)
	}()

	return locked, ctx, cancel, nil
}

func (l *redisLocker) renewLock(ctx context.Context, key string, ttl, maxHold time.Duration) {
	t1 := time.After(maxHold)
	t2 := time.NewTicker(max(time.Second, ttl-100*time.Millisecond))

	retry := 0
	unlock := func() {
		if _, err := l.Unlock(key); err != nil {
			logs.Warnw("failed to renew defer unlock", "key", key, "error", err)
		}
	}

	defer t2.Stop()

	for {
		select {
		case <-t1:
			logs.Infow("renew lock reached max hold duration", "key", key)
			unlock()
			return
		case <-t2.C:
			locked, err := l.ExpireLockIn(key, ttl)
			switch {
			case err != nil:
				if retry++; retry >= 3 {
					logs.Errorw("renew lock got too many retries, no more retry", "key", key, "error", err)
					return
				}
				logs.Warnw("renew lock got error, will retry", "key", key, "error", err)
			case !locked:
				logs.Infow("renew lock got non-ok, exiting", "key", key)
				return
			case locked:
				retry = 0
			}
		case <-ctx.Done():
			logs.Infow("renew lock context done", "key", key)
			unlock()
			return
		}
	}
}
