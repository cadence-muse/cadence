package transport

import (
	"context"
	"errors"
	"net/http"

	googleuuid "github.com/google/uuid"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"

	"cadence/api/server/publicapi"
	"cadence/pkg/cadence/app/service"
	"cadence/pkg/cadence/domain"
	commonogenerrors "cadence/pkg/common/ogenerrors"
)

func NewAPIServer(
	errorHandler ogenerrors.ErrorHandler,
	middlewares []middleware.Middleware,
	userService *service.UserService,
) (http.Handler, error) {
	apiHandler := newRESTHandler(userService)
	return publicapi.NewServer(
		apiHandler,
		publicapi.NewAuthHandler(),
		publicapi.WithErrorHandler(errorHandler),
		publicapi.WithMiddleware(middlewares...),
	)
}

func newRESTHandler(userService *service.UserService) publicapi.Handler {
	return &restHandler{userService: userService}
}

type restHandler struct {
	userService *service.UserService
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

func (h *restHandler) Login(_ context.Context, _ *publicapi.LoginRequestBody) (publicapi.LoginRes, error) {
	// TODO implement me
	panic("implement me")
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
