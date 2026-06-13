# errorx

`errorx` 用于在 Go 服务中统一定义、生成和使用业务错误码。

核心流程：

1. 用 YAML 定义业务错误码。
2. 用 `errorxgen` 生成 `errcode` Go 包。
3. 服务启动时注册错误码。
4. 业务代码通过生成的错误定义创建、包装、判断错误。

## 1. 定义错误码

建议把 YAML 放在业务服务内，例如：

```text
configs/error_code/task.yaml
```

示例：

```yaml
appCode: 6
bizCode: 12
errorCode:
  - name: TaskNotFound
    code: 1001
    message: task {task_id} not found
    description: task does not exist
    affectsStability: false

  - name: TaskStateInvalid
    code: 1002
    message: task {task_id} state invalid
    description: task state cannot satisfy current operation
```

字段说明：

- `appCode`：应用编码，范围 `1..9`。
- `bizCode`：业务模块编码，范围 `1..9999`。
- `errorCode`：错误码列表，不能为空。
- `name`：生成后的 Go 变量名，必须是导出的 Go 标识符，例如 `TaskNotFound`。
- `code`：当前业务模块内的错误子码，范围 `1..9999`。
- `message`：错误消息，支持 `{key}` 占位符。
- `description`：生成代码里的注释，可选。
- `affectsStability`：是否影响稳定性指标，可选，默认 `true`。

完整错误码规则：

```text
appCode(1位) + bizCode(4位) + code(4位)
```

例如：

```text
appCode=6, bizCode=12, code=1001 => 600121001
```

## 2. 生成 errcode 代码

在业务服务仓库中执行：

```bash
go run github.com/chaos-io/chaos/errorx/cmd/errorxgen \
  -out ./internal/errcode \
  -pkg errcode \
  ./configs/error_code
```

参数说明：

- `-out`：生成代码目录。
- `-pkg`：生成代码的包名，通常使用 `errcode`。
- `-errorx-import`：`errorx` 的 import path，默认 `github.com/chaos-io/chaos/errorx`。
- 最后一个参数：YAML 文件或目录。传目录时会递归读取 `.yaml` 和 `.yml` 文件。

也可以在业务服务中加入 `go:generate`：

```go
//go:generate go run github.com/chaos-io/chaos/errorx/cmd/errorxgen -out ./internal/errcode -pkg errcode ./configs/error_code
```

然后执行：

```bash
go generate ./...
```

生成后的代码大致如下：

```go
var TaskNotFound = errorx.Define(
    600121001,
    "task {task_id} not found",
    errorx.AffectsStability(false),
)

func RegisterAll() error {
    return errorx.Register(
        TaskNotFound,
        TaskStateInvalid,
    )
}
```

## 3. 服务启动时注册

```go
package service

import "your/service/internal/errcode"

func Init() error {
    return errcode.RegisterAll()
}
```

如果同一个错误码重复注册且定义不一致，`RegisterAll` 会返回错误；服务应直接启动失败。

## 4. 业务代码中使用

创建业务错误：

```go
return errcode.TaskNotFound.New(
    errorx.WithMessageParam("task_id", taskID),
    errorx.WithExtra(map[string]string{"task_id": taskID}),
)
```

包装底层错误：

```go
if err != nil {
    return errcode.TaskNotFound.Wrap(err, errorx.WithExtra(map[string]string{
        "task_id": taskID,
    }))
}
```

判断错误类型：

```go
if errcode.TaskNotFound.Is(err) {
    // handle task not found
}
```

在网关、日志、中间件等通用边界读取错误码：

```go
if coded, ok := errorx.From(err); ok {
    code := coded.Code()
    message := coded.Message()
    extra := coded.Extra()
    _, _, _ = code, message, extra
}
```

## 约束

- 同一批 YAML 中不能有重复的错误名。
- 同一批 YAML 中不能有重复的完整错误码。
- YAML 文件名生成 Go 文件名，文件名冲突会生成失败。
- `errorx` 只处理业务错误码，不绑定 HTTP/gRPC 状态码；协议映射应放在业务服务的 adapter 层。
