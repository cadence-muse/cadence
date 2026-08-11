package repo

import (
	"context"

	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/maybe"
	"cadence/pkg/common/postgresql"
	"cadence/pkg/common/uuid"
)

func NewBandRepository(
	ctx context.Context,
	client postgresql.ClientContext,
	subjectID maybe.Maybe[uuid.UUID],
) domain.BandRepository {
	return &bandRepository{
		ctx:       ctx,
		client:    client,
		subjectID: subjectID,
	}
}

type bandRepository struct {
	ctx       context.Context
	client    postgresql.ClientContext
	subjectID maybe.Maybe[uuid.UUID]
}

func (repo *bandRepository) NextID() domain.BandID {
	// TODO implement me
	panic("implement me")
}

func (repo *bandRepository) Store(_ *domain.Band) error {
	// TODO implement me
	panic("implement me")
}

func (repo *bandRepository) Get(_ domain.BandID) (*domain.Band, error) {
	// TODO implement me
	panic("implement me")
}
