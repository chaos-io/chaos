# Messaging

`messaging` 提供一层很薄的消息队列抽象，统一发布、订阅、消息上下文和确认语义。当前目录定义了通用模型和 `Client`，具体实现放在子包里，例如 [`nats`](./nats)。

## 适用场景

- 业务代码只依赖统一接口，不直接耦合到底层消息队列 SDK。
- 发布和订阅都通过同一套 `Message` / `Subscription` 模型表达。
- 订阅处理函数可以从 `context.Context` 中拿到 topic、message id 和 attributes。

## 核心类型

### `Message`

发布消息时使用的基础结构：

```go
type Message struct {
    Id         string
    TraceId    string
    SpanId     string
    Attributes map[string]any
    Data       string
}
```

其中 `Data` 是消息体字符串，额外元数据可以放在 `Attributes`。

### `Subscription`

订阅配置的公共模型：

```go
type Subscription struct {
    Name    string
    Topic   string
    Group   string
    AutoAck bool
}
```

- `Topic` 不能为空，`Client.Subscribe` 和各具体实现都会复用 `Validate()` 做校验。
- `Group` 非空时，底层实现可以按队列订阅方式消费。
- `AutoAck` 为 `true` 时，handler 成功返回后会自动确认消息。

### `SubMessage`

消费时传给 handler 的消息对象，继承 `Message` 字段，并带确认能力：

- `Ack()`：确认消息。
- `Nak()`：拒绝并允许重投。
- `Term()`：终止消息。
- `InProgress()`：告知底层消息仍在处理中。
- `Done()`：是否已经执行过一次终结动作。

`Ack` / `Nak` / `Term` 只会生效一次，避免重复确认。

## 快速开始

### 1. 通过配置初始化

推荐直接使用根包配置：

```go
package main

import (
    "log"

    "github.com/chaos-io/chaos/messaging"
    "github.com/chaos-io/chaos/messaging/nats"
)

func main() {
    nats.Register()

    client, err := messaging.New()
    if err != nil {
        log.Fatal(err)
    }
    defer client.Shutdown()
}
```

默认会读取 `messaging` 配置项，对应 [messaging/config/messaging.yaml](./config/messaging.yaml)：

```yaml
messaging:
  driver: nats
  nats:
    url: nats://127.0.0.1:4222
    jetStream: true
  subscriptions:
    - name: agent-consumer
      topic: demo.start-task
      group: start-task
      pull: true
      ackWait: 5m
      service: Agent
      method: start_task
```

### 2. 手动初始化底层队列

如果你不想依赖外部配置，也可以直接构造具体队列：

```go
package main

import (
    "log"

    "github.com/chaos-io/chaos/messaging"
    "github.com/chaos-io/chaos/messaging/nats"
)

func main() {
    queue, err := nats.NewWithConfig(&messaging.NatsConfig{
        URL:       "nats://127.0.0.1:4222",
        JetStream: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer queue.Shutdown()

    client, err := messaging.NewClient(queue)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Shutdown()
}
```

### 3. 发布消息

```go
err := client.Publish(ctx, "demo.user.created", &messaging.Message{
    Id:      "msg-1",
    TraceId: "trace-1",
    Attributes: map[string]any{
        "source": "api",
    },
    Data: `{"userId":"u-1"}`,
})
```

发布前会校验 topic 非空。

### 4. 订阅消息

```go
err := client.Subscribe(&messaging.Subscription{
    Name:    "user-created-worker",
    Topic:   "demo.user.created",
    Group:   "user-workers",
    AutoAck: true,
}, func(ctx context.Context, sub *messaging.Subscription, msg *messaging.SubMessage) error {
    topic := messaging.GetContextTopic(ctx)
    messageID := messaging.GetContextMessageId(ctx)
    attrs := messaging.GetContextMessageAttributes(ctx)

    _ = topic
    _ = messageID
    _ = attrs

    // 处理 msg.Data
    return nil
})
```

默认行为：

- handler 返回 `nil` 且 `AutoAck=true` 时，框架自动 `Ack()`。
- handler 返回错误且消息还没结束时，具体实现可以执行 `Nak()`；当前 `nats` 实现就是这个行为。
- 如果你在 handler 内主动调用了 `Ack()` / `Nak()` / `Term()`，后续不会重复执行。

## 手动确认示例

当你不想依赖 `AutoAck` 时，可以显式控制确认：

```go
err := client.Subscribe(&messaging.Subscription{
    Name:    "manual-worker",
    Topic:   "demo.task",
    AutoAck: false,
}, func(ctx context.Context, sub *messaging.Subscription, msg *messaging.SubMessage) error {
    if msg.Data == "" {
        msg.Term()
        return nil
    }

    if err := handle(msg.Data); err != nil {
        msg.Nak()
        return err
    }

    msg.Ack()
    return nil
})
```

## 上下文辅助函数

订阅回调里可以使用这些 helper：

- `GetContextTopic(ctx)`：获取当前 topic。
- `GetContextMessage(ctx)`：获取完整 `*SubMessage`。
- `GetContextMessageId(ctx)`：获取消息 ID。
- `GetContextMessageAttributes(ctx)`：获取消息 attributes。

如果你需要把消息信息继续透传到下游，也可以使用：

- `WithTopic(ctx, topic)`
- `WithMessage(ctx, msg)`
- `WithMessageID(ctx, id)`
- `WithMessageAttributes(ctx, attrs)`

## 错误约定

常见输入校验错误：

- `ErrNilClient`
- `ErrNilQueue`
- `ErrNilSubscription`
- `ErrNilHandler`
- `ErrEmptyTopic`

这些错误用于尽早暴露调用方传参问题。

## NATS 扩展

[`nats`](./nats) 子包在 `Subscription` 之外额外提供了 `Consumer`：

```go
type Consumer struct {
    Subscription      messaging.Subscription
    Pull              bool
    AckWait           time.Duration
    PullMaxWaiting    int
    PendingMsgLimit   int
    PendingBytesLimit int
}
```

适合需要 pull consumer、ack wait 或 pending limit 等 NATS 特有能力的场景。`Consumer.Validate()` 会复用 `Subscription.Validate()`。
