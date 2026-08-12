package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"

	redisinfra "cadence/pkg/cadence/infrastructure/persistence/redis"
	"cadence/pkg/common/postgresql"
	"cadence/pkg/common/redis"
)

func parseEnv() (*config, error) {
	c := new(config)
	if err := env.ParseWithOptions(c, env.Options{Prefix: strings.ToUpper(appID) + "_"}); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}
	return c, nil
}

type config struct {
	ServeRESTAddress string `env:"SERVE_REST_ADDRESS" envDefault:":8080"`

	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`

	DBMaxConn      int `env:"DB_MAX_CONN" envDefault:"10"`
	DBConnLifetime int `env:"DB_CONN_LIFETIME" envDefault:"60"`

	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     int    `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`

	SessionMaxPerUser int           `env:"SESSION_MAX_PER_USER" envDefault:"5"`
	SessionTTL        time.Duration `env:"SESSION_TTL" envDefault:"24h"`
}

func (c *config) dsn() postgresql.DSN {
	return postgresql.DSN{
		Host:     c.DBHost,
		Port:     c.DBPort,
		Database: c.DBName,
		User:     c.DBUser,
		Password: c.DBPassword,
	}
}

func (c *config) redisConfig() redis.Config {
	return redis.Config{
		Host:     c.RedisHost,
		Port:     c.RedisPort,
		Password: c.RedisPassword,
		DB:       c.RedisDB,
	}
}

func (c *config) sessionStoreConfig() redisinfra.Config {
	return redisinfra.Config{
		TTL:                c.SessionTTL,
		MaxSessionsPerUser: c.SessionMaxPerUser,
	}
}
