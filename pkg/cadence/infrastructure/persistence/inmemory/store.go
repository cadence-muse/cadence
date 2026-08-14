package inmemory

import (
	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/uuid"
)

func NewStore() *Store {
	return &Store{
		users:    newMapRepo[*domain.User](domain.ErrUserNotFound),
		bands:    newMapRepo[*domain.Band](domain.ErrBandNotFound),
		tracks:   newMapRepo[*domain.Track](domain.ErrTrackNotFound),
		setlists: newMapRepo[*domain.Setlist](domain.ErrSetlistNotFound),
	}
}

type Store struct {
	users    *mapRepo[*domain.User]
	bands    *mapRepo[*domain.Band]
	tracks   *mapRepo[*domain.Track]
	setlists *mapRepo[*domain.Setlist]
}

type identifiable interface {
	ID() uuid.UUID
}

func newMapRepo[T identifiable](notFound error) *mapRepo[T] {
	return &mapRepo[T]{items: make(map[uuid.UUID]T), notFound: notFound}
}

type mapRepo[T identifiable] struct {
	items    map[uuid.UUID]T
	notFound error
}

func (r *mapRepo[T]) NextID() uuid.UUID {
	return uuid.Generate()
}

func (r *mapRepo[T]) Store(item T) error {
	r.items[item.ID()] = item
	return nil
}

func (r *mapRepo[T]) Get(id uuid.UUID) (T, error) {
	item, ok := r.items[id]
	if !ok {
		var zero T
		return zero, r.notFound
	}
	return item, nil
}

func (r *mapRepo[T]) Remove(id uuid.UUID) error {
	if _, ok := r.items[id]; !ok {
		return r.notFound
	}
	delete(r.items, id)
	return nil
}
