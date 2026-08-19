package app

import (
	"context"
	"errors"

	"cadence/pkg/common/uuid"
)

var ErrSessionNotFound = errors.New("session not found")

type SessionStore interface {
	CreateSession(ctx context.Context, userID uuid.UUID) (token string, err error)
	ValidateSession(ctx context.Context, token string) (userID uuid.UUID, err error)
	DeleteSession(ctx context.Context, token string) error
}
