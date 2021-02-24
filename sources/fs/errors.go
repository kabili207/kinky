package fs

import "errors"

var (
	ErrNoConfig      = errors.New("no config given")
	ErrInvalidConfig = errors.New("invalid config")
	ErrNotFolder     = errors.New("not a folder")
)
