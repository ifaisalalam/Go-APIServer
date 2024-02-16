package auth

import (
	"net/http"
)

type StrategyName string

type Strategy interface {
	Authenticate(req *http.Request) (bool, error)
	Is(name StrategyName) bool
}

type StrategyFunc func(req *http.Request) (bool, error)

func (f StrategyFunc) Authenticate(req *http.Request) (bool, error) {
	return f(req)
}
