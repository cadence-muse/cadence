package query

import (
	"context"

	"cadence/pkg/cadence/app/query"
	"cadence/pkg/common/postgresql"
	"cadence/pkg/common/uuid"
)

func NewTrackQueryService(client postgresql.ClientContext) query.TrackQueryService {
	return &trackQueryService{client: client}
}

type trackQueryService struct {
	client postgresql.ClientContext
}

func (s *trackQueryService) ListBandTracks(ctx context.Context, bandID uuid.UUID) ([]query.TrackListItem, error) {
	// TODO implement
	panic("implement me")
}
