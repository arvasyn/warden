package apperr

import "errors"

var (
	ErrInvalidPath    = errors.New("invalid file path provided")
	ErrNotImplemented = errors.New("not implemented")
)
