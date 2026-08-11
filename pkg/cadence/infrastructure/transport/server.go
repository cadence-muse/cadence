package transport

import (
	"context"
	"net/http"

	"cadence/api/server/publicapi"

	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
)

func NewAPIServer(
	errorHandler ogenerrors.ErrorHandler,
	middlewares []middleware.Middleware,
) (http.Handler, error) {
	apiHandler := newRESTHandler()
	return publicapi.NewServer(
		apiHandler,
		publicapi.NewAuthHandler(),
		publicapi.WithErrorHandler(errorHandler),
		publicapi.WithMiddleware(middlewares...),
	)
}

func newRESTHandler() publicapi.Handler {
	return &restHandler{}
}

type restHandler struct{}

func (h *restHandler) ListBands(_ context.Context) (publicapi.ListBandsRes, error) {
	return &publicapi.ListBandsOKApplicationJSON{}, nil
}
