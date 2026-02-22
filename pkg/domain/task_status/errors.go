package task_status

import "errors"

var (
	ErrTaskStatusNotFound = errors.New("task status not found")
	ErrTaskStatusInUse    = errors.New("cannot delete task status that is still in use")
)
