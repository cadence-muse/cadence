package authorized // nolint:dupl

import (
	"context"

	"github.com/nightnoryu/go-kita/maybe"
	"github.com/nightnoryu/go-kita/transactional"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/app/query"
	"cadence/pkg/common/uuid"
)

func NewSetlistQueryService(inner query.SetlistQueryService, executor transactional.Executor[app.RepoProvider]) query.SetlistQueryService {
	return &authorizedSetlistQueryService{next: inner, executor: executor}
}

type authorizedSetlistQueryService struct {
	next     query.SetlistQueryService
	executor transactional.Executor[app.RepoProvider]
}

func (s *authorizedSetlistQueryService) ListBandSetlists(ctx context.Context, bandID uuid.UUID) ([]query.SetlistListItem, error) {
	requesterID, err := requesterIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return nil, err
	}
	return s.next.ListBandSetlists(ctx, bandID)
}

func (s *authorizedSetlistQueryService) ListUserSetlists(ctx context.Context, userID uuid.UUID, bandID maybe.Maybe[uuid.UUID], searchQuery maybe.Maybe[string]) ([]query.UserSetlistListItem, error) {
	return s.next.ListUserSetlists(ctx, userID, bandID, searchQuery)
}

func (s *authorizedSetlistQueryService) FindSetlist(ctx context.Context, bandID, setlistID uuid.UUID) (maybe.Maybe[query.SetlistData], error) {
	requesterID, err := requesterIDFromContext(ctx)
	if err != nil {
		return maybe.NewNone[query.SetlistData](), err
	}
	if err := requireMember(ctx, s.executor, bandID, requesterID); err != nil {
		return maybe.NewNone[query.SetlistData](), err
	}
	return s.next.FindSetlist(ctx, bandID, setlistID)
}
