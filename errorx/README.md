# errorx

`errorx` provides a production-oriented business error model for Go services:

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

The full code layout is:

```text
appCode(1 digit) + bizCode(4 digits) + code(4 digits)
```

Example:

```text
appCode=6, bizCode=12, code=1001 -> 600121001
```

## Generator

Run the generator from this module:

```bash
go run ./errorx/gen_error_code \
  -out ./internal/errcode \
  -pkg errcode \
  -errorx-import github.com/chaos-io/chaos/errorx \
  ./configs/error_code
```

Common `go:generate` usage inside a service:

```go
//go:generate go run github.com/chaos-io/chaos/errorx/gen_error_code -out ./internal/errcode ./configs/error_code
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

## Migration From Old API

Replace code-based calls with generated definitions:

```go
// Old
err := errorx.NewByCode(600121001)
err = errorx.WrapByCode(cause, 600121001)
status, ok := errorx.FromStatus(err)

// New
err := errcode.TaskNotFound.New()
err = errcode.TaskNotFound.Wrap(cause)
coded, ok := errorx.From(err)
```

Use `errorx.CodeOf(err)` or `errorx.Is(err, code)` only at generic boundaries. Business code should prefer generated definitions.
