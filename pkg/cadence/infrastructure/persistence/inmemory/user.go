package inmemory

import "cadence/pkg/cadence/domain"

func NewUserRepository(store *Store) domain.UserRepository {
	return &userRepository{mapRepo: store.users}
}

type userRepository struct {
	*mapRepo[*domain.User]
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	for _, user := range r.items {
		if user.Username() == username {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}
