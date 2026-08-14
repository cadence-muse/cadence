package transport

import (
	"context"
	"errors"
	"net/http"
	"time"

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
	userQueryService query.UserQueryService,
	bandService *service.BandService,
	bandQueryService query.BandQueryService,
	trackService *service.TrackService,
	trackQueryService query.TrackQueryService,
	setlistService *service.SetlistService,
	setlistQueryService query.SetlistQueryService,
	sessionStore app.SessionStore,
) (http.Handler, error) {
	apiHandler := newRESTHandler(
		userService,
		userQueryService,
		bandService,
		bandQueryService,
		trackService,
		trackQueryService,
		setlistService,
		setlistQueryService,
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
	userQueryService query.UserQueryService,
	bandService *service.BandService,
	bandQueryService query.BandQueryService,
	trackService *service.TrackService,
	trackQueryService query.TrackQueryService,
	setlistService *service.SetlistService,
	setlistQueryService query.SetlistQueryService,
	sessionStore app.SessionStore,
) publicapi.Handler {
	return &restHandler{
		userService:         userService,
		userQueryService:    userQueryService,
		bandService:         bandService,
		bandQueryService:    bandQueryService,
		trackService:        trackService,
		trackQueryService:   trackQueryService,
		setlistService:      setlistService,
		setlistQueryService: setlistQueryService,
		sessionStore:        sessionStore,
	}
}

type restHandler struct {
	userService         *service.UserService
	userQueryService    query.UserQueryService
	bandService         *service.BandService
	bandQueryService    query.BandQueryService
	trackService        *service.TrackService
	trackQueryService   query.TrackQueryService
	setlistService      *service.SetlistService
	setlistQueryService query.SetlistQueryService
	sessionStore        app.SessionStore
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
			return nil, commonogenerrors.NewPermissionDeniedError("invalid credentials")
		}
		return nil, err
	}

	token, err := h.sessionStore.CreateSession(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &publicapi.LoginResponseBody{Token: token}, nil
}

func (h *restHandler) GetUserProfile(ctx context.Context) (publicapi.GetUserProfileRes, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, commonogenerrors.NewPermissionDeniedError("user not authenticated")
	}
	user, err := h.userQueryService.FindUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	foundUser, ok := maybe.JustValid(user)
	if !ok {
		return nil, commonogenerrors.NewNotFoundError("user not found")
	}
	return &publicapi.UserProfile{ID: googleuuid.UUID(foundUser.ID), Username: foundUser.Username}, nil
}

func (h *restHandler) GetHomepageData(ctx context.Context) (publicapi.GetHomepageDataRes, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, commonogenerrors.NewPermissionDeniedError("user not authenticated")
	}
	user, err := h.userQueryService.FindUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	foundUser, ok := maybe.JustValid(user)
	if !ok {
		return nil, commonogenerrors.NewNotFoundError("user not found")
	}

	bandsCount, err := h.bandQueryService.CountUserBands(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &publicapi.HomepageData{
		Username:   foundUser.Username,
		BandsCount: bandsCount,
	}, nil
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
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, commonogenerrors.NewPermissionDeniedError("user not authenticated")
	}
	id, err := h.bandService.Create(ctx, service.CreateBandParams{
		OwnerID: userID,
		Name:    req.GetName(),
	})
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

func (h *restHandler) UpdateBand(ctx context.Context, req *publicapi.UpdateBandRequestBody, params publicapi.UpdateBandParams) (publicapi.UpdateBandRes, error) {
	err := h.bandService.Update(ctx, service.UpdateBandParams{
		BandID: uuid.UUID(params.BandId),
		Name:   maybeValueFromOpt[string](req.GetName()),
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBandNotFound):
			return nil, commonogenerrors.NewNotFoundError(err.Error())
		case errors.Is(err, domain.ErrEmptyBandName),
			errors.Is(err, domain.ErrBandNameTooLong):
			return nil, commonogenerrors.NewInvalidInputError(err.Error())
		default:
			return nil, err
		}
	}
	return &publicapi.UpdateBandOK{}, nil
}

