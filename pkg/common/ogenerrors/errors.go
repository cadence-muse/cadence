package ogenerrors

import (
	"fmt"
)

type ErrorCode string

const (
	ErrCodeInvalidInput      ErrorCode = "invalid_input"      // The client specified an invalid input argument
	ErrCodeNotFound          ErrorCode = "not_found"          // Requested entity was not found
	ErrCodePermissionDenied  ErrorCode = "permission_denied"  // The caller does not have permission to execute the specified operation
	ErrCodeOperationRejected ErrorCode = "operation_rejected" // The operation was rejected because the system is not in a state required for the operation's execution
	ErrCodeAlreadyExists     ErrorCode = "already_exists"     // The operation was rejected because the system is not in a state required for the operation's execution
)

type ErrorDetails map[string]interface{}

type Error struct {
	Code    ErrorCode
	Message string
	Details ErrorDetails
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewInvalidInputError(message string) *Error {
	return &Error{
		Code:    ErrCodeInvalidInput,
		Message: message,
	}
}

func NewNotFoundError(message string) *Error {
	return &Error{
		Code:    ErrCodeNotFound,
		Message: message,
	}
}

func NewPermissionDeniedError(message string) *Error {
	return &Error{
		Code:    ErrCodePermissionDenied,
		Message: message,
	}
}

func NewOperationRejectedError(message, code string) *Error {
	return &Error{
		Code:    ErrCodeOperationRejected,
		Message: message,
		Details: ErrorDetails{
			"code": code,
		},
	}
}

func NewAlreadyExistsError(message string) *Error {
	return &Error{
		Code:    ErrCodeAlreadyExists,
		Message: message,
	}
}
