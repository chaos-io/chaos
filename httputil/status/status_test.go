package status_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chaos-io/chaos/httputil/status"
)

func TestGetCodeGroup(t *testing.T) {
	for i := range make([]struct{}, 1000) {
		group := status.GetCodeGroup(i)
		switch i / 100 {
		case 1:
			assert.Equal(t, status.Informational, group)
		case 2:
			assert.Equal(t, status.Successful, group)
		case 3:
			assert.Equal(t, status.Redirection, group)
		case 4:
			assert.Equal(t, status.ClientError, group)
		case 5:
			assert.Equal(t, status.ServerError, group)
		default:
			assert.Equal(t, status.Unknown, group)
		}
	}
}
