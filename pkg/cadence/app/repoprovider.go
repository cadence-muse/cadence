package app

import "cadence/pkg/cadence/domain"

type RepoProvider interface {
	UserRepository() domain.UserRepository
	BandRepository() domain.BandRepository

	Complete(err error) error
}