func (h *restHandler) RemoveBand(ctx context.Context, params publicapi.RemoveBandParams) (publicapi.RemoveBandRes, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, commonogenerrors.NewPermissionDeniedError("user not authenticated")
	}
	err := h.bandService.Remove(ctx, uuid.UUID(params.BandId), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBandNotFound):
			return nil, commonogenerrors.NewNotFoundError(err.Error())
		case errors.Is(err, domain.ErrNotBandOwner):
			return nil, commonogenerrors.NewPermissionDeniedError(err.Error())
		default:
			return nil, err
		}
	}
	return &publicapi.RemoveBandNoContent{}, nil
}

func (h *restHandler) JoinBand(ctx context.Context, req *publicapi.JoinBandRequestBody) (publicapi.JoinBandRes, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, commonogenerrors.NewPermissionDeniedError("user not authenticated")
	}
	if err := h.bandService.JoinByInviteCode(ctx, userID, req.GetInviteCode()); err != nil {
		if errors.Is(err, domain.ErrBandNotFound) {
			return nil, commonogenerrors.NewNotFoundError(err.Error())
		}
		return nil, err
	}
	return &publicapi.JoinBandOK{}, nil
}

func (h *restHandler) RemoveBandMember(ctx context.Context, params publicapi.RemoveBandMemberParams) (publicapi.RemoveBandMemberRes, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, commonogenerrors.NewPermissionDeniedError("user not authenticated")
	}
	err := h.bandService.RemoveMember(ctx, uuid.UUID(params.BandId), uuid.UUID(params.UserId), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBandNotFound), errors.Is(err, domain.ErrBandMemberNotFound):
			return nil, commonogenerrors.NewNotFoundError(err.Error())
		case errors.Is(err, domain.ErrNotBandOwner):
			return nil, commonogenerrors.NewPermissionDeniedError(err.Error())
		case errors.Is(err, domain.ErrCannotRemoveOwner):
			return nil, commonogenerrors.NewInvalidInputError(err.Error())
		default:
			return nil, err
		}
	}
	return &publicapi.RemoveBandMemberNoContent{}, nil
}

func (h *restHandler) CreateBandTrack(ctx context.Context, req *publicapi.CreateBandTrackRequestBody, params publicapi.CreateBandTrackParams) (publicapi.CreateBandTrackRes, error) {
	key, err := maybeKeyFromOpt(req.GetKey())
	if err != nil {
		return nil, commonogenerrors.NewInvalidInputError(err.Error())
	}

	trackID, err := h.trackService.Create(ctx, service.CreateTrackParams{
		BandID:   uuid.UUID(params.BandId),
		Title:    req.GetTitle(),
		Artist:   req.GetArtist(),
		Duration: maybeDurationFromOpt(req.GetDurationSeconds()),
		Tempo:    maybeValueFromOpt[int](req.GetTempo()),
		Key:      key,
		Notes:    maybeValueFromOpt[string](req.GetNotes()),
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmptyTrackTitle),
			errors.Is(err, domain.ErrTrackTitleTooLong),
			errors.Is(err, domain.ErrEmptyTrackArtist),
			errors.Is(err, domain.ErrTrackArtistTooLong):
			return nil, commonogenerrors.NewInvalidInputError(err.Error())
		default:
			return nil, err
		}
	}
	return &publicapi.CreateBandTrackResponseBody{ID: googleuuid.UUID(trackID)}, nil
}

func (h *restHandler) GetBandTrack(ctx context.Context, params publicapi.GetBandTrackParams) (publicapi.GetBandTrackRes, error) {
	track, err := h.trackQueryService.FindTrack(ctx, uuid.UUID(params.BandId), uuid.UUID(params.TrackId))
	if err != nil {
		return nil, err
	}
	foundTrack, ok := maybe.JustValid(track)
	if !ok {
		return nil, commonogenerrors.NewNotFoundError("track not found")
	}
	return new(convertQueryTrackDataToAPI(foundTrack)), nil
}

func (h *restHandler) ListBandTracks(ctx context.Context, params publicapi.ListBandTracksParams) (publicapi.ListBandTracksRes, error) {
	tracks, err := h.trackQueryService.ListBandTracks(ctx, uuid.UUID(params.BandId))
	if err != nil {
		return nil, err
	}
	return &publicapi.ListBandTracksResponseBody{Items: slices.Map(tracks, convertQueryTrackListItemToAPI)}, nil
}

