package shortener

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/ifaisalalam/Go-awesome-service/pkg/cache"
)

// NewShortener returns a Shortener service implementation.
func NewShortener(cache cache.Cache, db *gorm.DB) Shortener {
	return &shortener{cache: cache, store: NewStore(db)}
}

type shortener struct {
	cache cache.Cache
	store Store
}

func (s *shortener) CreateShortURL(ctx context.Context, input *CreateShortURLInput) (*CreateShortURLOutput, error) {
	if err := s.validateCreateShortURL(input); err != nil {
		return nil, err
	}

	if _, err := s.cache.Get(ctx, input.ShortURL); !errors.Is(err, cache.ErrNil) {
		return nil, ErrShortURLAlreadyPresent
	}

	dto := new(ShortURLDTO)
	dto.ShortURL = input.ShortURL
	dto.TargetURL = input.LongURL
	if err := s.store.SaveShortURL(ctx, dto); err != nil {
		log.Println(err)
		return nil, err
	}

	if err := s.cache.Set(ctx, input.ShortURL, []byte(input.LongURL), time.Minute); err != nil {
		return nil, err
	}

	return &CreateShortURLOutput{ShortURL: input.ShortURL}, nil
}

func (s *shortener) GetTargetURL(ctx context.Context, input *GetTargetURLInput) (*GetTargetURLOutput, error) {
	value, err := s.cache.Get(ctx, input.ShortURL)
	if err == nil && string(value) != "" {
		return &GetTargetURLOutput{LongURL: string(value)}, nil
	}

	dto, err := s.store.Find(ctx, input.ShortURL)
	if err != nil || dto == nil || dto.TargetURL == "" {
		return nil, ErrShortURLDoesNotExist
	}

	_ = s.cache.Set(ctx, input.ShortURL, []byte(dto.TargetURL), time.Minute)

	return &GetTargetURLOutput{LongURL: dto.TargetURL}, nil
}

func (s *shortener) validateCreateShortURL(input *CreateShortURLInput) error {
	if !validateShortURL(input.ShortURL) {
		return ErrInvalidShortURL
	}
	if !validateTargetURL(input.LongURL) {
		return ErrInvalidTargetURL
	}
	return nil
}
