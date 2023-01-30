package test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chaos-io/chaos/core/resource"
)

func TestResource(t *testing.T) {
	assert.Equal(t, []byte("hello world"), resource.Get("/a.txt"))

	bindata, err := ioutil.ReadFile("testdata/b.bin")
	assert.NoError(t, err)
	assert.Equal(t, bindata, resource.Get("/b.bin"))

	assert.Equal(t, []byte("handle this"), resource.Get("testdata/collision.txt"))
}
