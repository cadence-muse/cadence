package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-faster/errors"

	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/maybe"
	"cadence/pkg/common/postgresql"
	"cadence/pkg/common/uuid"
	"cadence/pkg/common/valuetypes"
)

func NewTrackRepository(
	ctx context.Context,
	client postgresql.ClientContext,
	subjectID uuid.UUID,
) domain.TrackRepository {
	return &trackRepository{
		ctx:       ctx,
		client:    client,
		subjectID: subjectID,
	}
}

type trackRepository struct {
	ctx       context.Context
	client    postgresql.ClientContext
	subjectID uuid.UUID
}

func (repo *trackRepository) NextID() domain.TrackID {
	return uuid.Generate()
}

func (repo *trackRepository) Store(track *domain.Track) error {
	const sqlQuery = `
		INSERT INTO track (
			id, band_id, title, artist, duration_seconds, original_tempo, original_key,
			custom_tempo, custom_key, notes, created_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE
		SET band_id          = EXCLUDED.band_id,
		    title            = EXCLUDED.title,
		    artist           = EXCLUDED.artist,
		    duration_seconds = EXCLUDED.duration_seconds,
		    original_tempo   = EXCLUDED.original_tempo,
		    original_key     = EXCLUDED.original_key,
		    custom_tempo     = EXCLUDED.custom_tempo,
		    custom_key       = EXCLUDED.custom_key,
		    notes            = EXCLUDED.notes,
		    updated_at       = now(),
		    updated_by       = EXCLUDED.created_by
	`

	customTempo := maybe.NewNone[int]()
	if tempo, ok := maybe.JustValid(track.CustomTempo()); ok {
		customTempo = maybe.NewJust(tempo)
	}

	customKey := maybe.NewNone[string]()
	if key, ok := maybe.JustValid(track.CustomKey()); ok {
		customKey = maybe.NewJust(key.String())
	}

	_, err := repo.client.ExecContext(
		repo.ctx,
		sqlQuery,
		track.ID(),
		track.BandID(),
		track.Title(),
		track.Artist(),
		int(track.Duration()/time.Second),
		track.OriginalTempo(),
		track.OriginalKey().String(),
		customTempo,
		customKey,
		track.Notes(),
		repo.subjectID,
	)
	return err
}

func (repo *trackRepository) Get(id domain.TrackID) (*domain.Track, error) {
	const sqlQuery = `
		SELECT id, band_id, title, artist,
		       duration_seconds, original_tempo, original_key,
		       custom_tempo, custom_key, notes
		FROM track
		WHERE id = $1
		`

	var row sqlxTrack
	err := repo.client.GetContext(repo.ctx, &row, sqlQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTrackNotFound
		}
		return nil, err
	}

	originalKey, err := valuetypes.MakeKey(row.OriginalKey)
	if err != nil {
		return nil, err
	}

	var customKey maybe.Maybe[valuetypes.MusicalKey]
	if value, ok := maybe.JustValid(row.CustomKey); ok {
		key, keyErr := valuetypes.MakeKey(value)
		if keyErr != nil {
			return nil, err
		}
		customKey = maybe.NewJust(key)
	}

	return domain.LoadTrack(
		row.ID,
		row.BandID,
		row.Title,
		row.Artist,
		time.Duration(row.DurationSeconds)*time.Second,
		row.OriginalTempo,
		originalKey,
		row.CustomTempo,
		customKey,
		row.Notes,
	), nil
}

type sqlxTrack struct {
	ID              uuid.UUID           `db:"id"`
	BandID          uuid.UUID           `db:"band_id"`
	Title           string              `db:"title"`
	Artist          string              `db:"artist"`
	DurationSeconds int                 `db:"duration_seconds"`
	OriginalTempo   int                 `db:"original_tempo"`
	OriginalKey     string              `db:"original_key"`
	CustomTempo     maybe.Maybe[int]    `db:"custom_tempo"`
	CustomKey       maybe.Maybe[string] `db:"custom_key"`
	Notes           maybe.Maybe[string] `db:"notes"`
}
