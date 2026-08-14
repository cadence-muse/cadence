package inmemory

import "cadence/pkg/cadence/domain"

func NewTrackRepository(store *Store) domain.TrackRepository {
	return store.tracks
}
