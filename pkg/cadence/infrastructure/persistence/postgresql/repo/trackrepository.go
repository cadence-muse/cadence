package repo

import (
	"context"

	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/maybe"
	"cadence/pkg/common/postgresql"
	"cadence/pkg/common/uuid"
)

func NewTrackRepository(
	ctx context.Context,
	client postgresql.ClientContext,
	subjectID maybe.Maybe[uuid.UUID],
) domain.TrackRepository {
	return &trackRepository{
		ctx:       ctx,
		client:    client,
		subjectID: subjectID,
	}
}

type trackRepository struct {
	ctx       context.Context
	client    postgresql.ClientContext
	subjectID maybe.Maybe[uuid.UUID]
}

func (repo *trackRepository) NextID() domain.TrackID {
	// TODO implement me
	panic("implement me")
}

func (repo *trackRepository) Store(_ *domain.Track) error {
	// TODO implement me
	panic("implement me")
}

func (repo *trackRepository) Get(_ domain.TrackID) (*domain.Track, error) {
	// TODO implement me
	panic("implement me")
}
