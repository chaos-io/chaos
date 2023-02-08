package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObject_GetHttpHeaders(t *testing.T) {
	obj := &Object{
		Etag:         "1234",
		Key:          "",
		LastModified: nil,
		Size:         0,
		ContentType:  "application/zip",
	}

	headers := obj.GetHttpHeaders()
	fmt.Printf("==headers=%+v\n", headers)
	assert.Equal(t, 3, len(*headers))
}
