package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Do(t *testing.T) {
	do, err := Do(ctx, "PING")
	assert.NoError(t, err)
	assert.Equal(t, "PONG", do)
}
