package inmemory

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/common/transactional"
)

func NewTransactionFactory(store *Store) transactional.TransactionFactory[app.RepoProvider] {
	return &transactionFactory{store: store}
}

type transactionFactory struct {
	store *Store
}

func (f *transactionFactory) NewLockableTransaction(_ context.Context, _ string) (app.RepoProvider, error) {
	return NewRepoProvider(f.store), nil
}
