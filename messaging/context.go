package messaging

import "context"

type contextKey string

const (
	topicKey             contextKey = "@messaging/topic"
	messageKey           contextKey = "@messaging/message"
	messageIDKey         contextKey = "@messaging/message.id"
	messageAttributesKey contextKey = "@messaging/message.attributes"
)

func WithTopic(ctx context.Context, topic string) context.Context {
	return context.WithValue(ensureContext(ctx), topicKey, topic)
}

func WithMessage(ctx context.Context, m *SubMessage) context.Context {
	ctx = context.WithValue(ensureContext(ctx), messageKey, m)
	if m == nil {
		return ctx
	}
	ctx = context.WithValue(ctx, messageIDKey, m.Id)
	return context.WithValue(ctx, messageAttributesKey, m.Attributes)
}

func WithMessageID(ctx context.Context, id string) context.Context {
	return context.WithValue(ensureContext(ctx), messageIDKey, id)
}

func WithMessageAttributes(ctx context.Context, attributes map[string]any) context.Context {
	return context.WithValue(ensureContext(ctx), messageAttributesKey, attributes)
}

func GetContextTopic(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if value, ok := ctx.Value(topicKey).(string); ok {
		return value
	}
	return ""
}

func GetContextMessage(ctx context.Context) *SubMessage {
	if ctx == nil {
		return nil
	}
	if value, ok := ctx.Value(messageKey).(*SubMessage); ok {
		return value
	}
	return nil
}

func GetContextMessageId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if value, ok := ctx.Value(messageIDKey).(string); ok {
		return value
	}
	return ""
}

func GetContextMessageAttributes(ctx context.Context) map[string]any {
	if ctx == nil {
		return nil
	}
	if value, ok := ctx.Value(messageAttributesKey).(map[string]any); ok {
		return value
	}
	return nil
}

func ensureContext(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}
	return context.Background()
}
