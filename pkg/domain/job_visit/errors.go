package job_visit

import "errors"

var (
	ErrJobVisitNotFound = errors.New("job visit not found")
	ErrJobNotFound      = errors.New("job not found")
	ErrUserNotFound     = errors.New("user not found")
)
