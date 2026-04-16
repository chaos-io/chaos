package messaging

import "context"

type contextKey string

const (
	TopicKey             contextKey = "@messaging/topic"
	MessageKey           contextKey = "@messaging/message"
	MessageIdKey         contextKey = "@messaging/message.id"
	MessageAttributesKey contextKey = "@messaging/message.attributes"
)

func WithTopic(ctx context.Context, topic string) context.Context {
	return context.WithValue(ensureContext(ctx), TopicKey, topic)
}

func WithMessage(ctx context.Context, m *SubMessage) context.Context {
	ctx = context.WithValue(ensureContext(ctx), MessageKey, m)
	if m == nil {
		return ctx
	}
	ctx = context.WithValue(ctx, MessageIdKey, m.Id)
	return context.WithValue(ctx, MessageAttributesKey, m.Attributes)
}

func WithMessageID(ctx context.Context, id string) context.Context {
	return context.WithValue(ensureContext(ctx), MessageIdKey, id)
}

func WithMessageAttributes(ctx context.Context, attributes map[string]any) context.Context {
	return context.WithValue(ensureContext(ctx), MessageAttributesKey, attributes)
}

func GetContextTopic(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if value, ok := ctx.Value(TopicKey).(string); ok {
		return value
	}
	return ""
}

func GetContextMessage(ctx context.Context) *SubMessage {
	if ctx == nil {
		return nil
	}
	if value, ok := ctx.Value(MessageKey).(*SubMessage); ok {
		return value
	}
	return nil
}

func GetContextMessageId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if value, ok := ctx.Value(MessageIdKey).(string); ok {
		return value
	}
	return ""
}

func GetContextMessageAttributes(ctx context.Context) map[string]any {
	if ctx == nil {
		return nil
	}
	if value, ok := ctx.Value(MessageAttributesKey).(map[string]any); ok {
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
