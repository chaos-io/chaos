package redis

// import (
// 	"testing"
//
// 	"github.com/alicebob/miniredis/v2"
// 	"github.com/pkg/errors"
// 	"github.com/redis/go-redis/v9"
// )
//
// type tester interface {
// 	Errorf(format string, args ...interface{})
// 	Cleanup(func())
// }
//
// func NewTestRedis(t *testing.T) Cmdable {
// 	cli, err := newMiniRedis(t)
// 	if err != nil {
// 		t.Errorf("new test redis failed, err=%v", err)
// 		return nil
// 	}
// 	return cli
// }
//
// func newMiniRedis(t tester) (Cmdable, error) {
// 	m := miniredis.NewMiniRedis()
// 	if err := m.Start(); err != nil {
// 		return nil, errors.WithMessage(err, "start miniredis")
// 	}
//
// 	opts := &redis.Options{Addr: m.Addr()}
// 	p := NewClient(opts)
//
// 	t.Cleanup(m.Close)
// 	return p, nil
// }
