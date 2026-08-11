package main

import (
	"fmt"

	"github.com/caarlos0/env/v11"
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
	DBPorts    string `env:"db_post"`
	DBName     string `env:"db_name"`
	DBUser     string `env:"db_user"`
	DBPassword string `env:"db_password"`
}
