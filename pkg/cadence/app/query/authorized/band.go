package authorized

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/app/query"
	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/auth"
	"cadence/pkg/common/maybe"
	commonogenerrors "cadence/pkg/common/ogenerrors"
	"cadence/pkg/common/transactional"
	"cadence/pkg/common/uuid"
)

func NewBandQueryService(inner query.BandQueryService, executor transactional.Executor[app.RepoProvider]) query.BandQueryService {
	return &authorizedBandQueryService{next: inner, executor: executor}
}

type authorizedBandQueryService struct {
	next     query.BandQueryService
	executor transactional.Executor[app.RepoProvider]
}

func (s *authorizedBandQueryService) FindBand(ctx context.Context, id uuid.UUID) (maybe.Maybe[query.BandData], error) {
	requesterID, err := requesterIDFromContext(ctx)
	if err != nil {
		return maybe.NewNone[query.BandData](), err
	}
	err = requireMember(ctx, s.executor, id, requesterID)
	if err != nil {
		return maybe.NewNone[query.BandData](), err
	}
	return s.next.FindBand(ctx, id)
}

func (s *authorizedBandQueryService) ListUserBands(ctx context.Context, userID uuid.UUID) ([]query.BandListItem, error) {
	return s.next.ListUserBands(ctx, userID)
}

func (s *authorizedBandQueryService) CountUserBands(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.next.CountUserBands(ctx, userID)
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
