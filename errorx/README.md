# errorx

`errorx` provides a small production-oriented business error model for Go services:

1. Define stable business error codes in YAML.
2. Generate a small `errcode` package.
3. Register definitions at service startup.
4. Create, wrap, and match business errors through generated objects.

## Runtime API

Define an error manually:

```go
var TaskNotFound = errorx.Define(
    600121001,
    "task {task_id} not found",
    errorx.AffectsStability(false),
)
```

Create and wrap errors:

```go
err := TaskNotFound.New(
    errorx.WithMessageParam("task_id", "t-1"),
    errorx.WithExtra(map[string]string{"task_id": "t-1"}),
)

err = TaskNotFound.Wrap(dbErr)
```

Match and read errors:

```go
if TaskNotFound.Is(err) {
    // handle task-not-found
}

if coded, ok := errorx.From(err); ok {
    _ = coded.Code()
    _ = coded.Message()
    _ = coded.Extra()
}
```

Register definitions during service startup:

```go
func Init() error {
    return errcode.RegisterAll()
}
```

`Register` is idempotent for identical definitions and returns `ErrRegisterConflict` when the same code is registered with different metadata.

## YAML Format

```yaml
appCode: 6
bizCode: 12
errorCode:
  - name: TaskNotFound
    code: 1001
    message: task {task_id} not found
    description: task does not exist
    affectsStability: false
```

Fields:

- `appCode`: required, `1..9`.
- `bizCode`: required, `1..9999`.
- `errorCode`: required, non-empty list.
- `name`: required, exported Go identifier. It becomes the generated variable name.
- `code`: required, `1..9999`.
- `message`: optional business message. Empty values use `errorx.DefaultMessage`.
- `description`: optional generated Go comment.
- `affectsStability`: optional, defaults to `true`.

The full code layout is:

```text
appCode(1 digit) + bizCode(4 digits) + code(4 digits)
```

Example:

```text
appCode=6, bizCode=12, code=1001 -> 600121001
```

Use stable business names such as `TaskNotFound`, `UserAlreadyExists`, and `PaymentProviderUnavailable`.
Avoid transport-specific names such as `HTTP404` or `GRPCInvalidArgument`; protocol mapping belongs in each service adapter.

Generation fails when one input batch contains duplicate generated file names, duplicate error names, duplicate full error codes, invalid code ranges, or invalid Go identifiers.

## Generator

Run the generator:

```bash
go run ./errorx/cmd/errorxgen \
  -out ./internal/errcode \
  -pkg errcode \
  -errorx-import github.com/chaos-io/chaos/errorx \
  ./configs/error_code
```

Common `go:generate` usage inside a service:

```go
//go:generate go run github.com/chaos-io/chaos/errorx/cmd/errorxgen -out ./internal/errcode -pkg errcode ./configs/error_code
```

Generated code exposes one `errorx.Definition` per YAML item:

```go
var TaskNotFound = errorx.Define(
    600121001,
    "task {task_id} not found",
    errorx.AffectsStability(false),
)
```

It also generates:

```go
func RegisterAll() error
```

## Recommended Service Usage

```go
package service

import (
    "github.com/chaos-io/chaos/errorx"
    "your/service/internal/errcode"
)

func Init() error {
    return errcode.RegisterAll()
}

func FindTask(taskID string) error {
    return errcode.TaskNotFound.New(
        errorx.WithMessageParam("task_id", taskID),
        errorx.WithExtra(map[string]string{"task_id": taskID}),
    )
}

func LoadTask(taskID string) error {
    task, err := loadTask(taskID)
    if err != nil {
        return errcode.TaskNotFound.Wrap(err, errorx.WithExtra(map[string]string{
            "task_id": taskID,
        }))
    }
    _ = task
    return nil
}
```

Use `errorx.CodeOf(err)` or `errorx.Is(err, code)` only at generic boundaries. Business code should prefer generated definitions.
