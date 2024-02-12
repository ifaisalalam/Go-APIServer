package shortener

import (
	"context"
)

// Shortener service interface.
type Shortener interface {
	// CreateShortURL creates a new short URL. If the short URL value is already present, it returns a ErrShortURLAlreadyPresent.
	CreateShortURL(ctx context.Context, input *CreateShortURLInput) (*CreateShortURLOutput, error)
	// GetTargetURL returns the target long URL for the given ShortURL.
	GetTargetURL(ctx context.Context, input *GetTargetURLInput) (*GetTargetURLOutput, error)
}
