package repo

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/common/postgresql"
	"cadence/pkg/common/transactional"
	"cadence/pkg/common/uuid"
)

const advisoryLockSQL = `SELECT pg_advisory_xact_lock(hashtext($1))`

func NewTransactionFactory(connectionProvider postgresql.ConnectionProvider) transactional.TransactionFactory[app.RepoProvider] {
	return &transactionFactory{connectionProvider: connectionProvider}
}

type transactionFactory struct {
	connectionProvider postgresql.ConnectionProvider
}

func (f *transactionFactory) NewLockableTransaction(ctx context.Context, lockName string) (app.RepoProvider, error) {
	conn, err := f.connectionProvider.Connection(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := conn.BeginTransaction(ctx, nil)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	if lockName != "" {
		if _, err = tx.ExecContext(ctx, advisoryLockSQL, lockName); err != nil {
			_ = tx.Rollback()
			_ = conn.Close()
			return nil, err
		}
	}

	return NewRepoProvider(ctx, conn, tx, uuid.UUID{}), nil
}
