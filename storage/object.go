package storage

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Object struct {
	Etag         string    `json:"etag,omitempty"`
	Key          string    `json:"key,omitempty"`
	LastModified time.Time `json:"lastModified,omitempty"`
	Size         int64     `json:"size,omitempty"`
	ContentType  string    `json:"contentType,omitempty"`
	Content      []byte    `json:"content,omitempty"`
}

func (x *Object) GetHttpHeaders() *http.Header {
	if x != nil {
		headers := &http.Header{}
		if len(x.ContentType) > 0 {
			headers.Set("Content-Type", x.ContentType)
			if strings.Contains(x.ContentType, "zip") {
				headers.Set("Content-Encoding", "gzip")
			}
		}

		if len(x.Etag) > 0 {
			headers.Set("ETag", x.Etag)
		}

		return headers
	}

	return nil
}

func (x *Object) WriteHttpResponse(ctx context.Context, writer http.ResponseWriter) error {
	_ = ctx
	if x != nil {
		headers := x.GetHttpHeaders()
		for key, values := range *headers {
			if len(values) > 0 {
				for _, value := range values {
					writer.Header().Set(key, value)
				}
			}
		}

		if count, err := writer.Write(x.Content); err != nil {
			return err
		} else if count < len(x.Content) {
			return fmt.Errorf("failed to write completed content, expect :%d, actual: %d", len(x.Content), count)
		}

		return nil
	}

	return nil
}
