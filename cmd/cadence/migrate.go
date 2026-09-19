package main

import (
	"context"
	"errors"
	"time"

	"github.com/nightnoryu/go-kita/log"
	"github.com/nightnoryu/go-kita/postgresql"

	"cadence/data/migrations"
)

var errMigrationFinished = errors.New("migration finished without errors")

func migrate(ctx context.Context, config *config, logger log.Logger) error {
	_, err := newDatabaseMigrator(ctx, config, logger)
	if err != nil {
		return err
	}
	return errMigrationFinished
}

type databaseMigrator struct {
	connector postgresql.Connector
}

func newDatabaseMigrator(
	ctx context.Context,
	config *config,
	logger log.Logger,
) (*databaseMigrator, error) {
	connector := postgresql.NewConnector()
	err := openWithRetries(ctx, connector, config.postgresDSN(), config.DBMaxConn, config.DBConnLifetime, logger)
	if err != nil {
		return nil, err
	}

	m, err := connector.Migrator(logger, migrations.UpFS())
	if err != nil {
		return nil, err
	}

	err = m.MigrateUp(ctx)
	if err != nil {
		return nil, err
	}

	return &databaseMigrator{connector: connector}, nil
}

func openWithRetries(
	ctx context.Context,
	connector postgresql.Connector,
	dsn postgresql.DSN,
	dbMaxConn int,
	dbConnectionLifetime int,
	logger log.Logger,
) (err error) {
	const retryCount = 6
	const interval = time.Second * 5
	for i := 0; i < retryCount; i++ {
		err = connector.Open(ctx, dsn, postgresql.Config{
			MaxConnections:     dbMaxConn,
			ConnectionLifetime: time.Duration(dbConnectionLifetime) * time.Second,
		})
		if err == nil {
			return nil
		}

		logger.Info("Retrying connection to DB...")
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
		}
	}
	return err
}