func (h *restHandler) UpdateBandTrack(ctx context.Context, req *publicapi.UpdateBandTrackRequestBody, params publicapi.UpdateBandTrackParams) (publicapi.UpdateBandTrackRes, error) {
	key, err := maybeKeyFromOptNil(req.GetKey())
	if err != nil {
		return nil, commonogenerrors.NewInvalidInputError(err.Error())
	}

	err = h.trackService.Update(ctx, service.UpdateTrackParams{
		BandID:   uuid.UUID(params.BandId),
		TrackID:  uuid.UUID(params.TrackId),
		Title:    maybeValueFromOpt[string](req.GetTitle()),
		Artist:   maybeValueFromOpt[string](req.GetArtist()),
		Duration: maybeDurationFromOptNil(req.GetDurationSeconds()),
		Tempo:    maybeValueFromOptNil[int](req.GetTempo()),
		Key:      key,
		Notes:    maybeValueFromOptNil[string](req.GetNotes()),
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTrackNotFound):
			return nil, commonogenerrors.NewNotFoundError(err.Error())
		case errors.Is(err, domain.ErrEmptyTrackTitle),
			errors.Is(err, domain.ErrTrackTitleTooLong),
			errors.Is(err, domain.ErrEmptyTrackArtist),
			errors.Is(err, domain.ErrTrackArtistTooLong):
			return nil, commonogenerrors.NewInvalidInputError(err.Error())
		default:
			return nil, err
		}
	}
	return &publicapi.UpdateBandTrackOK{}, nil
}

func (h *restHandler) RemoveBandTrack(ctx context.Context, params publicapi.RemoveBandTrackParams) (publicapi.RemoveBandTrackRes, error) {
	err := h.trackService.Remove(ctx, uuid.UUID(params.BandId), uuid.UUID(params.TrackId))
	if err != nil {
		if errors.Is(err, domain.ErrTrackNotFound) {
			return nil, commonogenerrors.NewNotFoundError(err.Error())
		}
		return nil, err
	}
	return &publicapi.RemoveBandTrackNoContent{}, nil
}

func (h *restHandler) ListBandSetlists(ctx context.Context, params publicapi.ListBandSetlistsParams) (publicapi.ListBandSetlistsRes, error) {
	setlists, err := h.setlistQueryService.ListBandSetlists(ctx, uuid.UUID(params.BandId))
	if err != nil {
		return nil, err
	}
	return &publicapi.ListBandSetlistsResponseBody{Items: slices.Map(setlists, convertQuerySetlistListItemToAPI)}, nil
}

func (h *restHandler) CreateBandSetlist(ctx context.Context, req *publicapi.CreateBandSetlistsRequestBody, params publicapi.CreateBandSetlistParams) (publicapi.CreateBandSetlistRes, error) {
	setlistID, err := h.setlistService.Create(ctx, service.CreateSetlistParams{
		BandID:        uuid.UUID(params.BandId),
		Name:          req.GetName(),
		EventLocation: maybeValueFromOpt[string](req.GetEventLocation()),
		EventDate:     maybeValueFromOpt[time.Time](req.GetEventDate()),
		TrackIDs:      slices.Map(req.GetTrackIds(), func(id googleuuid.UUID) uuid.UUID { return uuid.UUID(id) }),
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBandNotFound), errors.Is(err, domain.ErrTrackNotFound):
			return nil, commonogenerrors.NewNotFoundError(err.Error())
		case errors.Is(err, domain.ErrEmptySetlistName),
			errors.Is(err, domain.ErrSetlistNameTooLong),
			errors.Is(err, domain.ErrSetlistEventLocationTooLong),
			errors.Is(err, domain.ErrDuplicateTrackIDs),
			errors.Is(err, domain.ErrTooManySetlistTracks):
			return nil, commonogenerrors.NewInvalidInputError(err.Error())
		default:
			return nil, err
		}
	}
	return &publicapi.CreateBandSetlistsResponseBody{ID: googleuuid.UUID(setlistID)}, nil
}

func (h *restHandler) GetBandSetlist(ctx context.Context, params publicapi.GetBandSetlistParams) (publicapi.GetBandSetlistRes, error) {
	setlist, err := h.setlistQueryService.FindSetlist(ctx, uuid.UUID(params.BandId), uuid.UUID(params.SetlistId))
	if err != nil {
		return nil, err
	}
	foundSetlist, ok := maybe.JustValid(setlist)
	if !ok {
		return nil, commonogenerrors.NewNotFoundError("setlist not found")
	}
	return new(convertQuerySetlistDataToAPI(foundSetlist)), nil
}

