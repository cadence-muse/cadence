//go:build e2e

package e2e

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nightnoryu/go-kita/jsonlog"
	"github.com/nightnoryu/go-kita/postgresql"
	"github.com/nightnoryu/go-kita/redis"
	"github.com/nightnoryu/go-kita/transactional"
	"github.com/stretchr/testify/require"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"

	"cadence/api/server/publicapi"
	"cadence/data/migrations"
	"cadence/pkg/cadence/app"
	authorizedquery "cadence/pkg/cadence/app/query/authorized"
	appservice "cadence/pkg/cadence/app/service"
	authorizedservice "cadence/pkg/cadence/app/service/authorized"
	pgquery "cadence/pkg/cadence/infrastructure/persistence/postgresql/query"
	"cadence/pkg/cadence/infrastructure/persistence/postgresql/repo"
	redisinfra "cadence/pkg/cadence/infrastructure/persistence/redis"
	"cadence/pkg/cadence/infrastructure/transport"
)

const (
	testDBName     = "cadence"
	testDBUser     = "cadence"
	testDBPassword = "cadence"

	sessionTTL         = 24 * time.Hour
	maxSessionsPerUser = 5
)

type testEnv struct {
	client *publicapi.Client
	sec    *sessionSecuritySource
}

func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx, "postgres:16-alpine",
		tcpostgres.WithDatabase(testDBName),
		tcpostgres.WithUsername(testDBUser),
		tcpostgres.WithPassword(testDBPassword),
		tcpostgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, pgContainer.Terminate(context.Background())) })

	redisContainer, err := tcredis.Run(ctx, "redis:7-alpine")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, redisContainer.Terminate(context.Background())) })

	pgHost, err := pgContainer.Host(ctx)
	require.NoError(t, err)
	pgPort, err := pgContainer.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	redisHost, err := redisContainer.Host(ctx)
	require.NoError(t, err)
	redisPort, err := redisContainer.MappedPort(ctx, "6379/tcp")
	require.NoError(t, err)

	logger := jsonlog.NewLogger(&jsonlog.Config{Level: jsonlog.InfoLevel, AppName: "cadence-e2e"})

	connector := postgresql.NewConnector()
	dsn := postgresql.DSN{
		Host:     pgHost,
		Port:     int(pgPort.Num()),
		Database: testDBName,
		User:     testDBUser,
		Password: testDBPassword,
	}
	require.NoError(t, connector.Open(dsn, postgresql.Config{MaxConnections: 10, ConnectionLifetime: time.Minute}))
	t.Cleanup(func() { _ = connector.Close() })

	migrator, err := connector.Migrator(logger, migrations.FS)
	require.NoError(t, err)
	require.NoError(t, migrator.MigrateUp())

	transactionalClient := connector.TransactionalClient()
	connectionProvider := postgresql.NewConnectionProvider(transactionalClient)
	transactionFactory := repo.NewTransactionFactory(connectionProvider)
	executor := transactional.NewExecutor[app.RepoProvider](transactionFactory)

	redisClient := redis.NewClient(redis.DSN{Host: redisHost, Port: int(redisPort.Num())}, redis.Config{})
	t.Cleanup(func() { require.NoError(t, redisClient.Close()) })

	sessionStore := redisinfra.NewSessionStore(redisClient, redisinfra.Config{
		TTL:                sessionTTL,
		MaxSessionsPerUser: maxSessionsPerUser,
	})

	userService := appservice.NewUserService(executor, sessionStore)
	userQueryService := pgquery.NewUserQueryService(transactionalClient)
	bandService := authorizedservice.NewBandService(appservice.NewBandService(executor), executor)
	bandQueryService := authorizedquery.NewBandQueryService(pgquery.NewBandQueryService(transactionalClient), executor)
	trackService := authorizedservice.NewTrackService(appservice.NewTrackService(executor), executor)
	trackQueryService := authorizedquery.NewTrackQueryService(pgquery.NewTrackQueryService(transactionalClient), executor)
	setlistService := authorizedservice.NewSetlistService(appservice.NewSetlistService(executor), executor)
	setlistQueryService := authorizedquery.NewSetlistQueryService(pgquery.NewSetlistQueryService(transactionalClient), executor)

	apiHandler, err := transport.NewAPIServer(
		transport.ErrorHandler,
		nil,
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
	require.NoError(t, err)

	server := httptest.NewServer(apiHandler)
	t.Cleanup(server.Close)

	sec := &sessionSecuritySource{}
	client, err := publicapi.NewClient(server.URL, sec)
	require.NoError(t, err)

	return &testEnv{client: client, sec: sec}
}

type sessionSecuritySource struct {
	token string
}

func (s *sessionSecuritySource) SessionAuth(_ context.Context, _ publicapi.OperationName) (publicapi.SessionAuth, error) {
	return publicapi.SessionAuth{APIKey: s.token}, nil
}
