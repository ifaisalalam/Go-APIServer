package auth

import (
	"context"
	"log"
	"net/http"
)

type ctxKey int

const (
	CtxKeyPrincipal ctxKey = iota
)

func NewNoOpAuthMiddleware(next http.Handler) http.Handler {
	return NewAuthMiddlewareWithStrategy(next, NewApiKeyNoOpStrategy())
}

func NewAuthMiddlewareWithStrategy(next http.Handler, strategy Strategy) http.Handler {
	return &handler{
		strategy: strategy,
		next:     next,
	}
}

type handler struct {
	strategy Strategy
	next     http.Handler
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	var err error
	if authenticationOk, err := h.strategy.Authenticate(req); err == nil && authenticationOk {
		ctx := context.WithValue(req.Context(), CtxKeyPrincipal, "Faisal")
		h.next.ServeHTTP(writer, req.WithContext(ctx))
		return
	}
	if err != nil {
		log.Println(err)
	}

	writer.WriteHeader(http.StatusUnauthorized)
}

var _ http.Handler = &handler{}
