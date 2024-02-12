package cache

import (
	"errors"
)

var (
	// ErrNil is returned when the Key does not exist in the cache.
	ErrNil = errors.New("nil")
)
