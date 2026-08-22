package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	maxUsernameLength = 100
	minPasswordLength = 8
	maxPasswordLength = 127
)

var (
	ErrEmptyUsername   = errors.New("username can not be empty")
	ErrUsernameTooLong = fmt.Errorf("username length must be less than or equal to %d", maxUsernameLength)

	ErrInvalidPassword = fmt.Errorf("password length must be more than or equal to %d and less than or equal to %d", minPasswordLength, maxPasswordLength)

	ErrUserNotFound       = errors.New("user not found")
	ErrUsernameTaken      = errors.New("username is already taken")
	ErrInvalidCredentials = errors.New("invalid username or password")
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
	password string,
) (*User, error) {
	err := validateUsernameLength(username)
	if err != nil {
		return nil, err
	}
	err = validatePasswordLength(password)
	if err != nil {
		return nil, err
	}
	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		id:           id,
		username:     username,
		passwordHash: hash,
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

func (u *User) ComparePassword(password string) error {
	if bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password)) != nil {
		return ErrInvalidCredentials
	}
	return nil
}

func (u *User) SetPassword(password string) error {
	if err := validatePasswordLength(password); err != nil {
		return err
	}
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}
	u.passwordHash = hash
	return nil
}

func validateUsernameLength(username string) error {
	return checkStringLimits(username, maxUsernameLength, ErrEmptyUsername, ErrUsernameTooLong)
}

func validatePasswordLength(password string) error {
	password = strings.TrimSpace(password)
	length := utf8.RuneCountInString(password)
	if length < minPasswordLength || length > maxPasswordLength {
		return ErrInvalidPassword
	}
	return nil
}

func hashPassword(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(result), err
}
