package main

import (
	"errors"
	"time"

	"cadence/pkg/common/log"
	"cadence/pkg/common/postgresql"
)

var errMigrationFinished = errors.New("migration finished without errors")

func migrate(config *config, logger log.Logger) error {
	_, err := newDatabaseMigrator(config, logger)
	if err != nil {
		return err
	}
	return errMigrationFinished
}

type databaseMigrator struct {
	connector postgresql.Connector
}

func newDatabaseMigrator(
	config *config,
	logger log.Logger,
) (*databaseMigrator, error) {
	connector := postgresql.NewConnector()
	err := openWithRetries(connector, config.dsn(), config.DBMaxConn, config.DBConnLifetime, logger)
	if err != nil {
		return nil, err
	}

	m, err := connector.Migrator(logger)
	if err != nil {
		return nil, err
	}

	err = m.MigrateUp()
	if err != nil {
		return nil, err
	}

	return &databaseMigrator{connector: connector}, nil
}

func openWithRetries(
	connector postgresql.Connector,
	dsn postgresql.DSN,
	dbMaxConn int,
	dbConnectionLifetime int,
	logger log.Logger,
) (err error) {
	const retryCount = 6
	const interval = time.Second * 5
	for i := 0; i < retryCount; i++ {
		err = connector.Open(dsn, postgresql.Config{
			MaxConnections:     dbMaxConn,
			ConnectionLifetime: time.Duration(dbConnectionLifetime) * time.Second,
		})
		if err == nil {
			return nil
		}

		logger.Info("Retrying connection to DB...")
		time.Sleep(interval)
	}
	return err
}
