package service

import (
	"context"
	"testing"
	"time"

	"github.com/nightnoryu/go-kita/maybe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/uuid"
	"cadence/pkg/common/valuetypes"
)

func TestTrackService_Create(t *testing.T) {
	t.Run("creates a track for the band", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewTrackService(executor)
		bandID := uuid.Generate()

		trackID, err := svc.Create(context.Background(), CreateTrackParams{
			BandID: bandID,
			Title:  "Song",
			Artist: "Artist",
		})
		require.NoError(t, err)

		track, err := executor.repoProvider().TrackRepository().Get(trackID)
		require.NoError(t, err)
		assert.Equal(t, bandID, track.BandID())
		assert.Equal(t, "Song", track.Title())
		assert.Equal(t, "Artist", track.Artist())
	})

	t.Run("domain validation error is propagated", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewTrackService(executor)

		_, err := svc.Create(context.Background(), CreateTrackParams{
			BandID: uuid.Generate(),
			Title:  "",
			Artist: "Artist",
		})
		require.ErrorIs(t, err, domain.ErrEmptyTrackTitle)
	})
}

func TestTrackService_Update(t *testing.T) {
	t.Run("updates an existing track", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewTrackService(executor)
		bandID := uuid.Generate()
		trackID := seedTrack(t, executor, bandID)

		err := svc.Update(context.Background(), UpdateTrackParams{
			BandID:  bandID,
			TrackID: trackID,
			Title:   maybe.NewJust("New Title"),
			Tempo:   maybe.NewJust(120),
		})
		require.NoError(t, err)

		track, err := executor.repoProvider().TrackRepository().Get(trackID)
		require.NoError(t, err)
		assert.Equal(t, "New Title", track.Title())
		assert.Equal(t, maybe.NewJust(120), track.Tempo())
	})

	t.Run("track not found", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewTrackService(executor)

		err := svc.Update(context.Background(), UpdateTrackParams{
			BandID:  uuid.Generate(),
			TrackID: uuid.Generate(),
			Title:   maybe.NewJust("New Title"),
		})
		require.ErrorIs(t, err, domain.ErrTrackNotFound)
	})

	t.Run("track belongs to a different band", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewTrackService(executor)
		trackID := seedTrack(t, executor, uuid.Generate())

		err := svc.Update(context.Background(), UpdateTrackParams{
			BandID:  uuid.Generate(),
			TrackID: trackID,
			Title:   maybe.NewJust("New Title"),
		})
		require.ErrorIs(t, err, domain.ErrTrackNotFound)
	})
}

func TestTrackService_Remove(t *testing.T) {
	t.Run("removes an existing track", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewTrackService(executor)
		bandID := uuid.Generate()
		trackID := seedTrack(t, executor, bandID)

		err := svc.Remove(context.Background(), bandID, trackID)
		require.NoError(t, err)

		_, err = executor.repoProvider().TrackRepository().Get(trackID)
		require.ErrorIs(t, err, domain.ErrTrackNotFound)
	})

	t.Run("track belongs to a different band", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewTrackService(executor)
		trackID := seedTrack(t, executor, uuid.Generate())

		err := svc.Remove(context.Background(), uuid.Generate(), trackID)
		require.ErrorIs(t, err, domain.ErrTrackNotFound)

		_, err = executor.repoProvider().TrackRepository().Get(trackID)
		require.NoError(t, err)
	})
}

func seedTrack(t *testing.T, executor *fakeExecutor, bandID uuid.UUID) uuid.UUID {
	t.Helper()

	repo := executor.repoProvider().TrackRepository()
	track, err := domain.NewTrack(
		repo.NextID(),
		bandID,
		"Track",
		"Artist",
		maybe.NewAbsent[time.Duration](),
		maybe.NewAbsent[int](),
		maybe.NewAbsent[valuetypes.MusicalKey](),
		maybe.NewAbsent[string](),
	)
	require.NoError(t, err)
	require.NoError(t, repo.Store(track))
	return track.ID()
}
