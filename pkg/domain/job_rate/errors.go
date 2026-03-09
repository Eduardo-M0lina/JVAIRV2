package job_rate

import "errors"

var (
	ErrNotFound              = errors.New("job rate not found")
	ErrJobNotFound           = errors.New("job not found")
	ErrUserNotFound          = errors.New("user not found")
	ErrJobRateStatusNotFound = errors.New("job rate status not found")
)
