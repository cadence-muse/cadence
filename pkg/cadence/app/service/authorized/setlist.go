package authorized

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/app/service"
	"cadence/pkg/common/transactional"
	"cadence/pkg/common/uuid"
)

func NewSetlistService(inner service.SetlistService, executor transactional.Executor[app.RepoProvider]) service.SetlistService {
	return &authorizedSetlistService{next: inner, executor: executor}
}

type authorizedSetlistService struct {
	next     service.SetlistService
	executor transactional.Executor[app.RepoProvider]
}

func (s *authorizedSetlistService) Create(ctx context.Context, params service.CreateSetlistParams, requesterID uuid.UUID) (uuid.UUID, error) {
	if err := requireMember(ctx, s.executor, params.BandID, requesterID); err != nil {
		return uuid.UUID{}, err
	}
	return s.next.Create(ctx, params, requesterID)
}

func (s *authorizedSetlistService) Update(ctx context.Context, params service.UpdateSetlistParams, requesterID uuid.UUID) error {
	if err := requireMember(ctx, s.executor, params.BandID, requesterID); err != nil {
		return err
	}
	return s.next.Update(ctx, params, requesterID)
}

func (s *authorizedSetlistService) Remove(ctx context.Context, bandID, setlistID, requesterID uuid.UUID) error {
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return err
	}
	return s.next.Remove(ctx, bandID, setlistID, requesterID)
}

func (s *authorizedSetlistService) AddTrack(ctx context.Context, bandID, setlistID, trackID, requesterID uuid.UUID) error {
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return err
	}
	return s.next.AddTrack(ctx, bandID, setlistID, trackID, requesterID)
}

func (s *authorizedSetlistService) AddTracks(ctx context.Context, bandID, setlistID uuid.UUID, trackIDs []uuid.UUID, requesterID uuid.UUID) error {
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return err
	}
	return s.next.AddTracks(ctx, bandID, setlistID, trackIDs, requesterID)
}

func (s *authorizedSetlistService) RemoveTrack(ctx context.Context, bandID, setlistID, trackID, requesterID uuid.UUID) error {
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return err
	}
	return s.next.RemoveTrack(ctx, bandID, setlistID, trackID, requesterID)
}

func (s *authorizedSetlistService) ReorderTracks(ctx context.Context, bandID, setlistID uuid.UUID, trackIDs []uuid.UUID, requesterID uuid.UUID) error {
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return err
	}
	return s.next.ReorderTracks(ctx, bandID, setlistID, trackIDs, requesterID)
}
