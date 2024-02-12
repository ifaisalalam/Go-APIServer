package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/ifaisalalam/Go-awesome-service/restapi/operations"
)

// NewLivenessHandler returns handler for the internal Ping call.
func NewLivenessHandler() operations.LivenessHandler {
	return operations.LivenessHandlerFunc(func(p operations.LivenessParams) middleware.Responder {
		return operations.NewLivenessOK()
	})
}
