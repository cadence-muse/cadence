package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"cadence/pkg/cadence/app/query"
	"cadence/pkg/common/maybe"
	"cadence/pkg/common/postgresql"
	"cadence/pkg/common/slices"
	"cadence/pkg/common/uuid"
)

func NewSetlistQueryService(client postgresql.ClientContext) query.SetlistQueryService {
	return &setlistQueryService{client: client}
}

type setlistQueryService struct {
	client postgresql.ClientContext
}

func (s *setlistQueryService) ListBandSetlists(ctx context.Context, bandID uuid.UUID) ([]query.SetlistListItem, error) {
	const sqlQuery = `
		SELECT s.id, s.name, s.event_date,
		       COUNT(t.id) AS tracks_count,
		       COALESCE(SUM(t.duration_seconds), 0) AS duration_seconds
		FROM setlist s
		LEFT JOIN setlist_track st ON st.setlist_id = s.id
		LEFT JOIN track t ON t.id = st.track_id AND t.deleted_at IS NULL
		WHERE s.band_id = $1 AND s.deleted_at IS NULL
		GROUP BY s.id, s.name, s.event_date
		ORDER BY s.event_date ASC NULLS LAST
	`
	var rows []sqlxSetlistListItem
	if err := s.client.SelectContext(ctx, &rows, sqlQuery, bandID); err != nil {
		return nil, err
	}

	return slices.Map(rows, func(row sqlxSetlistListItem) query.SetlistListItem {
		return query.SetlistListItem{
			ID:          row.ID,
			Name:        row.Name,
			TracksCount: row.TracksCount,
			Duration:    time.Duration(row.DurationSeconds) * time.Second,
			EventDate:   row.EventDate,
		}
	}), nil
}

func (s *setlistQueryService) FindSetlist(ctx context.Context, bandID, setlistID uuid.UUID) (maybe.Maybe[query.SetlistData], error) {
	const sqlQuery = `
		SELECT id, band_id, name, event_location, event_date
		FROM setlist
		WHERE id = $1 AND band_id = $2 AND deleted_at IS NULL
	`
	var row sqlxSetlistData
	err := s.client.GetContext(ctx, &row, sqlQuery, setlistID, bandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return maybe.NewNone[query.SetlistData](), nil
		}
		return maybe.NewNone[query.SetlistData](), err
	}

	tracks, err := s.listSetlistTracks(ctx, setlistID)
	if err != nil {
		return maybe.NewNone[query.SetlistData](), err
	}

	var duration time.Duration
	for _, track := range tracks {
		if value, ok := maybe.JustValid(track.Duration); ok {
			duration += value
		}
	}

	return maybe.NewJust(query.SetlistData{
		ID:            row.ID,
		BandID:        row.BandID,
		Name:          row.Name,
		Duration:      duration,
		EventLocation: row.EventLocation,
		EventDate:     row.EventDate,
		Tracks:        tracks,
	}), nil
}

func (s *setlistQueryService) listSetlistTracks(ctx context.Context, setlistID uuid.UUID) ([]query.SetlistTrackItem, error) {
	const sqlQuery = `
		SELECT st.track_id, (ROW_NUMBER() OVER (ORDER BY st.position) - 1)::int AS position,
		       t.title, t.artist, t.duration_seconds
		FROM setlist_track st
		JOIN track t ON t.id = st.track_id AND t.deleted_at IS NULL
		WHERE st.setlist_id = $1
		ORDER BY st.position
	`
	var rows []sqlxSetlistTrackItem
	if err := s.client.SelectContext(ctx, &rows, sqlQuery, setlistID); err != nil {
		return nil, err
	}

	return slices.Map(rows, func(row sqlxSetlistTrackItem) query.SetlistTrackItem {
		return query.SetlistTrackItem{
			TrackID:  row.TrackID,
			Position: row.Position,
			Title:    row.Title,
			Artist:   row.Artist,
			Duration: durationFromSeconds(row.DurationSeconds),
		}
	}), nil
}

type sqlxSetlistListItem struct {
	ID              uuid.UUID              `db:"id"`
	Name            string                 `db:"name"`
	TracksCount     int                    `db:"tracks_count"`
	DurationSeconds int                    `db:"duration_seconds"`
	EventDate       maybe.Maybe[time.Time] `db:"event_date"`
}

type sqlxSetlistData struct {
	ID            uuid.UUID              `db:"id"`
	BandID        uuid.UUID              `db:"band_id"`
	Name          string                 `db:"name"`
	EventLocation maybe.Maybe[string]    `db:"event_location"`
	EventDate     maybe.Maybe[time.Time] `db:"event_date"`
}

type sqlxSetlistTrackItem struct {
	TrackID         uuid.UUID        `db:"track_id"`
	Position        int              `db:"position"`
	Title           string           `db:"title"`
	Artist          string           `db:"artist"`
	DurationSeconds maybe.Maybe[int] `db:"duration_seconds"`
}
