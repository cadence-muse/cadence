package service

import (
	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/infrastructure/persistence/inmemory"
	"cadence/pkg/common/transactional"
)

func newFakeExecutor() *fakeExecutor {
	store := inmemory.NewStore()
	return &fakeExecutor{
		Executor: transactional.NewExecutor(inmemory.NewTransactionFactory(store)),
		store:    store,
	}
}

type fakeExecutor struct {
	transactional.Executor[app.RepoProvider]
	store *inmemory.Store
}

func (e *fakeExecutor) repoProvider() app.RepoProvider {
	return inmemory.NewRepoProvider(e.store)
}
