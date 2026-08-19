package query

import (
	"context"
	"time"

	"cadence/pkg/common/maybe"
	"cadence/pkg/common/uuid"
)

type TrackQueryService interface {
	ListBandTracks(ctx context.Context, bandID uuid.UUID) ([]TrackListItem, error)
	ListUserTracks(ctx context.Context, userID uuid.UUID, bandID maybe.Maybe[uuid.UUID], searchQuery maybe.Maybe[string]) ([]UserTrackListItem, error)
	FindTrack(ctx context.Context, bandID, trackID uuid.UUID) (maybe.Maybe[TrackData], error)
}

type TrackListItem struct {
	ID       uuid.UUID
	Title    string
	Artist   string
	Duration maybe.Maybe[time.Duration]
}

type UserTrackListItem struct {
	ID       uuid.UUID
	Title    string
	Artist   string
	Duration maybe.Maybe[time.Duration]
	BandID   uuid.UUID
	BandName string
}

type TrackData struct {
	ID       uuid.UUID
	BandID   uuid.UUID
	Title    string
	Artist   string
	Duration maybe.Maybe[time.Duration]
	Tempo    maybe.Maybe[int]
	Key      maybe.Maybe[string]
	Notes    maybe.Maybe[string]
}
