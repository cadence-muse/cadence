package main

import (
	"time"

	"github.com/nightnoryu/go-kita/postgresql"
	"github.com/nightnoryu/go-kita/redis"
)

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
	RedisPort     int    `env:"REDIS_PORT"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"`

	SessionMaxPerUser int           `env:"SESSION_MAX_PER_USER" envDefault:"5"`
	SessionTTL        time.Duration `env:"SESSION_TTL" envDefault:"24h"`

	CORSAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS"`
}

func (c *config) postgresDSN() postgresql.DSN {
	return postgresql.DSN{
		Host:     c.DBHost,
		Port:     c.DBPort,
		Database: c.DBName,
		User:     c.DBUser,
		Password: c.DBPassword,
	}
}

func (c *config) redisDSN() redis.DSN {
	return redis.DSN{
		Host:     c.RedisHost,
		Port:     c.RedisPort,
		Password: c.RedisPassword,
		DB:       c.RedisDB,
	}
}
