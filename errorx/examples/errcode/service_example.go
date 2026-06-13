package errcodeexample

import (
	"github.com/chaos-io/chaos/errorx"
	"github.com/chaos-io/chaos/errorx/examples/errcode/generated"
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
