package main

import (
	"github.com/go-faster/errors"
	"github.com/gorilla/mux"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"

	"cadence/pkg/cadence/app"
	appservice "cadence/pkg/cadence/app/service"
	"cadence/pkg/cadence/infrastructure/persistence/postgresql/repo"
	"cadence/pkg/cadence/infrastructure/transport"
	"cadence/pkg/common/log"
	"cadence/pkg/common/ogenmiddleware"
	"cadence/pkg/common/postgresql"
	"cadence/pkg/common/transactional"
)

type dependencyContainer struct {
	userService *appservice.UserService
}

func newDependencyContainer(
	config *config,
	logger log.Logger,
	router *mux.Router,
	errorHandler ogenerrors.ErrorHandler,
) (*dependencyContainer, error) {
	migrator, err := newDatabaseMigrator(config, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to migrate")
	}

	connectionProvider := postgresql.NewConnectionProvider(migrator.connector.TransactionalClient())
	transactionFactory := repo.NewTransactionFactory(connectionProvider)
	executor := transactional.NewExecutor[app.RepoProvider](transactionFactory)
	userService := appservice.NewUserService(executor)

	middlewares := []middleware.Middleware{
		ogenmiddleware.NewLoggingMiddleware(logger),
	}

	apiServer, err := transport.NewAPIServer(errorHandler, middlewares, userService)
	if err != nil {
		return nil, err
	}
	router.PathPrefix("/api").Handler(apiServer)

	return &dependencyContainer{userService: userService}, nil
}
