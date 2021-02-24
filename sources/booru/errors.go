package booru

import "errors"

var (
	ErrDialectNotSupported = errors.New("dialect not supported")
	ErrImageNotFound       = errors.New("image not found")
	ErrRequestFailed       = errors.New("request failed")
)
