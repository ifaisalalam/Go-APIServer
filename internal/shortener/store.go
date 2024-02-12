package shortener

import (
	"context"
	"errors"
	mysqlErr "github.com/go-mysql/errors"
	"github.com/ifaisalalam/Go-awesome-service/pkg/mysql"
	"gorm.io/gorm"
)

var (
	ErrConflict = errors.New("conflict error")
)

type ShortURLDTO struct {
	Id        uint64
	ShortURL  string
	TargetURL string
	CreatedAt []uint8
	UpdatedAt []uint8
}

func (ShortURLDTO) TableName() string {
	return "short_urls"
}

type Store interface {
	SaveShortURL(ctx context.Context, shortUrl *ShortURLDTO) error
	Find(ctx context.Context, shortUrl string) (*ShortURLDTO, error)
}

func NewStore(db *gorm.DB) Store {
	return &store{db: db}
}

type store struct {
	db *gorm.DB
}

func (d *store) SaveShortURL(ctx context.Context, shortUrl *ShortURLDTO) error {
	err := mysql.WithRetry(1, func() error {
		return d.db.WithContext(ctx).Create(shortUrl).Error
	})
	if ok, err := mysqlErr.Error(err); ok && errors.Is(err, mysqlErr.ErrDupeKey) {
		return ErrConflict
	}
	return err
}

func (d *store) Find(ctx context.Context, shortUrl string) (*ShortURLDTO, error) {
	dto := new(ShortURLDTO)

	err := mysql.WithRetry(1, func() error {
		return d.db.WithContext(ctx).Where(&ShortURLDTO{ShortURL: shortUrl}).Limit(1).Find(dto).Error
	})
	if err != nil {
		return nil, err
	}

	return dto, nil
}
