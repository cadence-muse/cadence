package authorized

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/app/query"
	"cadence/pkg/common/maybe"
	"cadence/pkg/common/transactional"
	"cadence/pkg/common/uuid"
)

func NewSetlistQueryService(inner query.SetlistQueryService, executor transactional.Executor[app.RepoProvider]) query.SetlistQueryService {
	return &authorizedSetlistQueryService{next: inner, executor: executor}
}

type authorizedSetlistQueryService struct {
	next     query.SetlistQueryService
	executor transactional.Executor[app.RepoProvider]
}

func (s *authorizedSetlistQueryService) ListBandSetlists(ctx context.Context, bandID, requesterID uuid.UUID) ([]query.SetlistListItem, error) {
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return nil, err
	}
	return s.next.ListBandSetlists(ctx, bandID, requesterID)
}

func (s *authorizedSetlistQueryService) ListUserSetlists(ctx context.Context, userID uuid.UUID, bandID maybe.Maybe[uuid.UUID]) ([]query.UserSetlistListItem, error) {
	return s.next.ListUserSetlists(ctx, userID, bandID)
}

func (s *authorizedSetlistQueryService) FindSetlist(ctx context.Context, bandID, setlistID, requesterID uuid.UUID) (maybe.Maybe[query.SetlistData], error) {
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return maybe.NewNone[query.SetlistData](), err
	}
	return s.next.FindSetlist(ctx, bandID, setlistID, requesterID)
}
