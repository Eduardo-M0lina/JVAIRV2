package job_task

import "errors"

var (
	ErrNotFound           = errors.New("job task not found")
	ErrJobNotFound        = errors.New("job not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrTaskStatusNotFound = errors.New("task status not found")
)
