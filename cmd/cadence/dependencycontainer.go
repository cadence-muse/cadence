package main

import (
	"context"

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

const healthCheckKey = "resilience:health-check"

type dependencyContainer struct {
	dbConnector         postgresql.Connector
	transactionalClient postgresql.TransactionalClient
	redisClient         redis.Client
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

	transactionalClient := migrator.connector.TransactionalClient()
	connectionProvider := postgresql.NewConnectionProvider(transactionalClient)
	transactionFactory := repo.NewTransactionFactory(connectionProvider)
	executor := transactional.NewExecutor[app.RepoProvider](transactionFactory)

	middlewares := []middleware.Middleware{
		ogenmiddleware.NewLoggingMiddleware(logger),
	}

	userService := appservice.NewUserService(executor)
	userQueryService := query.NewUserQueryService(transactionalClient)

	bandService := appservice.NewBandService(executor)
	bandQueryService := query.NewBandQueryService(transactionalClient)

	trackService := appservice.NewTrackService(executor)
	trackQueryService := query.NewTrackQueryService(transactionalClient)

	setlistService := appservice.NewSetlistService(executor)
	setlistQueryService := query.NewSetlistQueryService(transactionalClient)

	redisClient := redis.NewClient(config.redisConfig())
	sessionStore := redisinfra.NewSessionStore(redisClient, redisinfra.Config{
		TTL:                config.SessionTTL,
		MaxSessionsPerUser: config.SessionMaxPerUser,
	})

	apiServer, err := transport.NewAPIServer(
		errorHandler,
		middlewares,
		userService,
		userQueryService,
		bandService,
		bandQueryService,
		trackService,
		trackQueryService,
		setlistService,
		setlistQueryService,
		sessionStore,
	)
	if err != nil {
		return nil, err
	}
	router.PathPrefix("/api").Handler(corsMiddleware(config.CORSAllowedOrigins)(apiServer))

	return &dependencyContainer{
		dbConnector:         migrator.connector,
		transactionalClient: transactionalClient,
		redisClient:         redisClient,
	}, nil
}

// Ready reports whether dependencies are reachable and the service can serve traffic.
func (c *dependencyContainer) Ready(ctx context.Context) error {
	var dbAlive int
	if err := c.transactionalClient.GetContext(ctx, &dbAlive, "SELECT 1"); err != nil {
		return errors.Wrap(err, "database is not reachable")
	}
	if _, err := c.redisClient.Exists(ctx, healthCheckKey); err != nil {
		return errors.Wrap(err, "redis is not reachable")
	}
	return nil
}

func (c *dependencyContainer) Close() error {
	dbErr := c.dbConnector.Close()
	redisErr := c.redisClient.Close()
	return errors.Join(dbErr, redisErr)
}
