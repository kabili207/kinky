package redditbooru

import "errors"

var (
	ErrInvalidConfig = errors.New("invalid config")
	ErrImageNotFound = errors.New("image not found")
	ErrRequestFailed = errors.New("request failed")
)
