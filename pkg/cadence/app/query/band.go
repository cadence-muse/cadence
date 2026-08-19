package query

import (
	"context"

	"cadence/pkg/common/maybe"
	"cadence/pkg/common/uuid"
)

type BandQueryService interface {
	FindBand(ctx context.Context, id uuid.UUID) (maybe.Maybe[BandData], error)

	ListUserBands(ctx context.Context, userID uuid.UUID) ([]BandListItem, error)
	CountUserBands(ctx context.Context, userID uuid.UUID) (int, error)
}

type BandListItem struct {
	ID   uuid.UUID
	Name string
}

type BandData struct {
	ID         uuid.UUID
	Name       string
	OwnerID    uuid.UUID
	InviteCode string
	Members    []BandMemberData
}

type BandMemberData struct {
	ID       uuid.UUID
	Username string
	Role     BandMemberRole
}

type BandMemberRole string

const (
	BandMemberRoleOwner  BandMemberRole = "owner"
	BandMemberRoleMember BandMemberRole = "member"
)
