package query

import (
	"context"

	"github.com/nightnoryu/go-kita/maybe"

	"cadence/pkg/common/uuid"
)

type UserQueryService interface {
	FindUser(ctx context.Context, id uuid.UUID) (maybe.Maybe[UserData], error)
}

type UserData struct {
	ID       uuid.UUID
	Username string
}
