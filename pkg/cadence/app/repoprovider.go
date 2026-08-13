package app

import "cadence/pkg/cadence/domain"

type RepoProvider interface {
	UserRepository() domain.UserRepository
	BandRepository() domain.BandRepository
	TrackRepository() domain.TrackRepository

	Complete(err error) error
}
