package apperr

import "errors"

var (
	ErrFailedToReadManifest        = errors.New("failed to read application manifest")
	ErrFailedToParseManifest       = errors.New("failed to parse application manifest")
	ErrFailedToValidateManifest    = errors.New("failed to validate application manifest")
	ErrInvalidApplicationBundle    = errors.New("invalid application bundle provided")
	ErrApplicationNoDisableSandbox = errors.New("an application cannot disable the sandbox")
	ErrEmptyExec                   = errors.New("an application cannot have an empty exec")
	ErrBlacklistedPath             = errors.New("a path the application tried to mount is blacklisted")
	ErrAttemptedPathTraversal      = errors.New("the application attempted path traversal")
	ErrInvalidUID                  = errors.New("the provided UID is invalid")
	ErrInvalidGID                  = errors.New("the provided GID is invalid")
	ErrFailedToFindDBus            = errors.New("failed to find dbus location")
	ErrDBusProxyFailed             = errors.New("failed to start dbus proxy")
	ErrDBusProxyTimeoutReached     = errors.New("timeout reached while waiting for dbus proxy to start")
)
