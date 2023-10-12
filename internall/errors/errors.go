package errors

import "errors"

var (
	ErrNotFound       = errors.New("your requested item is not found")
	ErrConflict       = errors.New("your email already exist, must be unique")
	ErrBadParamInput  = errors.New("given param is not valid")
	ErrInternalServer = errors.New("internal Server Error")
)
