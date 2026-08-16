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

func NewTrackQueryService(client postgresql.ClientContext) query.TrackQueryService {
	return &trackQueryService{client: client}
}

type trackQueryService struct {
	client postgresql.ClientContext
}

func (s *trackQueryService) ListBandTracks(ctx context.Context, bandID uuid.UUID) ([]query.TrackListItem, error) {
	const sqlQuery = `
		SELECT id, title, artist, duration_seconds
		FROM track
		WHERE band_id = $1 AND deleted_at IS NULL
		ORDER BY title
	`
	var rows []sqlxTrackListItem
	if err := s.client.SelectContext(ctx, &rows, sqlQuery, bandID); err != nil {
		return nil, err
	}

	return slices.Map(rows, func(row sqlxTrackListItem) query.TrackListItem {
		return query.TrackListItem{
			ID:       row.ID,
			Title:    row.Title,
			Artist:   row.Artist,
			Duration: durationFromSeconds(row.DurationSeconds),
		}
	}), nil
}

func (s *trackQueryService) ListUserTracks(ctx context.Context, userID uuid.UUID, bandID maybe.Maybe[uuid.UUID]) ([]query.UserTrackListItem, error) {
	const sqlQuery = `
		SELECT t.id, t.title, t.artist, t.duration_seconds, t.band_id, b.name AS band_name
		FROM track t
		JOIN band_member bm ON bm.band_id = t.band_id
		JOIN band b ON b.id = t.band_id
		WHERE bm.user_id = $1 AND t.deleted_at IS NULL AND b.deleted_at IS NULL
		  AND ($2::uuid IS NULL OR t.band_id = $2)
		ORDER BY t.title
	`
	var rows []sqlxUserTrackListItem
	if err := s.client.SelectContext(ctx, &rows, sqlQuery, userID, bandID); err != nil {
		return nil, err
	}

	return slices.Map(rows, func(row sqlxUserTrackListItem) query.UserTrackListItem {
		return query.UserTrackListItem{
			ID:       row.ID,
			Title:    row.Title,
			Artist:   row.Artist,
			Duration: durationFromSeconds(row.DurationSeconds),
			BandID:   row.BandID,
			BandName: row.BandName,
		}
	}), nil
}

func (s *trackQueryService) FindTrack(ctx context.Context, bandID, trackID uuid.UUID) (maybe.Maybe[query.TrackData], error) {
	const sqlQuery = `
		SELECT id, band_id, title, artist, duration_seconds, tempo, key, notes
		FROM track
		WHERE id = $1 AND band_id = $2 AND deleted_at IS NULL
	`
	var row sqlxTrackData
	err := s.client.GetContext(ctx, &row, sqlQuery, trackID, bandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return maybe.NewNone[query.TrackData](), nil
		}
		return maybe.NewNone[query.TrackData](), err
	}

	return maybe.NewJust(query.TrackData{
		ID:       row.ID,
		BandID:   row.BandID,
		Title:    row.Title,
		Artist:   row.Artist,
		Duration: durationFromSeconds(row.DurationSeconds),
		Tempo:    row.Tempo,
		Key:      row.Key,
		Notes:    row.Notes,
	}), nil
}

func durationFromSeconds(seconds maybe.Maybe[int]) maybe.Maybe[time.Duration] {
	value, ok := maybe.JustValid(seconds)
	if !ok {
		return maybe.NewNone[time.Duration]()
	}
	return maybe.NewJust(time.Duration(value) * time.Second)
}

type sqlxTrackListItem struct {
	ID              uuid.UUID        `db:"id"`
	Title           string           `db:"title"`
	Artist          string           `db:"artist"`
	DurationSeconds maybe.Maybe[int] `db:"duration_seconds"`
}

type sqlxUserTrackListItem struct {
	ID              uuid.UUID        `db:"id"`
	Title           string           `db:"title"`
	Artist          string           `db:"artist"`
	DurationSeconds maybe.Maybe[int] `db:"duration_seconds"`
	BandID          uuid.UUID        `db:"band_id"`
	BandName        string           `db:"band_name"`
}

type sqlxTrackData struct {
	ID              uuid.UUID           `db:"id"`
	BandID          uuid.UUID           `db:"band_id"`
	Title           string              `db:"title"`
	Artist          string              `db:"artist"`
	DurationSeconds maybe.Maybe[int]    `db:"duration_seconds"`
	Tempo           maybe.Maybe[int]    `db:"tempo"`
	Key             maybe.Maybe[string] `db:"key"`
	Notes           maybe.Maybe[string] `db:"notes"`
}
