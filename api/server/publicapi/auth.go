package publicapi

import (
	"context"
)

func NewAuthHandler() SecurityHandler {
	return &authHandler{}
}

type authHandler struct{}

func (h *authHandler) HandleCookieAuth(ctx context.Context, _ OperationName, _ CookieAuth) (context.Context, error) {
	return ctx, nil
}

func (s *CookieAuth) CookieAuth(_ context.Context, _ OperationName) (CookieAuth, error) {
	return *s, nil
}
