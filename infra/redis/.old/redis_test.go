//go:build local
// +build local

package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	client := New(nil)
	do, err := client.Client.Do(ctx, "PING").Result()
	assert.NoError(t, err)
	assert.Equal(t, "PONG", do.(string))
}
