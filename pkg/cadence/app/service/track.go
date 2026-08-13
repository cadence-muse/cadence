package service

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/common/transactional"
	"cadence/pkg/common/uuid"
)

func NewTrackService(executor transactional.Executor[app.RepoProvider]) *TrackService {
	return &TrackService{executor: executor}
}

type TrackService struct {
	executor transactional.Executor[app.RepoProvider]
}

func (s *TrackService) Create(ctx context.Context) (trackID uuid.UUID, err error) {
	// TODO implement
	return [16]byte{}, nil
}
