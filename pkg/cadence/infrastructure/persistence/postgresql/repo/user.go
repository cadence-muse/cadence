package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/nightnoryu/go-kita/postgresql"

	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/uuid"
)

func NewUserRepository(
	ctx context.Context,
	client postgresql.ClientContext,
) domain.UserRepository {
	return &userRepository{
		ctx:    ctx,
		client: client,
	}
}

type userRepository struct {
	ctx    context.Context
	client postgresql.ClientContext
}

func (repo *userRepository) NextID() domain.UserID {
	return uuid.Generate()
}

func (repo *userRepository) Store(user *domain.User) error {
	const sqlQuery = `
		INSERT INTO "user" (id, username, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET username      = EXCLUDED.username,
		    password_hash = EXCLUDED.password_hash,
		    updated_at    = EXCLUDED.created_at
	`
	_, err := repo.client.ExecContext(
		repo.ctx,
		sqlQuery,
		user.ID(),
		user.Username(),
		user.PasswordHash(),
		time.Now(),
	)
	return err
}

func (repo *userRepository) Get(id domain.UserID) (*domain.User, error) {
	const sqlQuery = `SELECT id, username, password_hash FROM "user" WHERE id = $1`
	var row sqlxUser
	err := repo.client.GetContext(repo.ctx, &row, sqlQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return domain.LoadUser(row.ID, row.Username, row.PasswordHash), nil
}

func (repo *userRepository) FindByUsername(username string) (*domain.User, error) {
	const sqlQuery = `SELECT id, username, password_hash FROM "user" WHERE username = $1`
	var row sqlxUser
	err := repo.client.GetContext(repo.ctx, &row, sqlQuery, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return domain.LoadUser(row.ID, row.Username, row.PasswordHash), nil
}

type sqlxUser struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
}
