/*
Package config provides application configuration handling. It internally uses
viper so that we can switch providers easily later should we need to.

Configuration are looked up, in the following order:
 1. Environment Variables
 2. .env files
 3. .env files in parent folder (in a shallow monorepo setup)
 4. Values given to viper.SetDefault

To add a new configuration entry, please do the following:
 1. Adds a new key to the const list. The value is the environment variable name
    to load its value from.
 2. Adds a viper.SetDefault call to the initViper method.
 3. Adds a getter for accessing the values.
*/
package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	*log.Logger
	Viper *viper.Viper
}

const (
	yesKey         = "ALWAYS_YES"
	listenAddrKey  = "LISTEN_ADDR"
	databaseUrlKey = "DATABASE_URL"
	redisUrlKey    = "REDIS_URL"

	apiPrefixKey  = "API_PREFIX"
	adminTokenKey = "ADMIN_TOKEN"
)

func MustConfigure() *Config {
	if cfg, err := Configure(); err != nil {
		log.Fatalln(err)
		return nil
	} else {
		return cfg
	}
}

func Configure() (*Config, error) {
	v, err := initViper()
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return &Config{
		Logger: log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		Viper:  v,
	}, nil
}

func initViper() (*viper.Viper, error) {
	v := viper.New()
	v.AutomaticEnv()
	if _, err := os.Stat(".env"); !errors.Is(err, os.ErrNotExist) {
		v.SetConfigFile(".env")
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	} else if _, err := os.Stat("../.env"); !errors.Is(err, os.ErrNotExist) {
		v.SetConfigFile("../.env")
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	v.SetDefault(yesKey, "0")
	v.SetDefault(listenAddrKey, ":4000")
	v.SetDefault(databaseUrlKey, "postgres:///ties?sslmode=disable")
	v.SetDefault(redisUrlKey, "redis:///")

	v.SetDefault(apiPrefixKey, "http://0.0.0.0:4000")
	v.SetDefault(adminTokenKey, "c3b1466dc5ce8d98b11da92a8589778c")
	return v, nil
}

func (c *Config) AlwaysYes() bool     { return c.Viper.GetBool(yesKey) }
func (c *Config) ListenAddr() string  { return c.Viper.GetString(listenAddrKey) }
func (c *Config) DatabaseURL() string { return c.Viper.GetString(databaseUrlKey) }
func (c *Config) RedisURL() string    { return c.Viper.GetString(redisUrlKey) }

func (c *Config) APIPrefix() string  { return c.Viper.GetString(apiPrefixKey) }
func (c *Config) AdminToken() string { return c.Viper.GetString(adminTokenKey) }

func (c *Config) AllConfigurations() map[string]interface{} {
	m := map[string]interface{}{}
	for _, key := range c.Viper.AllKeys() {
		m[key] = c.Viper.Get(key)
	}
	return m
}
