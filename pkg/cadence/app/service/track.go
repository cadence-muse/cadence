package service

import (
	"context"
	"time"

	"cadence/pkg/cadence/app"
	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/maybe"
	"cadence/pkg/common/transactional"
	"cadence/pkg/common/uuid"
	"cadence/pkg/common/valuetypes"
)

func NewTrackService(executor transactional.Executor[app.RepoProvider]) *TrackService {
	return &TrackService{executor: executor}
}

type TrackService struct {
	executor transactional.Executor[app.RepoProvider]
}

type CreateTrackParams struct {
	BandID   uuid.UUID
	Title    string
	Artist   string
	Duration maybe.Maybe[time.Duration]
	Tempo    maybe.Maybe[int]
	Key      maybe.Maybe[valuetypes.MusicalKey]
	Notes    maybe.Maybe[string]
}

type UpdateTrackParams struct {
	BandID   uuid.UUID
	TrackID  uuid.UUID
	Title    maybe.Maybe[string]
	Artist   maybe.Maybe[string]
	Duration maybe.Maybe[time.Duration]
	Tempo    maybe.Maybe[int]
	Key      maybe.Maybe[valuetypes.MusicalKey]
	Notes    maybe.Maybe[string]
}

func (s *TrackService) Create(ctx context.Context, params CreateTrackParams) (trackID uuid.UUID, err error) {
	err = s.executor.Execute(ctx, func(repoProvider app.RepoProvider) error {
		trackRepo := repoProvider.TrackRepository()

		track, trackErr := domain.NewTrack(
			trackRepo.NextID(),
			params.BandID,
			params.Title,
			params.Artist,
			params.Duration,
			params.Tempo,
			params.Key,
			params.Notes,
		)
		if trackErr != nil {
			return trackErr
		}

		if storeErr := trackRepo.Store(track); storeErr != nil {
			return storeErr
		}

		trackID = track.ID()
		return nil
	})
	return trackID, err
}

func (s *TrackService) Update(ctx context.Context, params UpdateTrackParams) error {
	return s.executor.Execute(ctx, func(repoProvider app.RepoProvider) error {
		repo := repoProvider.TrackRepository()

		track, err := repo.Get(params.TrackID)
		if err != nil {
			return err
		}
		if track.BandID() != params.BandID {
			return domain.ErrTrackNotFound
		}

		if title, ok := maybe.JustValid(params.Title); ok {
			if err := track.SetTitle(title); err != nil {
				return err
			}
		}
		if artist, ok := maybe.JustValid(params.Artist); ok {
			if err := track.SetArtist(artist); err != nil {
				return err
			}
		}
		if maybe.IsSet(params.Duration) {
			track.SetDuration(params.Duration)
		}
		if maybe.IsSet(params.Tempo) {
			track.SetTempo(params.Tempo)
		}
		if maybe.IsSet(params.Key) {
			track.SetKey(params.Key)
		}
		if maybe.IsSet(params.Notes) {
			track.SetNotes(params.Notes)
		}

		return repo.Store(track)
	})
}

func (s *TrackService) Remove(ctx context.Context, bandID, trackID uuid.UUID) error {
	return s.executor.Execute(ctx, func(repoProvider app.RepoProvider) error {
		repo := repoProvider.TrackRepository()

		track, err := repo.Get(trackID)
		if err != nil {
			return err
		}
		if track.BandID() != bandID {
			return domain.ErrTrackNotFound
		}

		return repo.Remove(trackID)
	})
}
