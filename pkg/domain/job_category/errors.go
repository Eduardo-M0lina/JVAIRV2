package job_category

import "errors"

var (
	ErrJobCategoryNotFound = errors.New("job category not found")
	ErrJobCategoryInUse    = errors.New("cannot delete job category that is still in use")
)
