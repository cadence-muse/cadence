package inmemory

import "cadence/pkg/cadence/domain"

func NewSetlistRepository(store *Store) domain.SetlistRepository {
	return store.setlists
}
