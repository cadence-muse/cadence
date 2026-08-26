package service

import (
	"context"
	"errors"

	"github.com/nightnoryu/go-kita/transactional"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/uuid"
)

func NewUserService(
	executor transactional.Executor[app.RepoProvider],
	sessionStore app.SessionStore,
) *UserService {
	return &UserService{
		executor:     executor,
		sessionStore: sessionStore,
	}
}

type UserService struct {
	executor     transactional.Executor[app.RepoProvider]
	sessionStore app.SessionStore
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

		user, userErr := domain.NewUser(repo.NextID(), username, password)
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

func (s *UserService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	err := s.executor.Execute(ctx, func(repoProvider app.RepoProvider) error {
		repo := repoProvider.UserRepository()

		user, err := repo.Get(userID)
		if err != nil {
			return err
		}

		err = user.ComparePassword(currentPassword)
		if err != nil {
			return err
		}

		err = user.SetPassword(newPassword)
		if err != nil {
			return err
		}

		return repo.Store(user)
	})
	if err != nil {
		return err
	}

	return s.sessionStore.DeleteAllSessions(ctx, userID)
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

		err = user.ComparePassword(password)
		if err != nil {
			return domain.ErrInvalidCredentials
		}

		userID = user.ID()
		return nil
	})
	return userID, err
}
