package inmemory

import "cadence/pkg/cadence/domain"

func NewBandRepository(store *Store) domain.BandRepository {
	return &bandRepository{mapRepo: store.bands}
}

type bandRepository struct {
	*mapRepo[*domain.Band]
}

func (r *bandRepository) GetByInviteCode(inviteCode string) (*domain.Band, error) {
	for _, band := range r.items {
		if band.InviteCode() == inviteCode {
			return band, nil
		}
	}
	return nil, domain.ErrBandNotFound
}
