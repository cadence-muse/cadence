package publicapi

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/common/auth"
)

func NewAuthHandler(sessionStore app.SessionStore) SecurityHandler {
	return &authHandler{sessionStore: sessionStore}
}

type authHandler struct {
	sessionStore app.SessionStore
}

func (h *authHandler) HandleSessionAuth(ctx context.Context, _ OperationName, t SessionAuth) (context.Context, error) {
	userID, err := h.sessionStore.ValidateSession(ctx, t.APIKey)
	if err != nil {
		return nil, err
	}
	ctx = auth.InjectUserID(ctx, userID)
	ctx = auth.InjectSessionToken(ctx, t.APIKey)
	return ctx, nil
}

func (s *SessionAuth) SessionAuth(_ context.Context, _ OperationName) (SessionAuth, error) {
	return *s, nil
}
