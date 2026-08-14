package ogenerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInvalidInputError(t *testing.T) {
	err := NewInvalidInputError("invalid title")

	assert.Equal(t, ErrCodeInvalidInput, err.Code)
	assert.Equal(t, "invalid title", err.Message)
	assert.Nil(t, err.Details)
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("band not found")

	assert.Equal(t, ErrCodeNotFound, err.Code)
	assert.Equal(t, "band not found", err.Message)
	assert.Nil(t, err.Details)
}

func TestNewPermissionDeniedError(t *testing.T) {
	err := NewPermissionDeniedError("not allowed")

	assert.Equal(t, ErrCodePermissionDenied, err.Code)
	assert.Equal(t, "not allowed", err.Message)
	assert.Nil(t, err.Details)
}

func TestNewOperationRejectedError(t *testing.T) {
	err := NewOperationRejectedError("band already active", "band_active")

	assert.Equal(t, ErrCodeOperationRejected, err.Code)
	assert.Equal(t, "band already active", err.Message)
	require.NotNil(t, err.Details)
	assert.Equal(t, "band_active", err.Details["code"])
}

func TestNewAlreadyExistsError(t *testing.T) {
	err := NewAlreadyExistsError("track already exists")

	assert.Equal(t, ErrCodeAlreadyExists, err.Code)
	assert.Equal(t, "track already exists", err.Message)
	assert.Nil(t, err.Details)
}

func TestError_Error(t *testing.T) {
	err := &Error{Code: ErrCodeNotFound, Message: "band not found"}

	assert.Equal(t, "not_found: band not found", err.Error())
}

func TestError_ImplementsErrorInterface(t *testing.T) {
	var _ error = (*Error)(nil)

	var err error = NewNotFoundError("band not found")

	require.Error(t, err)
	assert.Equal(t, "not_found: band not found", err.Error())
}
