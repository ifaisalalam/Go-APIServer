package mysql

import (
	"time"

	driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Dsn          string
	Address      string
	User         string
	Password     string
	Database     string
	ConnTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewClient(config *Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: config.Dsn,
		DSNConfig: &driver.Config{
			User:         config.User,
			Passwd:       config.Password,
			Addr:         config.Address,
			DBName:       config.Database,
			Timeout:      config.ConnTimeout,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		},
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, err
}
