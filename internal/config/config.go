package config

import (
	"github.com/ifaisalalam/Go-awesome-service/pkg/cache"
	"github.com/ifaisalalam/Go-awesome-service/pkg/mysql"
)

// Config is the base config structure for the service.
type Config struct {
	Cache cache.Config
	MySQL mysql.Config
}
