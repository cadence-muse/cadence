package query

import (
	"context"

	"cadence/pkg/common/maybe"
	"cadence/pkg/common/uuid"
)

type BandQueryService interface {
	ListUserBands(ctx context.Context, userID uuid.UUID) ([]BandListItem, error)
	FindBand(ctx context.Context, id uuid.UUID) (maybe.Maybe[BandData], error)
}

type BandListItem struct {
	ID   uuid.UUID
	Name string
}

type BandData struct {
	ID         uuid.UUID
	Name       string
	InviteCode string
}
