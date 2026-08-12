package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"cadence/pkg/cadence/domain"
)

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

type UserService struct {
	repo domain.UserRepository
}

type RegisterInput struct {
	Username string
	Password string
}

func (s *UserService) Register(input RegisterInput) (*domain.User, error) {
	_, err := s.repo.FindByUsername(input.Username)
	if err == nil {
		return nil, domain.ErrUsernameTaken
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(s.repo.NextID(), input.Username, string(passwordHash))
	if err != nil {
		return nil, err
	}

	err = s.repo.Store(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
