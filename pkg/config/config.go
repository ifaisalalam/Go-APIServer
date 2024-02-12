// Package config has specific primitives for loading application configurations.
//
// Primitives:
//   - Application should have struct for containing configuration. E.g. refer
//     internal/config/config.go file.
//   - Application should have a directory holding default file and environment
//     specific file. E.g. refer application/configuration/* directory.
//
// Usage:
//   - E.g. NewDefaultConfig().Load("dev", &config), where config is a struct
//     where configuration gets unmarshalled into.
package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Default options for configuration loading.
const (
	DefaultConfigType     = "toml"
	DefaultConfigFileName = "default"
	DirPathEnv            = "CONFIG_FILE_PATH"
	FileNameEnv           = "CONFIG_FILE"
)

// Options is config options.
type Options struct {
	configType     string
	configDirPath  string
	configFileName string
}

func (o Options) WithConfigType(configType string) Options {
	o.configType = configType
	return o
}

func (o Options) WithConfigDirPath(configDirPath string) Options {
	o.configDirPath = configDirPath
	return o
}

func (o Options) WithConfigFileName(configFileName string) Options {
	o.configFileName = configFileName
	return o
}

// Config is a wrapper over an underlying config loader implementation.
type Config struct {
	opts  Options
	viper *viper.Viper
}

// NewDefaultOptions returns default options.
func NewDefaultOptions() Options {
	configDirPath := os.Getenv(DirPathEnv)
	fileName := os.Getenv(FileNameEnv)
	if fileName == "" {
		fileName = DefaultConfigFileName
	}
	return NewOptions(DefaultConfigType, configDirPath, fileName)
}

// NewOptions returns new Options struct.
func NewOptions(configType string, configDirPath string, configFileName string) Options {
	return Options{configType: configType, configDirPath: configDirPath, configFileName: configFileName}
}

// NewDefaultConfig returns new config struct with default options.
func NewDefaultConfig() *Config {
	return NewConfig(NewDefaultOptions())
}

// NewConfig returns new config struct.
func NewConfig(opts Options) *Config {
	return &Config{opts, viper.New()}
}

// Load reads environment specific configurations and along with the defaults
// unmarshalls into config.
func (c *Config) Load(env string, config interface{}) error {
	if err := c.loadByConfigName(c.opts.configFileName, config); err != nil {
		return err
	}
	return c.loadByConfigName(env, config)
}

// loadByConfigName reads configuration from file and unmarshalls into config.
func (c *Config) loadByConfigName(configName string, config interface{}) error {
	c.viper.SetEnvPrefix(strings.ToUpper("api"))
	c.viper.SetConfigName(configName)
	c.viper.SetConfigType(c.opts.configType)
	c.viper.AddConfigPath(c.opts.configDirPath)
	c.viper.AutomaticEnv()
	c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}
	return c.viper.Unmarshal(config)
}
