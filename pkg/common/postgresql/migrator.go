package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cadence/pkg/common/log"

	"github.com/jmoiron/sqlx"
)

type Migrator interface {
	MigrateUp() error
}

type migrator struct {
	db     *sqlx.DB
	logger log.Logger
}

func (m migrator) MigrateUp() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := m.db.Connx(ctx)
	if err != nil {
		err = fmt.Errorf("failed to open migrator connection: %w", err)
		return err
	}

	defer func() {
		closeErr := conn.Close()
		if closeErr != nil && !errors.Is(closeErr, sql.ErrConnDone) {
			err = errors.Join(err, conn.Close())
		}
	}()

	// TODO implement
}
