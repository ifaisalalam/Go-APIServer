package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/ifaisalalam/Go-awesome-service/restapi/operations"
)

// NewReadinessHandler returns handler for the internal Readiness call.
func NewReadinessHandler() operations.ReadinessHandler {
	return operations.ReadinessHandlerFunc(func(p operations.ReadinessParams) middleware.Responder {
		return operations.NewReadinessOK()
	})
}
