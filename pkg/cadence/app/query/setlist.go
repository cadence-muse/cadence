package query

import (
	"context"
	"time"

	"cadence/pkg/common/maybe"
	"cadence/pkg/common/uuid"
)

type SetlistQueryService interface {
	ListBandSetlists(ctx context.Context, bandID, requesterID uuid.UUID) ([]SetlistListItem, error)
	ListUserSetlists(ctx context.Context, userID uuid.UUID, bandID maybe.Maybe[uuid.UUID]) ([]UserSetlistListItem, error)
	FindSetlist(ctx context.Context, bandID, setlistID, requesterID uuid.UUID) (maybe.Maybe[SetlistData], error)
}

type SetlistListItem struct {
	ID          uuid.UUID
	Name        string
	TracksCount int
	Duration    time.Duration
	EventDate   maybe.Maybe[time.Time]
}

type UserSetlistListItem struct {
	ID          uuid.UUID
	Name        string
	TracksCount int
	Duration    time.Duration
	EventDate   maybe.Maybe[time.Time]
	BandID      uuid.UUID
	BandName    string
}

type SetlistData struct {
	ID            uuid.UUID
	BandID        uuid.UUID
	Name          string
	Duration      time.Duration
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
