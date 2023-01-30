package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chaos-io/chaos/core/resource"
)

func TestResource(t *testing.T) {
	assert.Equal(t, []byte("bar"), resource.Get("foo"))
	assert.Equal(t, []byte("baz"), resource.Get("bar"))
}