func (h *restHandler) UpdateBandSetlist(ctx context.Context, req *publicapi.UpdateBandSetlistRequestBody, params publicapi.UpdateBandSetlistParams) (publicapi.UpdateBandSetlistRes, error) {
	err := h.setlistService.Update(ctx, service.UpdateSetlistParams{
		BandID:        uuid.UUID(params.BandId),
		SetlistID:     uuid.UUID(params.SetlistId),
		Name:          maybeValueFromOpt[string](req.GetName()),
		EventLocation: maybeValueFromOptNil[string](req.GetEventLocation()),
		EventDate:     maybeValueFromOptNil[time.Time](req.GetEventDate()),
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSetlistNotFound):
			return nil, commonogenerrors.NewNotFoundError(err.Error())
		case errors.Is(err, domain.ErrEmptySetlistName),
			errors.Is(err, domain.ErrSetlistNameTooLong),
			errors.Is(err, domain.ErrSetlistEventLocationTooLong):
			return nil, commonogenerrors.NewInvalidInputError(err.Error())
		default:
			return nil, err
		}
	}
	return &publicapi.UpdateBandSetlistOK{}, nil
}

func (h *restHandler) RemoveBandSetlist(ctx context.Context, params publicapi.RemoveBandSetlistParams) (publicapi.RemoveBandSetlistRes, error) {
	err := h.setlistService.Remove(ctx, uuid.UUID(params.BandId), uuid.UUID(params.SetlistId))
	if err != nil {
		if errors.Is(err, domain.ErrSetlistNotFound) {
			return nil, commonogenerrors.NewNotFoundError(err.Error())
		}
		return nil, err
	}
	return &publicapi.RemoveBandSetlistNoContent{}, nil
}

func (h *restHandler) AddSetlistTrack(ctx context.Context, req *publicapi.AddSetlistTrackRequestBody, params publicapi.AddSetlistTrackParams) (publicapi.AddSetlistTrackRes, error) {
	err := h.setlistService.AddTrack(ctx, uuid.UUID(params.BandId), uuid.UUID(params.SetlistId), uuid.UUID(req.GetTrackId()))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSetlistNotFound), errors.Is(err, domain.ErrTrackNotFound):
			return nil, commonogenerrors.NewNotFoundError(err.Error())
		case errors.Is(err, domain.ErrTrackAlreadyInSetlist):
			return nil, commonogenerrors.NewAlreadyExistsError(err.Error())
		case errors.Is(err, domain.ErrTooManySetlistTracks):
			return nil, commonogenerrors.NewInvalidInputError(err.Error())
		default:
			return nil, err
		}
	}
	return &publicapi.AddSetlistTrackNoContent{}, nil
}

func (h *restHandler) RemoveSetlistTrack(ctx context.Context, params publicapi.RemoveSetlistTrackParams) (publicapi.RemoveSetlistTrackRes, error) {
	err := h.setlistService.RemoveTrack(ctx, uuid.UUID(params.BandId), uuid.UUID(params.SetlistId), uuid.UUID(params.TrackId))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSetlistNotFound), errors.Is(err, domain.ErrTrackNotInSetlist):
			return nil, commonogenerrors.NewNotFoundError(err.Error())
		default:
			return nil, err
		}
	}
	return &publicapi.RemoveSetlistTrackNoContent{}, nil
}

func (h *restHandler) ReorderSetlistTracks(ctx context.Context, req *publicapi.ReorderSetlistTracksRequestBody, params publicapi.ReorderSetlistTracksParams) (publicapi.ReorderSetlistTracksRes, error) {
	trackIDs := slices.Map(req.GetTrackIds(), func(id googleuuid.UUID) uuid.UUID { return uuid.UUID(id) })
	err := h.setlistService.ReorderTracks(ctx, uuid.UUID(params.BandId), uuid.UUID(params.SetlistId), trackIDs)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSetlistNotFound):
			return nil, commonogenerrors.NewNotFoundError(err.Error())
		case errors.Is(err, domain.ErrInvalidSetlistTrackOrder):
			return nil, commonogenerrors.NewInvalidInputError(err.Error())
		default:
			return nil, err
		}
	}
	return &publicapi.ReorderSetlistTracksNoContent{}, nil
}
