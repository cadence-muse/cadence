package authorized

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/app/service"
	"cadence/pkg/common/transactional"
	"cadence/pkg/common/uuid"
)

func NewTrackService(inner service.TrackService, executor transactional.Executor[app.RepoProvider]) service.TrackService {
	return &authorizedTrackService{next: inner, executor: executor}
}

type authorizedTrackService struct {
	next     service.TrackService
	executor transactional.Executor[app.RepoProvider]
}

func (s *authorizedTrackService) Create(ctx context.Context, params service.CreateTrackParams, requesterID uuid.UUID) (uuid.UUID, error) {
	if err := requireMember(ctx, s.executor, params.BandID, requesterID); err != nil {
		return uuid.UUID{}, err
	}
	return s.next.Create(ctx, params, requesterID)
}

func (s *authorizedTrackService) Update(ctx context.Context, params service.UpdateTrackParams, requesterID uuid.UUID) error {
	if err := requireMember(ctx, s.executor, params.BandID, requesterID); err != nil {
		return err
	}
	return s.next.Update(ctx, params, requesterID)
}

func (s *authorizedTrackService) Remove(ctx context.Context, bandID, trackID, requesterID uuid.UUID) error {
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return err
	}
	return s.next.Remove(ctx, bandID, trackID, requesterID)
}
