package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cadence/pkg/common/maybe"
	"cadence/pkg/common/uuid"
)

func TestNewSetlist(t *testing.T) {
	t.Run("valid setlist is created with given tracks in order", func(t *testing.T) {
		id := uuid.Generate()
		bandID := uuid.Generate()
		trackA, trackB := uuid.Generate(), uuid.Generate()

		setlist, err := NewSetlist(id, bandID, "Summer Show", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), []TrackID{trackA, trackB})
		require.NoError(t, err)
		assert.Equal(t, id, setlist.ID())
		assert.Equal(t, bandID, setlist.BandID())
		assert.Equal(t, "Summer Show", setlist.Name())
		assert.Equal(t, []TrackID{trackA, trackB}, setlist.TrackIDs())
	})

	t.Run("empty track list is accepted", func(t *testing.T) {
		setlist, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Empty", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), nil)
		require.NoError(t, err)
		assert.Empty(t, setlist.TrackIDs())
	})

	t.Run("empty name is rejected", func(t *testing.T) {
		_, err := NewSetlist(uuid.Generate(), uuid.Generate(), "", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), nil)
		assert.ErrorIs(t, err, ErrEmptySetlistName)
	})

	t.Run("name over the length limit is rejected", func(t *testing.T) {
		_, err := NewSetlist(uuid.Generate(), uuid.Generate(), strings.Repeat("a", maxSetlistNameLength+1), maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), nil)
		assert.ErrorIs(t, err, ErrSetlistNameTooLong)
	})

	t.Run("event location over the length limit is rejected", func(t *testing.T) {
		location := maybe.NewJust(strings.Repeat("a", maxSetlistEventLocationLength+1))
		_, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", location, maybe.NewAbsent[time.Time](), nil)
		assert.ErrorIs(t, err, ErrSetlistEventLocationTooLong)
	})

	t.Run("too many tracks are rejected", func(t *testing.T) {
		trackIDs := make([]TrackID, maxSetlistTracks+1)
		for i := range trackIDs {
			trackIDs[i] = uuid.Generate()
		}
		_, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), trackIDs)
		assert.ErrorIs(t, err, ErrTooManySetlistTracks)
	})

	t.Run("duplicate track ids are rejected", func(t *testing.T) {
		trackID := uuid.Generate()
		_, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), []TrackID{trackID, trackID})
		assert.ErrorIs(t, err, ErrDuplicateTrackIDs)
	})

	t.Run("mutating the returned track slice does not affect the setlist", func(t *testing.T) {
		trackID := uuid.Generate()
		setlist, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), []TrackID{trackID})
		require.NoError(t, err)

		trackIDs := setlist.TrackIDs()
		trackIDs[0] = uuid.Generate()

		assert.Equal(t, []TrackID{trackID}, setlist.TrackIDs())
	})
}

func TestSetlist_AddTrack(t *testing.T) {
	t.Run("appends a new track", func(t *testing.T) {
		setlist, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), nil)
		require.NoError(t, err)

		trackID := uuid.Generate()
		require.NoError(t, setlist.AddTrack(trackID))

		assert.Equal(t, []TrackID{trackID}, setlist.TrackIDs())
	})

	t.Run("adding an existing track is rejected", func(t *testing.T) {
		trackID := uuid.Generate()
		setlist, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), []TrackID{trackID})
		require.NoError(t, err)

		err = setlist.AddTrack(trackID)
		assert.ErrorIs(t, err, ErrTrackAlreadyInSetlist)
	})

	t.Run("adding a track beyond the limit is rejected", func(t *testing.T) {
		trackIDs := make([]TrackID, maxSetlistTracks)
		for i := range trackIDs {
			trackIDs[i] = uuid.Generate()
		}
		setlist, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), trackIDs)
		require.NoError(t, err)

		err = setlist.AddTrack(uuid.Generate())
		assert.ErrorIs(t, err, ErrTooManySetlistTracks)
	})
}

