package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/postgresql"
	"cadence/pkg/common/slices"
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
		INSERT INTO band (id, name, invite_code, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET name       = EXCLUDED.name,
		    updated_at = EXCLUDED.created_at,
		    updated_by = EXCLUDED.created_by
	`
	_, err := repo.client.ExecContext(
		repo.ctx,
		sqlQuery,
		band.ID(),
		band.Name(),
		band.InviteCode(),
		time.Now(),
		repo.subjectID,
	)
	if err != nil {
		return err
	}

	return repo.storeMembers(band)
}

func (repo *bandRepository) storeMembers(band *domain.Band) error {
	const deleteSQLQuery = `DELETE FROM band_member WHERE band_id = $1`
	if _, err := repo.client.ExecContext(repo.ctx, deleteSQLQuery, band.ID()); err != nil {
		return err
	}

	const insertSQLQuery = `
		INSERT INTO band_member (band_id, user_id, role)
		VALUES ($1, $2, $3)
	`
	for _, member := range band.Members() {
		if _, err := repo.client.ExecContext(repo.ctx, insertSQLQuery, band.ID(), member.UserID(), member.Role()); err != nil {
			return err
		}
	}
	return nil
}

func (repo *bandRepository) Get(id domain.BandID) (*domain.Band, error) {
	const sqlQuery = `SELECT id, name, invite_code FROM band WHERE id = $1 AND deleted_at IS NULL`
	var row sqlxBand
	err := repo.client.GetContext(repo.ctx, &row, sqlQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBandNotFound
		}
		return nil, err
	}

	members, err := repo.getMembers(row.ID)
	if err != nil {
		return nil, err
	}
	return domain.LoadBand(row.ID, row.Name, row.InviteCode, members), nil
}

func (repo *bandRepository) GetByInviteCode(inviteCode string) (*domain.Band, error) {
	const sqlQuery = `SELECT id, name, invite_code FROM band WHERE invite_code = $1 AND deleted_at IS NULL`
	var row sqlxBand
	err := repo.client.GetContext(repo.ctx, &row, sqlQuery, inviteCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBandNotFound
		}
		return nil, err
	}

	members, err := repo.getMembers(row.ID)
	if err != nil {
		return nil, err
	}
	return domain.LoadBand(row.ID, row.Name, row.InviteCode, members), nil
}

func (repo *bandRepository) Remove(id domain.BandID) error {
	deletedAt := time.Now()

	const sqlQuery = `
		UPDATE band
		SET deleted_at = $2,
		    deleted_by = $3
		WHERE id = $1 AND deleted_at IS NULL
	`
	if _, err := repo.client.ExecContext(repo.ctx, sqlQuery, id, deletedAt, repo.subjectID); err != nil {
		return err
	}

	const cascadeTracksSQLQuery = `
		UPDATE track
		SET deleted_at = $2,
		    deleted_by = $3
		WHERE band_id = $1 AND deleted_at IS NULL
	`
	if _, err := repo.client.ExecContext(repo.ctx, cascadeTracksSQLQuery, id, deletedAt, repo.subjectID); err != nil {
		return err
	}

	const cascadeSetlistsSQLQuery = `
		UPDATE setlist
		SET deleted_at = $2,
		    deleted_by = $3
		WHERE band_id = $1 AND deleted_at IS NULL
	`
	_, err := repo.client.ExecContext(repo.ctx, cascadeSetlistsSQLQuery, id, deletedAt, repo.subjectID)
	return err
}

func (repo *bandRepository) getMembers(bandID domain.BandID) ([]domain.BandMember, error) {
	const sqlQuery = `SELECT user_id, role FROM band_member WHERE band_id = $1`
	var rows []sqlxBandMember
	if err := repo.client.SelectContext(repo.ctx, &rows, sqlQuery, bandID); err != nil {
		return nil, err
	}

	return slices.Map(rows, func(row sqlxBandMember) domain.BandMember {
		return domain.LoadBandMember(row.UserID, domain.BandMemberRole(row.Role))
	}), nil
}

type sqlxBand struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	InviteCode string    `db:"invite_code"`
}

type sqlxBandMember struct {
	UserID uuid.UUID `db:"user_id"`
	Role   string    `db:"role"`
}
