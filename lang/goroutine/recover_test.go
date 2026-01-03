package goroutine

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	t.Run("recover from panic", func(t *testing.T) {
		ctx := context.Background()
		var err error

		func() {
			defer Recover(ctx, &err)
			panic("test panic")
		}()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "panic: ")
	})

	t.Run("recover with existing error", func(t *testing.T) {
		ctx := context.Background()
		originalErr := errors.New("original error")
		err := originalErr

		func() {
			defer Recover(ctx, &err)
			panic("test panic")
		}()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "panic: ")
		assert.Contains(t, err.Error(), originalErr.Error())
	})

	t.Run("recover with nil context", func(t *testing.T) {
		var err error

		func() {
			defer Recover(nil, &err)
			panic("test panic")
		}()

		assert.Error(t, err)
	})
}

func TestRecovery(t *testing.T) {
	t.Run("recover from panic", func(t *testing.T) {
		ctx := context.Background()

		func() {
			defer Recovery(ctx)
			panic("test panic")
		}()
		// No assertion needed as Recovery only logs the error
	})
}
