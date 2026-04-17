package errorx

import (
	"fmt"
	"strings"
)

type Option func(ws *withStatus)

func WithExtraMsg(extraMsg string) Option {
	return func(ws *withStatus) {
		if ws == nil || ws.status == nil || extraMsg == "" {
			return
		}
		ws.status.message = fmt.Sprintf("%s,%s", ws.status.message, extraMsg)
	}
}

func WithMsgParam(k, v string) Option {
	return func(ws *withStatus) {
		if ws == nil || ws.status == nil {
			return
		}
		ws.status.message = strings.ReplaceAll(ws.status.message, fmt.Sprintf("{%s}", k), v)
	}
}

func WithExtra(extra map[string]string) Option {
	return func(ws *withStatus) {
		if ws == nil || ws.status == nil || extra == nil {
			return
		}
		ws.status.extra = extra
	}
}
