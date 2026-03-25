package storage

import (
	"testing"
	"time"
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
	if len(*headers) != 3 {
		t.Errorf("headers length got %d, want %d", len(*headers), 3)
	}
}
