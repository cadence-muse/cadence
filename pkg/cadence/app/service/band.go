package service

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/transactional"
	"cadence/pkg/common/uuid"
)

func NewBandService(executor transactional.Executor[app.RepoProvider]) *BandService {
	return &BandService{executor: executor}
}

type BandService struct {
	executor transactional.Executor[app.RepoProvider]
}

func (s *BandService) Create(ctx context.Context, name string) (bandID uuid.UUID, err error) {
	err = s.executor.Execute(ctx, func(repoProvider app.RepoProvider) error {
		repo := repoProvider.BandRepository()

		band, bandErr := domain.NewBand(repo.NextID(), name)
		if bandErr != nil {
			return bandErr
		}

		if storeErr := repo.Store(band); storeErr != nil {
			return storeErr
		}

		bandID = band.ID()
		return nil
	})
	return bandID, err
}
