package shortener

import (
	"errors"
)

var (
	// ErrShortURLAlreadyPresent is returned when short URL cannot be created because of a conflict.
	ErrShortURLAlreadyPresent = errors.New("requested short url already exists")
	// ErrShortURLDoesNotExist is returned when the short URL does not exist in the data cache.
	ErrShortURLDoesNotExist = errors.New("requested short url is invalid")
	// ErrInvalidShortURL is returned when the entered short URL is not allowed to be used.
	ErrInvalidShortURL = errors.New("short url contains unallowed characters")
	// ErrInvalidTargetURL is returned when the target URL is not syntactically correct.
	ErrInvalidTargetURL = errors.New("target url is not correct")
)
