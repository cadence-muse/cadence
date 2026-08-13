package query

import (
	"context"

	"cadence/pkg/common/uuid"
)

type TrackQueryService interface {
	ListBandTracks(ctx context.Context, bandID uuid.UUID) ([]TrackListItem, error)
}

type TrackListItem struct {
	ID     uuid.UUID
	Title  string
	Artist string
}
