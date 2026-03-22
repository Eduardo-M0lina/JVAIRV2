package job_sms

import "errors"

var (
	ErrNotFound    = errors.New("job sms not found")
	ErrJobNotFound = errors.New("job not found")
)
