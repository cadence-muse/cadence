package query

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nightnoryu/go-kita/maybe"
	"github.com/nightnoryu/go-kita/postgresql"

	"cadence/pkg/cadence/app/query"
	"cadence/pkg/common/uuid"
)

func NewUserQueryService(client postgresql.ClientContext) query.UserQueryService {
	return &userQueryService{client: client}
}

type userQueryService struct {
	client postgresql.ClientContext
}

func (s *userQueryService) FindUser(ctx context.Context, id uuid.UUID) (maybe.Maybe[query.UserData], error) {
	const sqlQuery = `SELECT id, username FROM "user" WHERE id = $1`
	var row sqlxUserData
	err := s.client.GetContext(ctx, &row, sqlQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return maybe.NewNone[query.UserData](), nil
		}
		return maybe.NewNone[query.UserData](), err
	}

	return maybe.NewJust(query.UserData{ID: row.ID, Username: row.Username}), nil
}

type sqlxUserData struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`
}
