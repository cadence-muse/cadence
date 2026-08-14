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
)

func NewSetlistRepository(
	ctx context.Context,
	client postgresql.ClientContext,
	subjectID uuid.UUID,
) domain.SetlistRepository {
	return &setlistRepository{
		ctx:       ctx,
		client:    client,
		subjectID: subjectID,
	}
}

type setlistRepository struct {
	ctx       context.Context
	client    postgresql.ClientContext
	subjectID uuid.UUID
}

func (repo *setlistRepository) NextID() domain.SetlistID {
	return uuid.Generate()
}

func (repo *setlistRepository) Store(setlist *domain.Setlist) error {
	const sqlQuery = `
		INSERT INTO setlist (
			id, band_id, name, event_location, event_date,
			created_at, created_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE
		SET band_id        = EXCLUDED.band_id,
		    name           = EXCLUDED.name,
		    event_location = EXCLUDED.event_location,
		    event_date     = EXCLUDED.event_date,
		    updated_at     = EXCLUDED.created_at,
		    updated_by     = EXCLUDED.created_by
	`
	_, err := repo.client.ExecContext(
		repo.ctx,
		sqlQuery,
		setlist.ID(),
		setlist.BandID(),
		setlist.Name(),
		setlist.EventLocation(),
		setlist.EventDate(),
		time.Now(),
		repo.subjectID,
	)
	if err != nil {
		return err
	}

	return repo.replaceTracks(setlist.ID(), setlist.TrackIDs())
}

func (repo *setlistRepository) replaceTracks(setlistID domain.SetlistID, trackIDs []domain.TrackID) error {
	const deleteTracksSQL = `DELETE FROM setlist_track WHERE setlist_id = $1`
	if _, err := repo.client.ExecContext(repo.ctx, deleteTracksSQL, setlistID); err != nil {
		return err
	}

	const insertTrackSQL = `INSERT INTO setlist_track (setlist_id, track_id, position) VALUES ($1, $2, $3)`
	for position, trackID := range trackIDs {
		if _, err := repo.client.ExecContext(repo.ctx, insertTrackSQL, setlistID, trackID, position); err != nil {
			return err
		}
	}
	return nil
}

func (repo *setlistRepository) Get(id domain.SetlistID) (*domain.Setlist, error) {
	const sqlQuery = `
		SELECT id, band_id, name, event_location, event_date
		FROM setlist
		WHERE id = $1 AND deleted_at IS NULL
	`

	var row sqlxSetlist
	err := repo.client.GetContext(repo.ctx, &row, sqlQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSetlistNotFound
		}
		return nil, err
	}

	trackIDs, err := repo.getTrackIDs(id)
	if err != nil {
		return nil, err
	}

	return domain.LoadSetlist(
		row.ID,
		row.BandID,
		row.Name,
		row.EventLocation,
		row.EventDate,
		trackIDs,
	), nil
}

func (repo *setlistRepository) getTrackIDs(setlistID domain.SetlistID) ([]domain.TrackID, error) {
	const sqlQuery = `
		SELECT st.track_id
		FROM setlist_track st
		JOIN track t ON t.id = st.track_id AND t.deleted_at IS NULL
		WHERE st.setlist_id = $1
		ORDER BY st.position
	`
	var trackIDs []uuid.UUID
	if err := repo.client.SelectContext(repo.ctx, &trackIDs, sqlQuery, setlistID); err != nil {
		return nil, err
	}
	return trackIDs, nil
}

func (repo *setlistRepository) Remove(id domain.SetlistID) error {
	const sqlQuery = `
		UPDATE setlist
		SET deleted_at = $2,
		    deleted_by = $3
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := repo.client.ExecContext(repo.ctx, sqlQuery, id, time.Now(), repo.subjectID)
	return err
}

type sqlxSetlist struct {
	ID            uuid.UUID              `db:"id"`
	BandID        uuid.UUID              `db:"band_id"`
	Name          string                 `db:"name"`
	EventLocation maybe.Maybe[string]    `db:"event_location"`
	EventDate     maybe.Maybe[time.Time] `db:"event_date"`
}
