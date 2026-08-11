package main

import (
	"github.com/go-faster/errors"
	"github.com/gorilla/mux"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"

	"cadence/pkg/cadence/infrastructure/transport"
	"cadence/pkg/common/log"
	"cadence/pkg/common/ogenmiddleware"
)

type dependencyContainer struct{}

func newDependencyContainer(
	config *config,
	logger log.Logger,
	router *mux.Router,
	errorHandler ogenerrors.ErrorHandler,
) (*dependencyContainer, error) {
	_, err := newDatabaseMigrator(config, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to migrate")
	}

	middlewares := []middleware.Middleware{
		ogenmiddleware.NewLoggingMiddleware(logger),
	}

	apiServer, err := transport.NewAPIServer(errorHandler, middlewares)
	if err != nil {
		return nil, err
	}
	router.PathPrefix("/api").Handler(apiServer)

	return &dependencyContainer{}, nil
}
