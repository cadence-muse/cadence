package domain

import (
	"errors"
	"fmt"
)

const (
	maxUsernameLength = 100
)

var (
	ErrEmptyUsername     = errors.New("username can not be empty")
	ErrUsernameTooLong   = fmt.Errorf("username length should be less than or equal to %d", maxUsernameLength)
	ErrEmptyPasswordHash = errors.New("password hash can not be empty")
	ErrUserNotFound      = errors.New("user not found")
	ErrUsernameTaken     = errors.New("username is already taken")
)

type User struct {
	id           UserID
	username     string
	passwordHash string
}

type UserRepository interface {
	NextID() UserID
	Store(*User) error
	Get(UserID) (*User, error)
	FindByUsername(username string) (*User, error)
}

func NewUser(
	id UserID,
	username string,
	passwordHash string,
) (*User, error) {
	err := validateUsernameLength(username)
	if err != nil {
		return nil, err
	}
	err = validatePasswordHash(passwordHash)
	if err != nil {
		return nil, err
	}
	return &User{
		id:           id,
		username:     username,
		passwordHash: passwordHash,
	}, nil
}

func LoadUser(
	id UserID,
	username string,
	passwordHash string,
) *User {
	return &User{
		id:           id,
		username:     username,
		passwordHash: passwordHash,
	}
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Username() string {
	return u.username
}

func (u *User) PasswordHash() string {
	return u.passwordHash
}

func validateUsernameLength(username string) error {
	return checkStringLimits(username, maxUsernameLength, ErrEmptyUsername, ErrUsernameTooLong)
}

func validatePasswordHash(passwordHash string) error {
	if passwordHash == "" {
		return ErrEmptyPasswordHash
	}
	return nil
}
