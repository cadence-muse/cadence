package main

import (
	"context"
	"errors"

	"cadence/pkg/common/log"
)

var errMigrationFinished = errors.New("migration finished without errors")

func migrate(_ context.Context, _ *config, _ log.Logger) error {
	return nil
}
