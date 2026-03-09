package alert

import "errors"

var (
	ErrNotFound     = errors.New("alert not found")
	ErrUserNotFound = errors.New("user not found")
)
