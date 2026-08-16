package service

import (
	"context"
	"fmt"
	"time"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/maybe"
	"cadence/pkg/common/transactional"
	"cadence/pkg/common/uuid"
)

func NewSetlistService(executor transactional.Executor[app.RepoProvider]) *SetlistService {
	return &SetlistService{executor: executor}
}

type SetlistService struct {
	executor transactional.Executor[app.RepoProvider]
}

type CreateSetlistParams struct {
	BandID        uuid.UUID
	Name          string
	EventLocation maybe.Maybe[string]
	EventDate     maybe.Maybe[time.Time]
	TrackIDs      []uuid.UUID
}

type UpdateSetlistParams struct {
	BandID        uuid.UUID
	SetlistID     uuid.UUID
	Name          maybe.Maybe[string]
	EventLocation maybe.Maybe[string]
	EventDate     maybe.Maybe[time.Time]
}

func (s *SetlistService) Create(ctx context.Context, params CreateSetlistParams) (setlistID uuid.UUID, err error) {
	err = s.executor.Execute(ctx, func(repoProvider app.RepoProvider) error {
		if _, bandErr := repoProvider.BandRepository().Get(params.BandID); bandErr != nil {
			return bandErr
		}
		if trackErr := verifyTracksBelongToBand(repoProvider, params.BandID, params.TrackIDs); trackErr != nil {
			return trackErr
		}

		repo := repoProvider.SetlistRepository()

		setlist, setlistErr := domain.NewSetlist(
			repo.NextID(),
			params.BandID,
			params.Name,
			params.EventLocation,
			params.EventDate,
			params.TrackIDs,
		)
		if setlistErr != nil {
			return setlistErr
		}

		if storeErr := repo.Store(setlist); storeErr != nil {
			return storeErr
		}

		setlistID = setlist.ID()
		return nil
	})
	return setlistID, err
}

func (s *SetlistService) Update(ctx context.Context, params UpdateSetlistParams) error {
	return s.executor.ExecuteWithLock(ctx, getSetlistLockName(params.SetlistID), func(repoProvider app.RepoProvider) error {
		setlist, err := getBandSetlist(repoProvider, params.BandID, params.SetlistID)
		if err != nil {
			return err
		}

		if name, ok := maybe.JustValid(params.Name); ok {
			if err := setlist.SetName(name); err != nil {
				return err
			}
		}
		if maybe.IsSet(params.EventLocation) {
			if err := setlist.SetEventLocation(params.EventLocation); err != nil {
				return err
			}
		}
		if maybe.IsSet(params.EventDate) {
			setlist.SetEventDate(params.EventDate)
		}

		return repoProvider.SetlistRepository().Store(setlist)
	})
}

func (s *SetlistService) Remove(ctx context.Context, bandID, setlistID uuid.UUID) error {
	return s.executor.ExecuteWithLock(ctx, getSetlistLockName(setlistID), func(repoProvider app.RepoProvider) error {
		if _, err := getBandSetlist(repoProvider, bandID, setlistID); err != nil {
			return err
		}

		return repoProvider.SetlistRepository().Remove(setlistID)
	})
}

func (s *SetlistService) AddTrack(ctx context.Context, bandID, setlistID, trackID uuid.UUID) error {
	return s.executor.ExecuteWithLock(ctx, getSetlistLockName(setlistID), func(repoProvider app.RepoProvider) error {
		setlist, err := getBandSetlist(repoProvider, bandID, setlistID)
		if err != nil {
			return err
		}

		track, err := repoProvider.TrackRepository().Get(trackID)
		if err != nil {
			return err
		}
		if track.BandID() != bandID {
			return domain.ErrTrackNotFound
		}

		if err := setlist.AddTrack(trackID); err != nil {
			return err
		}

		return repoProvider.SetlistRepository().Store(setlist)
	})
}

func (s *SetlistService) AddTracks(ctx context.Context, bandID, setlistID uuid.UUID, trackIDs []uuid.UUID) error {
	return s.executor.ExecuteWithLock(ctx, getSetlistLockName(setlistID), func(repoProvider app.RepoProvider) error {
		setlist, err := getBandSetlist(repoProvider, bandID, setlistID)
		if err != nil {
			return err
		}

		if err := verifyTracksBelongToBand(repoProvider, bandID, trackIDs); err != nil {
			return err
		}

		for _, trackID := range trackIDs {
			if err := setlist.AddTrack(trackID); err != nil {
				return err
			}
		}

		return repoProvider.SetlistRepository().Store(setlist)
	})
}

func (s *SetlistService) RemoveTrack(ctx context.Context, bandID, setlistID, trackID uuid.UUID) error {
	return s.executor.ExecuteWithLock(ctx, getSetlistLockName(setlistID), func(repoProvider app.RepoProvider) error {
		setlist, err := getBandSetlist(repoProvider, bandID, setlistID)
		if err != nil {
			return err
		}

		if err := setlist.RemoveTrack(trackID); err != nil {
			return err
		}

		return repoProvider.SetlistRepository().Store(setlist)
	})
}

func (s *SetlistService) ReorderTracks(ctx context.Context, bandID, setlistID uuid.UUID, trackIDs []uuid.UUID) error {
	return s.executor.ExecuteWithLock(ctx, getSetlistLockName(setlistID), func(repoProvider app.RepoProvider) error {
		setlist, err := getBandSetlist(repoProvider, bandID, setlistID)
		if err != nil {
			return err
		}

		if err := setlist.Reorder(trackIDs); err != nil {
			return err
		}

		return repoProvider.SetlistRepository().Store(setlist)
	})
}

func getBandSetlist(repoProvider app.RepoProvider, bandID, setlistID uuid.UUID) (*domain.Setlist, error) {
	setlist, err := repoProvider.SetlistRepository().Get(setlistID)
	if err != nil {
		return nil, err
	}
	if setlist.BandID() != bandID {
		return nil, domain.ErrSetlistNotFound
	}
	return setlist, nil
}

func verifyTracksBelongToBand(repoProvider app.RepoProvider, bandID uuid.UUID, trackIDs []uuid.UUID) error {
	trackRepo := repoProvider.TrackRepository()
	checked := make(map[uuid.UUID]struct{}, len(trackIDs))
	for _, trackID := range trackIDs {
		if _, alreadyChecked := checked[trackID]; alreadyChecked {
			continue
		}
		checked[trackID] = struct{}{}

		track, err := trackRepo.Get(trackID)
		if err != nil {
			return err
		}
		if track.BandID() != bandID {
			return domain.ErrTrackNotFound
		}
	}
	return nil
}

func getSetlistLockName(id uuid.UUID) string {
	return fmt.Sprintf("setlist_%s", id.String())
}
