package messaging

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

type Handler func(ctx context.Context, subscription *Subscription, m *SubMessage) error

func JSONEndpoint[T any](endpoint func(context.Context, any) (any, error)) Handler {
	if endpoint == nil {
		return nil
	}

	return func(ctx context.Context, subscription *Subscription, message *SubMessage) error {
		request := new(T)
		if err := jsoniter.ConfigFastest.UnmarshalFromString(message.Data, request); err != nil {
			message.Term()
			return fmt.Errorf(
				"messaging: decode topic %q for method %q: %w",
				subscription.Topic,
				subscription.Endpoint.Method,
				err,
			)
		}

		_, err := endpoint(ctx, request)
		return err
	}
}

//go:generate mockgen -destination=mocks/subscriber.go -package=mocks . Subscriber
type Subscriber interface {
	Subscribe(subscription *Subscription, handler Handler) error
}
