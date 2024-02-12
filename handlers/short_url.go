package handlers

import (
	"errors"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	shortenerService "github.com/ifaisalalam/Go-awesome-service/internal/shortener"
	"github.com/ifaisalalam/Go-awesome-service/models"
	"github.com/ifaisalalam/Go-awesome-service/restapi/operations"
)

// ShortURLHandler ...
type ShortURLHandler interface {
	CreateShortURLHandler() operations.CreateShortURLHandler
	GetTargetURLHandler() operations.GetTargetURLHandler
}

// NewShortURLHandler returns handler for the internal Readiness call.
func NewShortURLHandler(shortener shortenerService.Shortener) ShortURLHandler {
	return &shortURLHandler{
		shortener: shortener,
	}
}

type shortURLHandler struct {
	shortener shortenerService.Shortener
}

func (h *shortURLHandler) CreateShortURLHandler() operations.CreateShortURLHandler {
	return operations.CreateShortURLHandlerFunc(func(p operations.CreateShortURLParams) middleware.Responder {

		input := &shortenerService.CreateShortURLInput{
			LongURL:  swag.StringValue(p.Body.LongURL),
			ShortURL: swag.StringValue(p.Body.ShortURL),
		}

		resp, err := h.shortener.CreateShortURL(p.HTTPRequest.Context(), input)
		if errors.Is(err, shortenerService.ErrShortURLAlreadyPresent) {
			return operations.NewCreateShortURLConflict().WithPayload(&models.ErrorResponse{
				Message: swag.String(err.Error()),
			})
		}
		if errors.Is(err, shortenerService.ErrInvalidTargetURL) || errors.Is(err, shortenerService.ErrInvalidShortURL) {
			return operations.NewCreateShortURLBadRequest().WithPayload(&models.ErrorResponse{
				Message: swag.String(err.Error()),
			})
		}
		if err != nil {
			return operations.NewCreateShortURLInternalServerError()
		}

		return operations.NewCreateShortURLCreated().WithPayload(&models.CreateShortURLResponse{
			ShortURL: swag.String(resp.ShortURL),
		})
	})
}

func (h *shortURLHandler) GetTargetURLHandler() operations.GetTargetURLHandler {
	return operations.GetTargetURLHandlerFunc(func(p operations.GetTargetURLParams) middleware.Responder {

		input := &shortenerService.GetTargetURLInput{
			ShortURL: p.ShortURL,
		}
		resp, err := h.shortener.GetTargetURL(p.HTTPRequest.Context(), input)
		if errors.Is(err, shortenerService.ErrShortURLDoesNotExist) {
			return operations.NewGetTargetURLNotFound().WithPayload(&models.ErrorResponse{
				Message: swag.String(err.Error()),
			})
		}
		if err != nil {
			return operations.NewGetTargetURLInternalServerError()
		}

		return operations.NewGetTargetURLOK().WithPayload(&models.GetTargetURLResponse{
			TargetURL: swag.String(resp.LongURL),
		})
	})
}
