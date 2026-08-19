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

func (s *authorizedSetlistService) Create(ctx context.Context, params service.CreateSetlistParams) (uuid.UUID, error) {
	requesterID, err := requesterIDFromContext(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}
	if err := requireMember(ctx, s.executor, params.BandID, requesterID); err != nil {
		return uuid.UUID{}, err
	}
	return s.next.Create(ctx, params)
}

func (s *authorizedSetlistService) Update(ctx context.Context, params service.UpdateSetlistParams) error {
	requesterID, err := requesterIDFromContext(ctx)
	if err != nil {
		return err
	}
	if err := requireMember(ctx, s.executor, params.BandID, requesterID); err != nil {
		return err
	}
	return s.next.Update(ctx, params)
}

func (s *authorizedSetlistService) Remove(ctx context.Context, bandID, setlistID uuid.UUID) error {
	requesterID, err := requesterIDFromContext(ctx)
	if err != nil {
		return err
	}
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return err
	}
	return s.next.Remove(ctx, bandID, setlistID)
}

func (s *authorizedSetlistService) AddTracks(ctx context.Context, bandID, setlistID uuid.UUID, trackIDs []uuid.UUID) error {
	requesterID, err := requesterIDFromContext(ctx)
	if err != nil {
		return err
	}
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return err
	}
	return s.next.AddTracks(ctx, bandID, setlistID, trackIDs)
}

func (s *authorizedSetlistService) RemoveTracks(ctx context.Context, bandID, setlistID uuid.UUID, trackIDs []uuid.UUID) error {
	requesterID, err := requesterIDFromContext(ctx)
	if err != nil {
		return err
	}
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return err
	}
	return s.next.RemoveTracks(ctx, bandID, setlistID, trackIDs)
}

func (s *authorizedSetlistService) ReorderTracks(ctx context.Context, bandID, setlistID uuid.UUID, trackIDs []uuid.UUID) error {
	requesterID, err := requesterIDFromContext(ctx)
	if err != nil {
		return err
	}
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return err
	}
	return s.next.ReorderTracks(ctx, bandID, setlistID, trackIDs)
}
