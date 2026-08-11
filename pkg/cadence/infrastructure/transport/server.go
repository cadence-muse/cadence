package transport

import (
	"context"
	"net/http"

	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"

	"cadence/api/server/publicapi"
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

func (h *restHandler) Login(_ context.Context, _ *publicapi.LoginRequestBody) (publicapi.LoginRes, error) {
	// TODO implement me
	panic("implement me")
}

func (h *restHandler) UpdateBand(_ context.Context, _ *publicapi.UpdateBandRequestBody, _ publicapi.UpdateBandParams) (publicapi.UpdateBandRes, error) {
	// TODO implement me
	panic("implement me")
}

func (h *restHandler) UpdateBandTrack(_ context.Context, _ *publicapi.UpdateBandTrackRequestBody, _ publicapi.UpdateBandTrackParams) (publicapi.UpdateBandTrackRes, error) {
	// TODO implement me
	panic("implement me")
}
