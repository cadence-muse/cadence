package publicapi

import (
	"context"

	"cadence/pkg/cadence/app/service"
	"cadence/pkg/common/auth"
)

func NewAuthHandler(sessionService *service.SessionService) SecurityHandler {
	return &authHandler{sessionService: sessionService}
}

type authHandler struct {
	sessionService *service.SessionService
}

func (h *authHandler) HandleCookieAuth(ctx context.Context, _ OperationName, t CookieAuth) (context.Context, error) {
	userID, err := h.sessionService.ValidateSession(ctx, t.APIKey)
	if err != nil {
		return nil, err
	}
	return auth.InjectUserID(ctx, userID), nil
}

func (s *CookieAuth) CookieAuth(_ context.Context, _ OperationName) (CookieAuth, error) {
	return *s, nil
}
