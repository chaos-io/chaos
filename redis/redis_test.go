package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testNew(t *testing.T) {
	client := New(nil)
	do, err := client.Client.Do(ctx, "PING").Result()
	assert.NoError(t, err)
	assert.Equal(t, "PONG", do.(string))
}
