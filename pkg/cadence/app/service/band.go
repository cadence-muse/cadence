package service

import (
	"context"
	"fmt"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/maybe"
	"cadence/pkg/common/transactional"
	"cadence/pkg/common/uuid"
)

func NewBandService(executor transactional.Executor[app.RepoProvider]) BandService {
	return &bandService{executor: executor}
}

type BandService interface {
	Create(ctx context.Context, params CreateBandParams) (uuid.UUID, error)
	Update(ctx context.Context, params UpdateBandParams, requesterID uuid.UUID) error
	JoinByInviteCode(ctx context.Context, userID uuid.UUID, inviteCode string) error
	Remove(ctx context.Context, bandID, requesterID uuid.UUID) error
	RemoveMember(ctx context.Context, bandID, targetUserID, requesterID uuid.UUID) error
}

type bandService struct {
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

func (s *bandService) Create(ctx context.Context, params CreateBandParams) (bandID uuid.UUID, err error) {
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

func (s *bandService) Update(ctx context.Context, params UpdateBandParams, _ uuid.UUID) (err error) {
	return s.executor.ExecuteWithLock(ctx, getBandLockName(params.BandID), func(repoProvider app.RepoProvider) error {
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

func (s *bandService) JoinByInviteCode(ctx context.Context, userID uuid.UUID, inviteCode string) error {
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

func (s *bandService) Remove(ctx context.Context, bandID, requesterID uuid.UUID) error {
	return s.executor.ExecuteWithLock(ctx, getBandLockName(bandID), func(repoProvider app.RepoProvider) error {
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

func (s *bandService) RemoveMember(ctx context.Context, bandID, targetUserID, requesterID uuid.UUID) error {
	return s.executor.ExecuteWithLock(ctx, getBandLockName(bandID), func(repoProvider app.RepoProvider) error {
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

func getBandLockName(id uuid.UUID) string {
	return fmt.Sprintf("band_%s", id.String())
}
