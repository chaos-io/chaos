# errorx

`errorx` 用于两件事：

1. 在运行时按错误码创建、包装、判断业务错误
2. 从 YAML 生成业务错误码代码

## 快速开始

### 1. 服务启动时注册错误码

如果你使用生成代码：

```go
if err := errcode.RegisterAll(); err != nil {
    return err
}
```

如果你手工注册：

```go
if err := errorx.Register(
    600121001,
    "task not found",
    errorx.WithAffectsStability(true),
); err != nil {
    return err
}
```

### 2. 业务里创建错误

直接按错误码创建：

```go
err := errorx.NewByCode(
    600121001,
    errorx.WithExtra(map[string]string{"task_id": "123"}),
)
```

包装已有错误：

```go
err := errorx.WrapByCode(dbErr, 600121001)
```

生成普通错误：

```go
err := errorx.New("load config failed: %v", cause)
err = errorx.Wrapf(err, "init service")
```

### 3. 判断错误类型

```go
status, ok := errorx.FromStatus(err)
if ok && status.Code() == 600121001 {
    // handle
}
```

如果你使用生成代码，优先用生成函数：

```go
err := errcode.NewTaskNotFound()

if errcode.IsTaskNotFound(err) {
    // handle
}
```

## 运行时 API

最常用的入口：

- `Register`
- `MustRegister`
- `NewByCode`
- `WrapByCode`
- `FromStatus`
- `New`
- `Wrapf`

常用扩展项：

- `WithAffectsStability`
- `WithExtraMsg`
- `WithMsgParam`
- `WithExtra`

行为说明：

- 相同 `code` 且定义完全一致：`Register` 返回 `nil`
- 相同 `code` 但定义不一致：返回 `ErrRegisterConflict`
- `MustRegister` 会在冲突时直接 `panic`

## 生成器使用方式

生成器目录：

- `./gen_error_code`

### YAML 格式

使用小驼峰字段：

```yaml
appCode: 6
bizCode: 12
errorCode:
  - name: TaskNotFound
    code: 1001
    message: task not found
    description: task does not exist
    affectsStability: true
```

约束：

- `appCode` 范围：`1..9`
- `bizCode` 范围：`1..9999`
- `errorCode[*].code` 范围：`1..9999`
- `name` 必须是合法 Go 标识符

完整错误码格式：

```text
appCode(1位) + bizCode(4位) + subCode(4位)
```

例如：

```text
appCode=6, bizCode=12, code=1001 -> 600121001
```

### 生成命令

```bash
go run ./gen_error_code -out ./internal/errcode ./configs/error_code
```

或：

```bash
go generate ./gen_error_code
```

### 生成结果

每个 YAML 会生成一个业务文件，另外统一生成一个注册入口：

- `loop_task.go`
- `user.go`
- `register_all.go`

生成包名固定为 `errcode`。

每个业务文件通常包含：

- `TaskNotFoundCode`
- `NewTaskNotFound(opts ...errorx.Option) error`
- `IsTaskNotFound(err error) bool`

统一入口文件包含：

- `RegisterAll() error`

## 推荐接入方式

### 启动阶段

1. 调用 `errcode.RegisterAll()`
2. 如果返回错误，直接终止启动

### 业务阶段

1. 优先使用生成的 `NewXxx`
2. 需要保留底层原因时使用 `WrapByCode`
3. 判断错误类型时优先使用生成的 `IsXxx`
4. 跨层取错误码时使用 `FromStatus`

## 最小示例

假设生成目录是 `./internal/errcode`：

```go
package service

import (
    "chaos-io/chaos/errorx"
    "your/service/internal/errcode"
)

func Init() error {
    return errcode.RegisterAll()
}

func FindTask() error {
    return errcode.NewTaskNotFound(
        errorx.WithExtra(map[string]string{"task_id": "123"}),
    )
}
```
