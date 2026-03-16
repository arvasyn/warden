package apperr

import "errors"

var (
	ErrApplicationNoDisableSandbox = errors.New("an application cannot disable the sandbox")
	ErrEmptyExec                   = errors.New("an application cannot have an empty exec")
	ErrBlacklistedPath             = errors.New("a path the application tried to mount is blacklisted")
)
