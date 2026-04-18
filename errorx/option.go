package errorx

import (
	"fmt"
	"strings"
)

type Option func(err *statusError)

func WithExtraMsg(extraMsg string) Option {
	return func(err *statusError) {
		if err == nil || err.status == nil || extraMsg == "" {
			return
		}
		err.status.message = fmt.Sprintf("%s,%s", err.status.message, extraMsg)
	}
}

func WithMsgParam(k, v string) Option {
	return func(err *statusError) {
		if err == nil || err.status == nil {
			return
		}
		err.status.message = strings.ReplaceAll(err.status.message, fmt.Sprintf("{%s}", k), v)
	}
}

func WithExtra(extra map[string]string) Option {
	return func(err *statusError) {
		if err == nil || err.status == nil || extra == nil {
			return
		}
		err.status.extra = extra
	}
}
