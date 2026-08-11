package main

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"cadence/pkg/common/postgresql"
)

func parseEnv() (*config, error) {
	c := new(config)
	if err := env.ParseWithOptions(c, env.Options{Prefix: appID}); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}
	return c, nil
}

type config struct {
	ServeRESTAddress string `env:"serve_rest_address" envDefault:":8080"`

	DBHost     string `env:"db_host"`
	DBPort     int    `env:"db_port"`
	DBName     string `env:"db_name"`
	DBUser     string `env:"db_user"`
	DBPassword string `env:"db_password"`
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
