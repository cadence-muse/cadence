package service

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/common/uuid"
)

func NewSessionService(store app.SessionStore) *SessionService {
	return &SessionService{store: store}
}

type SessionService struct {
	store app.SessionStore
}

func (s *SessionService) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	return s.store.CreateSession(ctx, userID)
}

func (s *SessionService) ValidateSession(ctx context.Context, token string) (uuid.UUID, error) {
	return s.store.ValidateSession(ctx, token)
}
