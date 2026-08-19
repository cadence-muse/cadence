package authorized

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/app/service"
	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/auth"
	commonogenerrors "cadence/pkg/common/ogenerrors"
	"cadence/pkg/common/transactional"
	"cadence/pkg/common/uuid"
)

func NewBandService(inner service.BandService, executor transactional.Executor[app.RepoProvider]) service.BandService {
	return &authorizedBandService{next: inner, executor: executor}
}

type authorizedBandService struct {
	next     service.BandService
	executor transactional.Executor[app.RepoProvider]
}

func (s *authorizedBandService) Create(ctx context.Context, params service.CreateBandParams) (uuid.UUID, error) {
	return s.next.Create(ctx, params)
}

func (s *authorizedBandService) Update(ctx context.Context, params service.UpdateBandParams) error {
	requesterID, err := requesterIDFromContext(ctx)
	if err != nil {
		return err
	}
	if err := requireMember(ctx, s.executor, params.BandID, requesterID); err != nil {
		return err
	}
	return s.next.Update(ctx, params)
}

func (s *authorizedBandService) JoinByInviteCode(ctx context.Context, userID uuid.UUID, inviteCode string) error {
	return s.next.JoinByInviteCode(ctx, userID, inviteCode)
}

func (s *authorizedBandService) Remove(ctx context.Context, bandID, requesterID uuid.UUID) error {
	return s.next.Remove(ctx, bandID, requesterID)
}

func (s *authorizedBandService) RemoveMember(ctx context.Context, bandID, targetUserID, requesterID uuid.UUID) error {
	return s.next.RemoveMember(ctx, bandID, targetUserID, requesterID)
}

func (s *authorizedBandService) TransferOwnership(ctx context.Context, bandID, requesterID, newOwnerID uuid.UUID) error {
	return s.next.TransferOwnership(ctx, bandID, requesterID, newOwnerID)
}

func (s *authorizedBandService) RegenerateInviteCode(ctx context.Context, bandID, requesterID uuid.UUID) (string, error) {
	return s.next.RegenerateInviteCode(ctx, bandID, requesterID)
}

func requireMember(ctx context.Context, executor transactional.Executor[app.RepoProvider], bandID, requesterID uuid.UUID) error {
	return executor.Execute(ctx, func(repoProvider app.RepoProvider) error {
		band, err := repoProvider.BandRepository().Get(bandID)
		if err != nil {
			return err
		}
		if !band.HasMember(requesterID) {
			return domain.ErrNotBandMember
		}
		return nil
	})
}

func requesterIDFromContext(ctx context.Context) (uuid.UUID, error) {
	requesterID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return uuid.UUID{}, commonogenerrors.NewPermissionDeniedError("user not authenticated")
	}
	return requesterID, nil
}
