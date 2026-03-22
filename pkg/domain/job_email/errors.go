package job_email

import "errors"

var (
	ErrNotFound    = errors.New("job email not found")
	ErrJobNotFound = errors.New("job not found")
)
