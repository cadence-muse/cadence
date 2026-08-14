package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cadence/pkg/cadence/domain"
)

func TestUserService_Register(t *testing.T) {
	t.Run("registers a new user", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewUserService(executor)

		userID, err := svc.Register(context.Background(), "alice", "s3cret-password")
		require.NoError(t, err)

		user, err := executor.repoProvider().UserRepository().Get(userID)
		require.NoError(t, err)
		assert.Equal(t, "alice", user.Username())
		assert.NotEqual(t, "s3cret-password", user.PasswordHash())
	})

	t.Run("username already taken is rejected", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewUserService(executor)

		_, err := svc.Register(context.Background(), "alice", "s3cret-password")
		require.NoError(t, err)

		_, err = svc.Register(context.Background(), "alice", "another-password")
		require.ErrorIs(t, err, domain.ErrUsernameTaken)
	})
}

func TestUserService_Authenticate(t *testing.T) {
	t.Run("authenticates with correct credentials", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewUserService(executor)
		registeredID, err := svc.Register(context.Background(), "alice", "s3cret-password")
		require.NoError(t, err)

		userID, err := svc.Authenticate(context.Background(), "alice", "s3cret-password")
		require.NoError(t, err)
		assert.Equal(t, registeredID, userID)
	})

	t.Run("wrong password is rejected", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewUserService(executor)
		_, err := svc.Register(context.Background(), "alice", "s3cret-password")
		require.NoError(t, err)

		_, err = svc.Authenticate(context.Background(), "alice", "wrong-password")
		require.ErrorIs(t, err, domain.ErrInvalidCredentials)
	})

	t.Run("unknown username is rejected with the same error as a wrong password", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewUserService(executor)

		_, err := svc.Authenticate(context.Background(), "unknown", "any-password")
		require.ErrorIs(t, err, domain.ErrInvalidCredentials)
		assert.NotErrorIs(t, err, domain.ErrUserNotFound)
	})
}
