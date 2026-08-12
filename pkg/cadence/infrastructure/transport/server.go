package transport

import (
	"context"
	"errors"
	"net/http"

	googleuuid "github.com/google/uuid"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"

	"cadence/api/server/publicapi"
	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/app/service"
	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/auth"
	commonogenerrors "cadence/pkg/common/ogenerrors"
)

func NewAPIServer(
	errorHandler ogenerrors.ErrorHandler,
	middlewares []middleware.Middleware,
	userService *service.UserService,
	sessionStore app.SessionStore,
) (http.Handler, error) {
	apiHandler := newRESTHandler(userService, sessionStore)
	return publicapi.NewServer(
		apiHandler,
		publicapi.NewAuthHandler(sessionStore),
		publicapi.WithErrorHandler(errorHandler),
		publicapi.WithMiddleware(middlewares...),
	)
}

func newRESTHandler(
	userService *service.UserService,
	sessionStore app.SessionStore,
) publicapi.Handler {
	return &restHandler{
		userService:  userService,
		sessionStore: sessionStore,
	}
}

type restHandler struct {
	userService  *service.UserService
	sessionStore app.SessionStore
}

func (h *restHandler) Register(ctx context.Context, req *publicapi.RegisterRequestBody) (publicapi.RegisterRes, error) {
	userID, err := h.userService.Register(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUsernameTaken):
			return nil, commonogenerrors.NewAlreadyExistsError(err.Error())
		case errors.Is(err, domain.ErrEmptyUsername),
			errors.Is(err, domain.ErrUsernameTooLong),
			errors.Is(err, domain.ErrEmptyPasswordHash):
			return nil, commonogenerrors.NewInvalidInputError(err.Error())
		default:
			return nil, err
		}
	}
	return &publicapi.RegisterResponseBody{ID: googleuuid.UUID(userID)}, nil
}

func (h *restHandler) Login(ctx context.Context, req *publicapi.LoginRequestBody) (publicapi.LoginRes, error) {
	userID, err := h.userService.Authenticate(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return &publicapi.LoginUnauthorized{}, nil
		}
		return nil, err
	}

	token, err := h.sessionStore.CreateSession(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &publicapi.LoginResponseBody{Token: token}, nil
}

func (h *restHandler) Logout(ctx context.Context) (publicapi.LogoutRes, error) {
	token, ok := auth.SessionTokenFromContext(ctx)
	if !ok {
		return nil, errors.New("session token missing from context")
	}
	if err := h.sessionStore.DeleteSession(ctx, token); err != nil {
		return nil, err
	}
	return &publicapi.LogoutOK{}, nil
}

func (h *restHandler) CreateBand(_ context.Context, _ *publicapi.CreateBandRequestBody) (publicapi.CreateBandRes, error) {
	// TODO implement me
	panic("implement me")
}

func (h *restHandler) CreateBandTrack(_ context.Context, _ *publicapi.CreateBandTrackRequestBody, _ publicapi.CreateBandTrackParams) (publicapi.CreateBandTrackRes, error) {
	// TODO implement me
	panic("implement me")
}

func (h *restHandler) GetBand(_ context.Context, _ publicapi.GetBandParams) (publicapi.GetBandRes, error) {
	// TODO implement me
	panic("implement me")
}

func (h *restHandler) GetBandTrack(_ context.Context, _ *publicapi.BandTrack, _ publicapi.GetBandTrackParams) (publicapi.GetBandTrackRes, error) {
	// TODO implement me
	panic("implement me")
}

func (h *restHandler) ListBandTracks(_ context.Context, _ publicapi.ListBandTracksParams) (publicapi.ListBandTracksRes, error) {
	return &publicapi.TrackList{}, nil
}

func (h *restHandler) ListBands(_ context.Context) (publicapi.ListBandsRes, error) {
	return &publicapi.BandList{}, nil
}

func (h *restHandler) UpdateBand(_ context.Context, _ *publicapi.UpdateBandRequestBody, _ publicapi.UpdateBandParams) (publicapi.UpdateBandRes, error) {
	// TODO implement me
	panic("implement me")
}

func (h *restHandler) UpdateBandTrack(_ context.Context, _ *publicapi.UpdateBandTrackRequestBody, _ publicapi.UpdateBandTrackParams) (publicapi.UpdateBandTrackRes, error) {
	// TODO implement me
	panic("implement me")
}
