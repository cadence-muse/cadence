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
		WHERE bm.user_id = $1
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

func (s *bandQueryService) FindBand(ctx context.Context, id uuid.UUID) (maybe.Maybe[query.BandData], error) {
	const sqlQuery = `SELECT id, name, invite_code FROM band WHERE id = $1`
	var row sqlxBandData
	err := s.client.GetContext(ctx, &row, sqlQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return maybe.NewNone[query.BandData](), nil
		}
		return maybe.NewNone[query.BandData](), err
	}

	return maybe.NewJust(query.BandData{
		ID:         row.ID,
		Name:       row.Name,
		InviteCode: row.InviteCode,
	}), nil
}

type sqlxBandListItem struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type sqlxBandData struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	InviteCode string    `db:"invite_code"`
}
