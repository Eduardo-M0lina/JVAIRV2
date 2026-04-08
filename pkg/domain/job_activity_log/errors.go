package job_activity_log

import "errors"

var (
	ErrNotFound     = errors.New("job activity log not found")
	ErrJobNotFound  = errors.New("job not found")
	ErrUserNotFound = errors.New("user not found")
)
