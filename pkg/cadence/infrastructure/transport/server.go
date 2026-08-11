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

func (h *restHandler) CreateBand(ctx context.Context, req *publicapi.CreateBandRequestBody) (publicapi.CreateBandRes, error) {
	//TODO implement me
	panic("implement me")
}

func (h *restHandler) CreateBandTrack(ctx context.Context, req *publicapi.CreateBandTrackRequestBody, params publicapi.CreateBandTrackParams) (publicapi.CreateBandTrackRes, error) {
	//TODO implement me
	panic("implement me")
}

func (h *restHandler) GetBand(ctx context.Context, params publicapi.GetBandParams) (publicapi.GetBandRes, error) {
	//TODO implement me
	panic("implement me")
}

func (h *restHandler) GetBandTrack(ctx context.Context, req *publicapi.BandTrack, params publicapi.GetBandTrackParams) (publicapi.GetBandTrackRes, error) {
	//TODO implement me
	panic("implement me")
}

func (h *restHandler) ListBandTracks(ctx context.Context, params publicapi.ListBandTracksParams) (publicapi.ListBandTracksRes, error) {
	//TODO implement me
	panic("implement me")
}

func (h *restHandler) ListBands(ctx context.Context) (publicapi.ListBandsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (h *restHandler) Login(ctx context.Context, req *publicapi.LoginRequestBody) (publicapi.LoginRes, error) {
	//TODO implement me
	panic("implement me")
}

func (h *restHandler) UpdateBand(ctx context.Context, req *publicapi.UpdateBandRequestBody, params publicapi.UpdateBandParams) (publicapi.UpdateBandRes, error) {
	//TODO implement me
	panic("implement me")
}

func (h *restHandler) UpdateBandTrack(ctx context.Context, req *publicapi.UpdateBandTrackRequestBody, params publicapi.UpdateBandTrackParams) (publicapi.UpdateBandTrackRes, error) {
	//TODO implement me
	panic("implement me")
}
