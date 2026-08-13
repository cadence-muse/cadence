package repo

import (
	"context"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/postgresql"
	"cadence/pkg/common/uuid"
)

func NewRepoProvider(
	ctx context.Context,
	conn postgresql.TransactionalConnection,
	tx postgresql.Transaction,
	subjectID uuid.UUID,
) app.RepoProvider {
	return &repoProvider{
		ctx:       ctx,
		conn:      conn,
		tx:        tx,
		subjectID: subjectID,
	}
}

type repoProvider struct {
	ctx       context.Context
	conn      postgresql.TransactionalConnection
	tx        postgresql.Transaction
	subjectID uuid.UUID
}

func (p *repoProvider) UserRepository() domain.UserRepository {
	return NewUserRepository(p.ctx, p.tx)
}

func (p *repoProvider) BandRepository() domain.BandRepository {
	return NewBandRepository(p.ctx, p.tx, p.subjectID)
}

func (p *repoProvider) TrackRepository() domain.TrackRepository {
	return NewTrackRepository(p.ctx, p.tx, p.subjectID)
}

func (p *repoProvider) Complete(err error) error {
	if err != nil {
		if rollbackErr := p.tx.Rollback(); rollbackErr != nil {
			err = rollbackErr
		}
	} else {
		err = p.tx.Commit()
	}

	closeErr := p.conn.Close()
	if err != nil {
		return err
	}
	return closeErr
}
