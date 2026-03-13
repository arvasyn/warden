package apperr

import "errors"

// General Errors
var (
	ErrInvalidPath    = errors.New("invalid file path provided")
	ErrNotImplemented = errors.New("not implemented")
)

// Prism Errors
var (
	ErrInvalidAnchorFormat = errors.New("invalid anchor format")
)
