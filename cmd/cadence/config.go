package main

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"

	"cadence/pkg/common/postgresql"
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
