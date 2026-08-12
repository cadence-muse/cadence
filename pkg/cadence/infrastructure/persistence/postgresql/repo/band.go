package repo

import (
	"context"
	"database/sql"
	"errors"

	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/postgresql"
	"cadence/pkg/common/uuid"
)

func NewBandRepository(
	ctx context.Context,
	client postgresql.ClientContext,
	subjectID uuid.UUID,
) domain.BandRepository {
	return &bandRepository{
		ctx:       ctx,
		client:    client,
		subjectID: subjectID,
	}
}

type bandRepository struct {
	ctx       context.Context
	client    postgresql.ClientContext
	subjectID uuid.UUID
}

func (repo *bandRepository) NextID() domain.BandID {
	return uuid.Generate()
}

func (repo *bandRepository) Store(band *domain.Band) error {
	const sqlQuery = `
		INSERT INTO band (id, name, created_by)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET name       = EXCLUDED.name,
		    updated_at = now(),
		    updated_by = EXCLUDED.created_by
	`
	_, err := repo.client.ExecContext(repo.ctx, sqlQuery, band.ID(), band.Name(), repo.subjectID)
	return err
}

func (repo *bandRepository) Get(id domain.BandID) (*domain.Band, error) {
	const sqlQuery = `SELECT id, name FROM band WHERE id = $1`
	var row sqlxBand
	err := repo.client.GetContext(repo.ctx, &row, sqlQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBandNotFound
		}
		return nil, err
	}
	return domain.LoadBand(row.ID, row.Name), nil
}

type sqlxBand struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
