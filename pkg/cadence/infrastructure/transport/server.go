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
	"cadence/pkg/cadence/app/query"
	"cadence/pkg/cadence/app/service"
	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/auth"
	"cadence/pkg/common/maybe"
	commonogenerrors "cadence/pkg/common/ogenerrors"
	"cadence/pkg/common/slices"
	"cadence/pkg/common/uuid"
)

func NewAPIServer(
	errorHandler ogenerrors.ErrorHandler,
	middlewares []middleware.Middleware,
	userService *service.UserService,
	bandService *service.BandService,
	bandQueryService query.BandQueryService,
	sessionStore app.SessionStore,
) (http.Handler, error) {
	apiHandler := newRESTHandler(
		userService,
		bandService,
		bandQueryService,
		sessionStore,
	)
	return publicapi.NewServer(
		apiHandler,
		publicapi.NewAuthHandler(sessionStore),
		publicapi.WithErrorHandler(errorHandler),
		publicapi.WithMiddleware(middlewares...),
	)
}

func newRESTHandler(
	userService *service.UserService,
	bandService *service.BandService,
	bandQueryService query.BandQueryService,
	sessionStore app.SessionStore,
) publicapi.Handler {
	return &restHandler{
		userService:      userService,
		bandService:      bandService,
		bandQueryService: bandQueryService,
		sessionStore:     sessionStore,
	}
}

type restHandler struct {
	userService      *service.UserService
	bandService      *service.BandService
	bandQueryService query.BandQueryService
	sessionStore     app.SessionStore
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

func (h *restHandler) ListBands(ctx context.Context) (publicapi.ListBandsRes, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, commonogenerrors.NewPermissionDeniedError("user not authenticated")
	}
	bands, err := h.bandQueryService.ListUserBands(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &publicapi.ListBandsResponseBody{Items: slices.Map(bands, convertQueryBandListItemToAPI)}, nil
}

func (h *restHandler) CreateBand(ctx context.Context, req *publicapi.CreateBandRequestBody) (publicapi.CreateBandRes, error) {
	id, err := h.bandService.Create(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	return &publicapi.CreateBandResponseBody{ID: googleuuid.UUID(id)}, nil
}

func (h *restHandler) GetBand(ctx context.Context, params publicapi.GetBandParams) (publicapi.GetBandRes, error) {
	band, err := h.bandQueryService.FindBand(ctx, uuid.UUID(params.BandId))
	if err != nil {
		return nil, err
	}
	foundBand, ok := maybe.JustValid(band)
	if !ok {
		return nil, commonogenerrors.NewNotFoundError("band not found")
	}
	return new(convertQueryBandDataToAPI(foundBand)), nil
}

func (h *restHandler) UpdateBand(_ context.Context, _ *publicapi.UpdateBandRequestBody, _ publicapi.UpdateBandParams) (publicapi.UpdateBandRes, error) {
	panic("implement me")
}

func (h *restHandler) CreateBandTrack(_ context.Context, _ *publicapi.CreateBandTrackRequestBody, _ publicapi.CreateBandTrackParams) (publicapi.CreateBandTrackRes, error) {
	panic("implement me")
}

func (h *restHandler) GetBandTrack(_ context.Context, _ *publicapi.BandTrack, _ publicapi.GetBandTrackParams) (publicapi.GetBandTrackRes, error) {
	panic("implement me")
}

func (h *restHandler) ListBandTracks(_ context.Context, _ publicapi.ListBandTracksParams) (publicapi.ListBandTracksRes, error) {
	return &publicapi.TrackList{}, nil
}

func (h *restHandler) UpdateBandTrack(_ context.Context, _ *publicapi.UpdateBandTrackRequestBody, _ publicapi.UpdateBandTrackParams) (publicapi.UpdateBandTrackRes, error) {
	panic("implement me")
}
