package postgresql

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // registers the "postgres" database/sql driver
	"github.com/jmoiron/sqlx"

	"cadence/pkg/common/log"
)

const dbDriverName = "pgx"

type Config struct {
	MaxConnections     int
	ConnectionLifetime time.Duration
	ConnectTimeout     time.Duration
}

func NewConnector() Connector {
	return &connector{}
}

type Connector interface {
	Open(dsn DSN, cfg Config) error
	TransactionalClient() TransactionalClient
	Migrator(logger log.Logger) (Migrator, error)
	Close() error
}

type connector struct {
	db *sqlx.DB
}

func (c *connector) Open(dsn DSN, cfg Config) (err error) {
	c.db, err = sqlx.Open(dbDriverName, dsn.String())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	c.db.SetMaxOpenConns(cfg.MaxConnections)
	c.db.SetConnMaxLifetime(cfg.ConnectionLifetime)
	return nil
}

func (c *connector) TransactionalClient() TransactionalClient {
	return &transactionalClient{c.db}
}

func (c *connector) Migrator(logger log.Logger) (Migrator, error) {
	if c.db == nil {
		return nil, errors.New("DB not initialized")
	}
	return &migrator{db: c.db, logger: logger}, nil
}

func (c *connector) Close() error {
	if c.db != nil {
		err := c.db.Close()
		return fmt.Errorf("failed to disconnect: %w", err)
	}
	return errors.New("DB not initialized")
}
