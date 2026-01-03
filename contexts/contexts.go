package contexts

import (
	"context"

	"github.com/chaos-io/chaos/consts"
)

type ctxLocalKeyType struct{}

var ctxLocalKey = ctxLocalKeyType{}

func WithLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, ctxLocalKey, locale)
}

func CtxLocale(ctx context.Context) string {
	locale, ok := ctx.Value(ctxLocalKey).(string)
	if !ok || locale == "" {
		return consts.LocaleDefault
	}
	return locale
}
