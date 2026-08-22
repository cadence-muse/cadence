package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cadence/pkg/common/uuid"
)

func TestNewUser(t *testing.T) {
	t.Run("valid user is created", func(t *testing.T) {
		id := uuid.Generate()

		user, err := NewUser(id, "alice", "correct horse battery staple")
		require.NoError(t, err)
		assert.Equal(t, id, user.ID())
		assert.Equal(t, "alice", user.Username())
		assert.NotEqual(t, "correct horse battery staple", user.PasswordHash())
	})

	t.Run("empty username is rejected", func(t *testing.T) {
		_, err := NewUser(uuid.Generate(), "", "hash")
		assert.ErrorIs(t, err, ErrEmptyUsername)
	})

	t.Run("username over the length limit is rejected", func(t *testing.T) {
		_, err := NewUser(uuid.Generate(), strings.Repeat("a", maxUsernameLength+1), "correct horse battery staple")
		assert.ErrorIs(t, err, ErrUsernameTooLong)
	})

	t.Run("username at the length limit is accepted", func(t *testing.T) {
		_, err := NewUser(uuid.Generate(), strings.Repeat("a", maxUsernameLength), "correct horse battery staple")
		assert.NoError(t, err)
	})

	t.Run("empty password hash is rejected", func(t *testing.T) {
		_, err := NewUser(uuid.Generate(), "alice", "")
		assert.ErrorIs(t, err, ErrInvalidPassword)
	})
}

func TestLoadUser(t *testing.T) {
	id := uuid.Generate()

	user := LoadUser(id, "loaded", "hash")

	assert.Equal(t, id, user.ID())
	assert.Equal(t, "loaded", user.Username())
	assert.Equal(t, "hash", user.PasswordHash())
}
