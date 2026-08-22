package inmemory

import (
	"context"

	"github.com/nightnoryu/go-kita/transactional"

	"cadence/pkg/cadence/app"
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
