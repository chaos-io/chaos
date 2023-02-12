package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestObject_GetHttpHeaders(t *testing.T) {
	obj := &Object{
		Etag:         "1234",
		Key:          "",
		LastModified: time.Time{},
		Size:         0,
		ContentType:  "application/zip",
	}

	headers := obj.GetHttpHeaders()
	assert.Equal(t, 3, len(*headers))
}
