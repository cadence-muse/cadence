package query

import (
	"context"
	"database/sql"
	"errors"

	"cadence/pkg/cadence/app/query"
	"cadence/pkg/common/maybe"
	"cadence/pkg/common/postgresql"
	"cadence/pkg/common/slices"
	"cadence/pkg/common/uuid"
)

func NewBandQueryService(client postgresql.ClientContext) query.BandQueryService {
	return &bandQueryService{client: client}
}

type bandQueryService struct {
	client postgresql.ClientContext
}

func (s *bandQueryService) ListUserBands(ctx context.Context, userID uuid.UUID) ([]query.BandListItem, error) {
	const sqlQuery = `
		SELECT b.id, b.name
		FROM band b
		JOIN band_member bm ON bm.band_id = b.id
		WHERE bm.user_id = $1 AND b.deleted_at IS NULL
		ORDER BY b.name
	`
	var rows []sqlxBandListItem
	if err := s.client.SelectContext(ctx, &rows, sqlQuery, userID); err != nil {
		return nil, err
	}

	return slices.Map(rows, func(row sqlxBandListItem) query.BandListItem {
		return query.BandListItem(row)
	}), nil
}

func (s *bandQueryService) CountUserBands(ctx context.Context, userID uuid.UUID) (int, error) {
	const sqlQuery = `
		SELECT COUNT(*)
		FROM band b
		JOIN band_member bm ON bm.band_id = b.id
		WHERE bm.user_id = $1 AND b.deleted_at IS NULL
	`
	var result int
	if err := s.client.GetContext(ctx, &result, sqlQuery, userID); err != nil {
		return 0, err
	}

	return result, nil
}

func (s *bandQueryService) FindBand(ctx context.Context, id uuid.UUID) (maybe.Maybe[query.BandData], error) {
	const sqlQuery = `
		SELECT b.id, b.name, bm.user_id AS owner_id, b.invite_code
		FROM band b
		JOIN band_member bm ON bm.band_id = b.id AND bm.role = 'owner'
		WHERE b.id = $1 AND b.deleted_at IS NULL
	`
	var row sqlxBandData
	err := s.client.GetContext(ctx, &row, sqlQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return maybe.NewNone[query.BandData](), nil
		}
		return maybe.NewNone[query.BandData](), err
	}

	members, err := s.listBandMembers(ctx, id)
	if err != nil {
		return maybe.NewNone[query.BandData](), err
	}

	return maybe.NewJust(query.BandData{
		ID:         row.ID,
		Name:       row.Name,
		OwnerID:    row.OwnerID,
		InviteCode: row.InviteCode,
		Members:    members,
	}), nil
}

func (s *bandQueryService) listBandMembers(ctx context.Context, bandID uuid.UUID) ([]query.BandMemberData, error) {
	const sqlQuery = `
		SELECT u.id, u.username, bm.role
		FROM band_member bm
		JOIN "user" u ON u.id = bm.user_id
		WHERE bm.band_id = $1
		ORDER BY u.username
	`
	var rows []sqlxBandMember
	if err := s.client.SelectContext(ctx, &rows, sqlQuery, bandID); err != nil {
		return nil, err
	}

	return slices.Map(rows, func(row sqlxBandMember) query.BandMemberData {
		return query.BandMemberData{
			ID:       row.ID,
			Username: row.Username,
			Role:     query.BandMemberRole(row.Role),
		}
	}), nil
}

type sqlxBandListItem struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type sqlxBandData struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	OwnerID    uuid.UUID `db:"owner_id"`
	InviteCode string    `db:"invite_code"`
}

type sqlxBandMember struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Role     string    `db:"role"`
}
