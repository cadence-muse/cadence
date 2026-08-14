package service

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/maybe"
	"cadence/pkg/common/transactional"
	"cadence/pkg/common/uuid"
)

func NewBandService(executor transactional.Executor[app.RepoProvider]) *BandService {
	return &BandService{executor: executor}
}

type BandService struct {
	executor transactional.Executor[app.RepoProvider]
}

type CreateBandParams struct {
	OwnerID uuid.UUID
	Name    string
}

type UpdateBandParams struct {
	BandID uuid.UUID
	Name   maybe.Maybe[string]
}

func (s *BandService) Create(ctx context.Context, params CreateBandParams) (bandID uuid.UUID, err error) {
	err = s.executor.Execute(ctx, func(repoProvider app.RepoProvider) error {
		repo := repoProvider.BandRepository()

		band, bandErr := domain.NewBand(repo.NextID(), params.Name, params.OwnerID)
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

func (s *BandService) Update(ctx context.Context, params UpdateBandParams) (err error) {
	return s.executor.Execute(ctx, func(repoProvider app.RepoProvider) error {
		repo := repoProvider.BandRepository()

		band, err := repo.Get(params.BandID)
		if err != nil {
			return err
		}

		if name, ok := maybe.JustValid(params.Name); ok {
			if err := band.SetName(name); err != nil {
				return err
			}
		}

		return repo.Store(band)
	})
}

func (s *BandService) JoinByInviteCode(ctx context.Context, userID uuid.UUID, inviteCode string) error {
	return s.executor.Execute(ctx, func(repoProvider app.RepoProvider) error {
		repo := repoProvider.BandRepository()

		band, err := repo.GetByInviteCode(inviteCode)
		if err != nil {
			return err
		}

		band.AddMember(userID, domain.BandRoleMember)
		return repo.Store(band)
	})
}

func (s *BandService) Remove(ctx context.Context, bandID, requesterID uuid.UUID) error {
	return s.executor.Execute(ctx, func(repoProvider app.RepoProvider) error {
		repo := repoProvider.BandRepository()

		band, err := repo.Get(bandID)
		if err != nil {
			return err
		}

		if !band.IsOwner(requesterID) {
			return domain.ErrNotBandOwner
		}

		return repo.Remove(bandID)
	})
}

func (s *BandService) RemoveMember(ctx context.Context, bandID, targetUserID, requesterID uuid.UUID) error {
	return s.executor.Execute(ctx, func(repoProvider app.RepoProvider) error {
		repo := repoProvider.BandRepository()

		band, err := repo.Get(bandID)
		if err != nil {
			return err
		}

		if requesterID != targetUserID && !band.IsOwner(requesterID) {
			return domain.ErrNotBandOwner
		}
		if !band.HasMember(targetUserID) {
			return domain.ErrBandMemberNotFound
		}
		if band.IsOwner(targetUserID) {
			return domain.ErrCannotRemoveOwner
		}

		band.RemoveMember(targetUserID)
		return repo.Store(band)
	})
}
