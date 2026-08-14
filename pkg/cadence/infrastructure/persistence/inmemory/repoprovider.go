package inmemory

import (
	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/domain"
)

func NewRepoProvider(store *Store) app.RepoProvider {
	return &repoProvider{store: store}
}

type repoProvider struct {
	store *Store
}

func (p *repoProvider) UserRepository() domain.UserRepository {
	return NewUserRepository(p.store)
}

func (p *repoProvider) BandRepository() domain.BandRepository {
	return NewBandRepository(p.store)
}

func (p *repoProvider) TrackRepository() domain.TrackRepository {
	return NewTrackRepository(p.store)
}

func (p *repoProvider) SetlistRepository() domain.SetlistRepository {
	return NewSetlistRepository(p.store)
}

func (p *repoProvider) Complete(err error) error {
	return err
}
