package main

import (
	"context"
	"errors"

	"cadence/pkg/common/log"
)

var errMigrationFinished = errors.New("migration finished without errors")

func migrate(ctx context.Context, config *config, logger log.Logger) error {

}
