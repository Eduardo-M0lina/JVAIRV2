package job_priority

import "errors"

var (
	ErrJobPriorityNotFound = errors.New("job priority not found")
	ErrJobPriorityInUse    = errors.New("cannot delete job priority that is still in use")
)
