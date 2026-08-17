package authorized // nolint:dupl

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/app/query"
	"cadence/pkg/common/maybe"
	"cadence/pkg/common/transactional"
	"cadence/pkg/common/uuid"
)

func NewTrackQueryService(inner query.TrackQueryService, executor transactional.Executor[app.RepoProvider]) query.TrackQueryService {
	return &authorizedTrackQueryService{next: inner, executor: executor}
}

type authorizedTrackQueryService struct {
	next     query.TrackQueryService
	executor transactional.Executor[app.RepoProvider]
}

func (s *authorizedTrackQueryService) ListBandTracks(ctx context.Context, bandID, requesterID uuid.UUID) ([]query.TrackListItem, error) {
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return nil, err
	}
	return s.next.ListBandTracks(ctx, bandID, requesterID)
}

func (s *authorizedTrackQueryService) ListUserTracks(ctx context.Context, userID uuid.UUID, bandID maybe.Maybe[uuid.UUID]) ([]query.UserTrackListItem, error) {
	return s.next.ListUserTracks(ctx, userID, bandID)
}

func (s *authorizedTrackQueryService) FindTrack(ctx context.Context, bandID, trackID, requesterID uuid.UUID) (maybe.Maybe[query.TrackData], error) {
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return maybe.NewNone[query.TrackData](), err
	}
	return s.next.FindTrack(ctx, bandID, trackID, requesterID)
}
