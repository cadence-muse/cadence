package query

import (
	"context"
	"time"

	"cadence/pkg/common/maybe"
	"cadence/pkg/common/uuid"
)

type SetlistQueryService interface {
	ListBandSetlists(ctx context.Context, bandID uuid.UUID) ([]SetlistListItem, error)
	FindSetlist(ctx context.Context, bandID, setlistID uuid.UUID) (maybe.Maybe[SetlistData], error)
}

type SetlistListItem struct {
	ID        uuid.UUID
	Name      string
	EventDate maybe.Maybe[time.Time]
}

type SetlistData struct {
	ID            uuid.UUID
	BandID        uuid.UUID
	Name          string
	EventLocation maybe.Maybe[string]
	EventDate     maybe.Maybe[time.Time]
	Tracks        []SetlistTrackItem
}

type SetlistTrackItem struct {
	TrackID  uuid.UUID
	Position int
	Title    string
	Artist   string
	Duration maybe.Maybe[time.Duration]
}
