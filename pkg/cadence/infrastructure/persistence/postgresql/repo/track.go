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
			id, band_id, title, artist,
			duration_seconds, tempo, key, notes,
		    created_at, created_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE
		SET band_id          = EXCLUDED.band_id,
		    title            = EXCLUDED.title,
		    artist           = EXCLUDED.artist,
		    duration_seconds = EXCLUDED.duration_seconds,
		    tempo            = EXCLUDED.tempo,
		    key              = EXCLUDED.key,
		    notes            = EXCLUDED.notes,
		    updated_at       = EXCLUDED.created_at,
		    updated_by       = EXCLUDED.created_by
	`

	tempo := maybe.NewNone[int]()
	if tempoValue, ok := maybe.JustValid(track.Tempo()); ok {
		tempo = maybe.NewJust(tempoValue)
	}

	key := maybe.NewNone[string]()
	if keyValue, ok := maybe.JustValid(track.Key()); ok {
		key = maybe.NewJust(keyValue.String())
	}

	duration := maybe.NewNone[int]()
	if durationValue, ok := maybe.JustValid(track.Duration()); ok {
		duration = maybe.NewJust(int(durationValue / time.Second))
	}

	_, err := repo.client.ExecContext(
		repo.ctx,
		sqlQuery,
		track.ID(),
		track.BandID(),
		track.Title(),
		track.Artist(),
		duration,
		tempo,
		key,
		track.Notes(),
		time.Now(),
		repo.subjectID,
	)
	return err
}

func (repo *trackRepository) Get(id domain.TrackID) (*domain.Track, error) {
	const sqlQuery = `
		SELECT id, band_id, title, artist,
		       duration_seconds, tempo, key, notes
		FROM track
		WHERE id = $1 AND deleted_at IS NULL
		`

	var row sqlxTrack
	err := repo.client.GetContext(repo.ctx, &row, sqlQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTrackNotFound
		}
		return nil, err
	}

	var key maybe.Maybe[valuetypes.MusicalKey]
	if value, ok := maybe.JustValid(row.Key); ok {
		keyValue, keyErr := valuetypes.MakeKey(value)
		if keyErr != nil {
			return nil, err
		}
		key = maybe.NewJust(keyValue)
	}

	var duration maybe.Maybe[time.Duration]
	if value, ok := maybe.JustValid(row.DurationSeconds); ok {
		duration = maybe.NewJust(time.Duration(value) * time.Second)
	}

	return domain.LoadTrack(
		row.ID,
		row.BandID,
		row.Title,
		row.Artist,
		duration,
		row.Tempo,
		key,
		row.Notes,
	), nil
}

func (repo *trackRepository) Remove(id domain.TrackID) error {
	const sqlQuery = `
		UPDATE track
		SET deleted_at = $2,
		    deleted_by = $3
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := repo.client.ExecContext(repo.ctx, sqlQuery, id, time.Now(), repo.subjectID)
	return err
}

type sqlxTrack struct {
	ID              uuid.UUID           `db:"id"`
	BandID          uuid.UUID           `db:"band_id"`
	Title           string              `db:"title"`
	Artist          string              `db:"artist"`
	DurationSeconds maybe.Maybe[int]    `db:"duration_seconds"`
	Tempo           maybe.Maybe[int]    `db:"tempo"`
	Key             maybe.Maybe[string] `db:"key"`
	Notes           maybe.Maybe[string] `db:"notes"`
}
