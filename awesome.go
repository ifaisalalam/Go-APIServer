package awesome

import (
	"context"
	"github.com/ifaisalalam/Go-awesome-service/pkg/mysql"
	"log"

	"github.com/ifaisalalam/Go-awesome-service/internal/config"
	"github.com/ifaisalalam/Go-awesome-service/internal/shortener"
	"github.com/ifaisalalam/Go-awesome-service/pkg/cache"
	configreader "github.com/ifaisalalam/Go-awesome-service/pkg/config"
)

// Service contains an initialized instance of all the supported services.
type Service struct {
	Shortener shortener.Shortener
}

// NewService initializes all the services and dependencies.
func NewService(ctx context.Context, env string) *Service {
	var appConfig config.Config
	if err := configreader.NewDefaultConfig().Load(env, &appConfig); err != nil {
		log.Fatalln(err)
	}

	cacheConfig := appConfig.Cache
	redisConfig := cache.RedisConfig{
		Host:     cacheConfig.RedisConfig.Host,
		Port:     cacheConfig.RedisConfig.Port,
		Database: cacheConfig.RedisConfig.Database,
		Password: cacheConfig.RedisConfig.Password,
	}
	cacheStore, err := cache.NewCache(&cache.Config{
		Driver:      cacheConfig.Driver,
		RedisConfig: redisConfig,
	})
	if err != nil {
		log.Fatalln(err)
	}

	db := appConfig.MySQL
	dbConfig := &mysql.Config{
		Dsn:          db.Dsn,
		Address:      db.Address,
		User:         db.User,
		Password:     db.Password,
		Database:     db.Database,
		ConnTimeout:  db.ConnTimeout,
		ReadTimeout:  db.ReadTimeout,
		WriteTimeout: db.WriteTimeout,
	}
	dbStore, err := mysql.NewClient(dbConfig)
	if err != nil {
		log.Fatalln(err)
	}

	shortenerService := shortener.NewShortener(cacheStore, dbStore)

	return &Service{
		Shortener: shortenerService,
	}
}
