package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/transactional"
	"cadence/pkg/common/uuid"
)

func NewUserService(executor transactional.Executor[app.RepoProvider]) *UserService {
	return &UserService{executor: executor}
}

type UserService struct {
	executor transactional.Executor[app.RepoProvider]
}

func (s *UserService) Authenticate(ctx context.Context, username, password string) (userID uuid.UUID, err error) {
	err = s.executor.Execute(ctx, func(repoProvider app.RepoProvider) error {
		repo := repoProvider.UserRepository()

		user, findErr := repo.FindByUsername(username)
		if findErr != nil {
			if errors.Is(findErr, domain.ErrUserNotFound) {
				return domain.ErrInvalidCredentials
			}
			return findErr
		}

		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash()), []byte(password)) != nil {
			return domain.ErrInvalidCredentials
		}

		userID = user.ID()
		return nil
	})
	return userID, err
}

func (s *UserService) Register(ctx context.Context, username, password string) (userID uuid.UUID, err error) {
	err = s.executor.Execute(ctx, func(repoProvider app.RepoProvider) error {
		repo := repoProvider.UserRepository()

		_, findErr := repo.FindByUsername(username)
		if findErr == nil {
			return domain.ErrUsernameTaken
		}
		if !errors.Is(findErr, domain.ErrUserNotFound) {
			return findErr
		}

		passwordHash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}

		user, userErr := domain.NewUser(repo.NextID(), username, string(passwordHash))
		if userErr != nil {
			return userErr
		}

		if storeErr := repo.Store(user); storeErr != nil {
			return storeErr
		}

		userID = user.ID()
		return nil
	})
	return userID, err
}
