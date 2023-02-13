package recovery

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chaos-io/chaos/core/logs"
)

func TestWithLogger(t *testing.T) {
	var mw middleware

	logger := logs.New(&logs.Config{})
	opt := WithLogger(logger)
	opt(&mw)

	assert.Same(t, logger, mw.l)
}

func TestWithCallBack(t *testing.T) {
	var mw middleware

	callback := func(_ http.ResponseWriter, _ *http.Request, _ error) {}

	opt := WithCallBack(callback)
	opt(&mw)

	assert.Equal(t, fmt.Sprintf("%p", callback), fmt.Sprintf("%p", mw.panicCallback))
}
