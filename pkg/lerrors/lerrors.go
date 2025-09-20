package lerrors

import "errors"

var (
	ErrPortConflict = errors.New("port conflict")
	ErrInvalidUsage = errors.New("invalid usage for command line")
)