func TestSetlist_RemoveTrack(t *testing.T) {
	t.Run("removes an existing track", func(t *testing.T) {
		trackA, trackB := uuid.Generate(), uuid.Generate()
		setlist, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), []TrackID{trackA, trackB})
		require.NoError(t, err)

		require.NoError(t, setlist.RemoveTrack(trackA))

		assert.Equal(t, []TrackID{trackB}, setlist.TrackIDs())
	})

	t.Run("removing a track not in the setlist is rejected", func(t *testing.T) {
		setlist, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), nil)
		require.NoError(t, err)

		err = setlist.RemoveTrack(uuid.Generate())
		assert.ErrorIs(t, err, ErrTrackNotInSetlist)
	})
}

func TestSetlist_RemoveTracks(t *testing.T) {
	t.Run("removes multiple existing tracks", func(t *testing.T) {
		trackA, trackB, trackC := uuid.Generate(), uuid.Generate(), uuid.Generate()
		setlist, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), []TrackID{trackA, trackB, trackC})
		require.NoError(t, err)

		require.NoError(t, setlist.RemoveTracks([]TrackID{trackA, trackC}))

		assert.Equal(t, []TrackID{trackB}, setlist.TrackIDs())
	})

	t.Run("one id not in the setlist rejects the whole batch", func(t *testing.T) {
		trackA := uuid.Generate()
		setlist, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), []TrackID{trackA})
		require.NoError(t, err)

		err = setlist.RemoveTracks([]TrackID{trackA, uuid.Generate()})
		require.ErrorIs(t, err, ErrTrackNotInSetlist)
		assert.Equal(t, []TrackID{trackA}, setlist.TrackIDs())
	})
}

func TestSetlist_Reorder(t *testing.T) {
	t.Run("reorders tracks to the given order", func(t *testing.T) {
		trackA, trackB, trackC := uuid.Generate(), uuid.Generate(), uuid.Generate()
		setlist, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), []TrackID{trackA, trackB, trackC})
		require.NoError(t, err)

		require.NoError(t, setlist.Reorder([]TrackID{trackC, trackA, trackB}))

		assert.Equal(t, []TrackID{trackC, trackA, trackB}, setlist.TrackIDs())
	})

	t.Run("empty setlist accepts empty order", func(t *testing.T) {
		setlist, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), nil)
		require.NoError(t, err)

		assert.NoError(t, setlist.Reorder(nil))
	})

	t.Run("order with a different length is rejected", func(t *testing.T) {
		trackA, trackB := uuid.Generate(), uuid.Generate()
		setlist, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), []TrackID{trackA, trackB})
		require.NoError(t, err)

		err = setlist.Reorder([]TrackID{trackA})
		assert.ErrorIs(t, err, ErrInvalidSetlistTrackOrder)
	})

	t.Run("order with a track from outside the setlist is rejected", func(t *testing.T) {
		trackA, trackB := uuid.Generate(), uuid.Generate()
		setlist, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), []TrackID{trackA})
		require.NoError(t, err)

		err = setlist.Reorder([]TrackID{trackB})
		assert.ErrorIs(t, err, ErrInvalidSetlistTrackOrder)
	})

	t.Run("order with a duplicated track id is rejected", func(t *testing.T) {
		trackA, trackB := uuid.Generate(), uuid.Generate()
		setlist, err := NewSetlist(uuid.Generate(), uuid.Generate(), "Setlist", maybe.NewAbsent[string](), maybe.NewAbsent[time.Time](), []TrackID{trackA, trackB})
		require.NoError(t, err)

		err = setlist.Reorder([]TrackID{trackA, trackA})
		assert.ErrorIs(t, err, ErrInvalidSetlistTrackOrder)
	})
}

func TestLoadSetlist(t *testing.T) {
	id := uuid.Generate()
	bandID := uuid.Generate()
	trackID := uuid.Generate()

	setlist := LoadSetlist(id, bandID, "Loaded", maybe.NewJust("Venue"), maybe.NewJust(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)), []TrackID{trackID})

	assert.Equal(t, id, setlist.ID())
	assert.Equal(t, bandID, setlist.BandID())
	assert.Equal(t, "Loaded", setlist.Name())
	location, ok := maybe.JustValid(setlist.EventLocation())
	require.True(t, ok)
	assert.Equal(t, "Venue", location)
	assert.Equal(t, []TrackID{trackID}, setlist.TrackIDs())
}
