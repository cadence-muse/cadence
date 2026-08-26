package service

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/common/uuid"
)

func newFakeSessionStore() *fakeSessionStore {
	return &fakeSessionStore{sessions: make(map[string]uuid.UUID)}
}

type fakeSessionStore struct {
	sessions map[string]uuid.UUID
}

func (f *fakeSessionStore) CreateSession(_ context.Context, userID uuid.UUID) (string, error) {
	token := uuid.Generate().String()
	f.sessions[token] = userID
	return token, nil
}

func (f *fakeSessionStore) ValidateSession(_ context.Context, token string) (uuid.UUID, error) {
	userID, ok := f.sessions[token]
	if !ok {
		return uuid.UUID{}, app.ErrSessionNotFound
	}
	return userID, nil
}

func (f *fakeSessionStore) DeleteSession(_ context.Context, token string) error {
	delete(f.sessions, token)
	return nil
}

func (f *fakeSessionStore) DeleteAllSessions(_ context.Context, userID uuid.UUID) error {
	for token, id := range f.sessions {
		if id == userID {
			delete(f.sessions, token)
		}
	}
	return nil
}
