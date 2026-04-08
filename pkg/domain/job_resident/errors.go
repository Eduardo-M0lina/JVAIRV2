package job_resident

import "errors"

var (
	ErrNotFound    = errors.New("job resident not found")
	ErrJobNotFound = errors.New("job not found")
)
