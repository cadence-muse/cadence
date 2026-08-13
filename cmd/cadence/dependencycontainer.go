package main

import (
	"github.com/go-faster/errors"
	"github.com/gorilla/mux"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"

	"cadence/pkg/cadence/app"
	appservice "cadence/pkg/cadence/app/service"
	"cadence/pkg/cadence/infrastructure/persistence/postgresql/query"
	"cadence/pkg/cadence/infrastructure/persistence/postgresql/repo"
	redisinfra "cadence/pkg/cadence/infrastructure/persistence/redis"
	"cadence/pkg/cadence/infrastructure/transport"
	"cadence/pkg/common/log"
	"cadence/pkg/common/ogenmiddleware"
	"cadence/pkg/common/postgresql"
	"cadence/pkg/common/redis"
	"cadence/pkg/common/transactional"
)

type dependencyContainer struct {
	dbConnector postgresql.Connector
	redisClient redis.Client
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

	middlewares := []middleware.Middleware{
		ogenmiddleware.NewLoggingMiddleware(logger),
	}

	userService := appservice.NewUserService(executor)

	bandService := appservice.NewBandService(executor)
	bandQueryService := query.NewBandQueryService(migrator.connector.TransactionalClient())

	redisClient := redis.NewClient(config.redisConfig())
	sessionStore := redisinfra.NewSessionStore(redisClient, redisinfra.Config{
		TTL:                config.SessionTTL,
		MaxSessionsPerUser: config.SessionMaxPerUser,
	})

	apiServer, err := transport.NewAPIServer(
		errorHandler,
		middlewares,
		userService,
		bandService,
		bandQueryService,
		sessionStore,
	)
	if err != nil {
		return nil, err
	}
	router.PathPrefix("/api").Handler(corsMiddleware(config.CORSAllowedOrigins)(apiServer))

	return &dependencyContainer{
		dbConnector: migrator.connector,
		redisClient: redisClient,
	}, nil
}

func (c *dependencyContainer) Close() error {
	dbErr := c.dbConnector.Close()
	redisErr := c.redisClient.Close()
	return errors.Join(dbErr, redisErr)
}
