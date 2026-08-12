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

func (h *authHandler) HandleCookieAuth(ctx context.Context, _ OperationName, t CookieAuth) (context.Context, error) {
	userID, err := h.sessionStore.ValidateSession(ctx, t.APIKey)
	if err != nil {
		return nil, err
	}
	return auth.InjectUserID(ctx, userID), nil
}

func (s *CookieAuth) CookieAuth(_ context.Context, _ OperationName) (CookieAuth, error) {
	return *s, nil
}
