package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/maybe"
	"cadence/pkg/common/uuid"
)

func TestSetlistService_Create(t *testing.T) {
	t.Run("creates a setlist with valid tracks", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewSetlistService(executor)
		bandID := seedBand(t, executor, uuid.Generate())
		trackID := seedTrack(t, executor, bandID)

		setlistID, err := svc.Create(context.Background(), CreateSetlistParams{
			BandID:   bandID,
			Name:     "Summer Show",
			TrackIDs: []uuid.UUID{trackID},
		})
		require.NoError(t, err)

		setlist, err := executor.repoProvider().SetlistRepository().Get(setlistID)
		require.NoError(t, err)
		assert.Equal(t, bandID, setlist.BandID())
		assert.Equal(t, "Summer Show", setlist.Name())
		assert.Equal(t, []uuid.UUID{trackID}, setlist.TrackIDs())
	})

	t.Run("track not belonging to the band is rejected", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewSetlistService(executor)
		bandID := seedBand(t, executor, uuid.Generate())
		otherBandTrackID := seedTrack(t, executor, uuid.Generate())

		_, err := svc.Create(context.Background(), CreateSetlistParams{
			BandID:   bandID,
			Name:     "Summer Show",
			TrackIDs: []uuid.UUID{otherBandTrackID},
		})
		require.ErrorIs(t, err, domain.ErrTrackNotFound)
	})

	t.Run("band not found", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewSetlistService(executor)

		_, err := svc.Create(context.Background(), CreateSetlistParams{
			BandID: uuid.Generate(),
			Name:   "Summer Show",
		})
		require.ErrorIs(t, err, domain.ErrBandNotFound)
	})
}

func TestSetlistService_Update(t *testing.T) {
	t.Run("updates an existing setlist", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewSetlistService(executor)
		bandID := uuid.Generate()
		setlistID := seedSetlist(t, executor, bandID, nil)

		err := svc.Update(context.Background(), UpdateSetlistParams{
			BandID:    bandID,
			SetlistID: setlistID,
			Name:      maybe.NewJust("New Name"),
		})
		require.NoError(t, err)

		setlist, err := executor.repoProvider().SetlistRepository().Get(setlistID)
		require.NoError(t, err)
		assert.Equal(t, "New Name", setlist.Name())
	})

	t.Run("setlist belongs to a different band", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewSetlistService(executor)
		setlistID := seedSetlist(t, executor, uuid.Generate(), nil)

		err := svc.Update(context.Background(), UpdateSetlistParams{
			BandID:    uuid.Generate(),
			SetlistID: setlistID,
			Name:      maybe.NewJust("New Name"),
		})
		require.ErrorIs(t, err, domain.ErrSetlistNotFound)
	})
}

func TestSetlistService_Remove(t *testing.T) {
	t.Run("removes an existing setlist", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewSetlistService(executor)
		bandID := uuid.Generate()
		setlistID := seedSetlist(t, executor, bandID, nil)

		err := svc.Remove(context.Background(), bandID, setlistID)
		require.NoError(t, err)

		_, err = executor.repoProvider().SetlistRepository().Get(setlistID)
		require.ErrorIs(t, err, domain.ErrSetlistNotFound)
	})

	t.Run("setlist belongs to a different band", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewSetlistService(executor)
		setlistID := seedSetlist(t, executor, uuid.Generate(), nil)

		err := svc.Remove(context.Background(), uuid.Generate(), setlistID)
		require.ErrorIs(t, err, domain.ErrSetlistNotFound)
	})
}

func TestSetlistService_AddTrack(t *testing.T) {
	t.Run("adds a track from the same band", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewSetlistService(executor)
		bandID := uuid.Generate()
		setlistID := seedSetlist(t, executor, bandID, nil)
		trackID := seedTrack(t, executor, bandID)

		err := svc.AddTrack(context.Background(), bandID, setlistID, trackID)
		require.NoError(t, err)

		setlist, err := executor.repoProvider().SetlistRepository().Get(setlistID)
		require.NoError(t, err)
		assert.Equal(t, []uuid.UUID{trackID}, setlist.TrackIDs())
	})

	t.Run("track from a different band is rejected", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewSetlistService(executor)
		bandID := uuid.Generate()
		setlistID := seedSetlist(t, executor, bandID, nil)
		otherBandTrackID := seedTrack(t, executor, uuid.Generate())

		err := svc.AddTrack(context.Background(), bandID, setlistID, otherBandTrackID)
		require.ErrorIs(t, err, domain.ErrTrackNotFound)
	})
}

func TestSetlistService_RemoveTrack(t *testing.T) {
	t.Run("removes a track already in the setlist", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewSetlistService(executor)
		bandID := uuid.Generate()
		trackID := seedTrack(t, executor, bandID)
		setlistID := seedSetlist(t, executor, bandID, []uuid.UUID{trackID})

		err := svc.RemoveTrack(context.Background(), bandID, setlistID, trackID)
		require.NoError(t, err)

		setlist, err := executor.repoProvider().SetlistRepository().Get(setlistID)
		require.NoError(t, err)
		assert.Empty(t, setlist.TrackIDs())
	})

	t.Run("track not in the setlist propagates domain error", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewSetlistService(executor)
		bandID := uuid.Generate()
		setlistID := seedSetlist(t, executor, bandID, nil)

		err := svc.RemoveTrack(context.Background(), bandID, setlistID, uuid.Generate())
		require.ErrorIs(t, err, domain.ErrTrackNotInSetlist)
	})
}

func TestSetlistService_ReorderTracks(t *testing.T) {
	t.Run("reorders tracks currently in the setlist", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewSetlistService(executor)
		bandID := uuid.Generate()
		trackA := seedTrack(t, executor, bandID)
		trackB := seedTrack(t, executor, bandID)
		setlistID := seedSetlist(t, executor, bandID, []uuid.UUID{trackA, trackB})

		err := svc.ReorderTracks(context.Background(), bandID, setlistID, []uuid.UUID{trackB, trackA})
		require.NoError(t, err)

		setlist, err := executor.repoProvider().SetlistRepository().Get(setlistID)
		require.NoError(t, err)
		assert.Equal(t, []uuid.UUID{trackB, trackA}, setlist.TrackIDs())
	})

	t.Run("invalid order propagates domain error", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewSetlistService(executor)
		bandID := uuid.Generate()
		trackA := seedTrack(t, executor, bandID)
		setlistID := seedSetlist(t, executor, bandID, []uuid.UUID{trackA})

		err := svc.ReorderTracks(context.Background(), bandID, setlistID, []uuid.UUID{trackA, uuid.Generate()})
		require.ErrorIs(t, err, domain.ErrInvalidSetlistTrackOrder)
	})
}

func seedSetlist(t *testing.T, executor *fakeExecutor, bandID uuid.UUID, trackIDs []uuid.UUID) uuid.UUID {
	t.Helper()

	repo := executor.repoProvider().SetlistRepository()
	setlist, err := domain.NewSetlist(
		repo.NextID(),
		bandID,
		"Setlist",
		maybe.NewAbsent[string](),
		maybe.NewAbsent[time.Time](),
		trackIDs,
	)
	require.NoError(t, err)
	require.NoError(t, repo.Store(setlist))
	return setlist.ID()
}
