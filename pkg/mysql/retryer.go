package mysql

import (
	mysqlErr "github.com/go-mysql/errors"
)

const (
	DefaultRetryAttempt = 1
)

type Retryer interface {
	WithRetry(f func() error) error
}

func NewDefaultRetryer() Retryer {
	return retryer{Attempts: DefaultRetryAttempt}
}

func NewRetryer(attempts uint) Retryer {
	return retryer{Attempts: attempts}
}

type retryer struct {
	Attempts uint
}

func (r retryer) WithRetry(f func() error) (err error) {
	if r.Attempts == 0 {
		return f()
	}

	var i uint
	for i = 0; i < r.Attempts; i++ {
		if err = f(); err == nil || !mysqlErr.CanRetry(err) {
			return
		}
	}
	return
}

func WithRetry(attempt int, f func() error) (err error) {
	if attempt <= 0 {
		return f()
	}

	for i := 0; i < attempt; i++ {
		if err = f(); err == nil || !mysqlErr.CanRetry(err) {
			return
		}
	}
	return
}
