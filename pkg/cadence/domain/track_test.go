package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cadence/pkg/common/maybe"
	"cadence/pkg/common/uuid"
	"cadence/pkg/common/valuetypes"
)

func TestNewTrack(t *testing.T) {
	t.Run("valid track is created", func(t *testing.T) {
		id := uuid.Generate()
		bandID := uuid.Generate()
		key, err := valuetypes.MakeKey("C")
		require.NoError(t, err)

		track, err := NewTrack(
			id, bandID, "Song", "Artist",
			maybe.NewJust(3*time.Minute),
			maybe.NewJust(120),
			maybe.NewJust(key),
			maybe.NewJust("notes"),
		)
		require.NoError(t, err)
		assert.Equal(t, id, track.ID())
		assert.Equal(t, bandID, track.BandID())
		assert.Equal(t, "Song", track.Title())
		assert.Equal(t, "Artist", track.Artist())

		duration, ok := maybe.JustValid(track.Duration())
		require.True(t, ok)
		assert.Equal(t, 3*time.Minute, duration)

		tempo, ok := maybe.JustValid(track.Tempo())
		require.True(t, ok)
		assert.Equal(t, 120, tempo)

		notes, ok := maybe.JustValid(track.Notes())
		require.True(t, ok)
		assert.Equal(t, "notes", notes)
	})

	t.Run("optional fields can be absent", func(t *testing.T) {
		track, err := NewTrack(
			uuid.Generate(), uuid.Generate(), "Song", "Artist",
			maybe.NewAbsent[time.Duration](),
			maybe.NewAbsent[int](),
			maybe.NewAbsent[valuetypes.MusicalKey](),
			maybe.NewAbsent[string](),
		)
		require.NoError(t, err)
		_, ok := maybe.JustValid(track.Duration())
		assert.False(t, ok)
	})

	t.Run("empty title is rejected", func(t *testing.T) {
		_, err := NewTrack(
			uuid.Generate(), uuid.Generate(), "", "Artist",
			maybe.NewAbsent[time.Duration](), maybe.NewAbsent[int](),
			maybe.NewAbsent[valuetypes.MusicalKey](), maybe.NewAbsent[string](),
		)
		assert.ErrorIs(t, err, ErrEmptyTrackTitle)
	})

	t.Run("title over the length limit is rejected", func(t *testing.T) {
		_, err := NewTrack(
			uuid.Generate(), uuid.Generate(), strings.Repeat("a", maxTrackTitleLength+1), "Artist",
			maybe.NewAbsent[time.Duration](), maybe.NewAbsent[int](),
			maybe.NewAbsent[valuetypes.MusicalKey](), maybe.NewAbsent[string](),
		)
		assert.ErrorIs(t, err, ErrTrackTitleTooLong)
	})

	t.Run("title at the length limit is accepted", func(t *testing.T) {
		_, err := NewTrack(
			uuid.Generate(), uuid.Generate(), strings.Repeat("a", maxTrackTitleLength), "Artist",
			maybe.NewAbsent[time.Duration](), maybe.NewAbsent[int](),
			maybe.NewAbsent[valuetypes.MusicalKey](), maybe.NewAbsent[string](),
		)
		assert.NoError(t, err)
	})

	t.Run("empty artist is rejected", func(t *testing.T) {
		_, err := NewTrack(
			uuid.Generate(), uuid.Generate(), "Song", "",
			maybe.NewAbsent[time.Duration](), maybe.NewAbsent[int](),
			maybe.NewAbsent[valuetypes.MusicalKey](), maybe.NewAbsent[string](),
		)
		assert.ErrorIs(t, err, ErrEmptyTrackArtist)
	})

	t.Run("artist over the length limit is rejected", func(t *testing.T) {
		_, err := NewTrack(
			uuid.Generate(), uuid.Generate(), "Song", strings.Repeat("a", maxTrackArtistLength+1),
			maybe.NewAbsent[time.Duration](), maybe.NewAbsent[int](),
			maybe.NewAbsent[valuetypes.MusicalKey](), maybe.NewAbsent[string](),
		)
		assert.ErrorIs(t, err, ErrTrackArtistTooLong)
	})
}

func newTestTrack(t *testing.T) *Track {
	t.Helper()
	track, err := NewTrack(
		uuid.Generate(), uuid.Generate(), "Song", "Artist",
		maybe.NewAbsent[time.Duration](), maybe.NewAbsent[int](),
		maybe.NewAbsent[valuetypes.MusicalKey](), maybe.NewAbsent[string](),
	)
	require.NoError(t, err)
	return track
}

func TestLoadTrack(t *testing.T) {
	id := uuid.Generate()
	bandID := uuid.Generate()

	track := LoadTrack(
		id, bandID, "Loaded", "Artist",
		maybe.NewJust(2*time.Minute),
		maybe.NewJust(90),
		maybe.NewAbsent[valuetypes.MusicalKey](),
		maybe.NewAbsent[string](),
	)

	assert.Equal(t, id, track.ID())
	assert.Equal(t, bandID, track.BandID())
	assert.Equal(t, "Loaded", track.Title())
	assert.Equal(t, "Artist", track.Artist())
	duration, ok := maybe.JustValid(track.Duration())
	require.True(t, ok)
	assert.Equal(t, 2*time.Minute, duration)
}

func TestTrack_SetTitle(t *testing.T) {
	t.Run("valid title is set", func(t *testing.T) {
		track := newTestTrack(t)

		require.NoError(t, track.SetTitle("New Title"))
		assert.Equal(t, "New Title", track.Title())
	})

	t.Run("empty title is rejected and leaves title unchanged", func(t *testing.T) {
		track := newTestTrack(t)

		err := track.SetTitle("")
		require.ErrorIs(t, err, ErrEmptyTrackTitle)
		assert.Equal(t, "Song", track.Title())
	})
}

func TestTrack_SetArtist(t *testing.T) {
	t.Run("valid artist is set", func(t *testing.T) {
		track := newTestTrack(t)

		require.NoError(t, track.SetArtist("New Artist"))
		assert.Equal(t, "New Artist", track.Artist())
	})

	t.Run("empty artist is rejected and leaves artist unchanged", func(t *testing.T) {
		track := newTestTrack(t)

		err := track.SetArtist("")
		require.ErrorIs(t, err, ErrEmptyTrackArtist)
		assert.Equal(t, "Artist", track.Artist())
	})
}

func TestTrack_SetDuration(t *testing.T) {
	track := newTestTrack(t)

	track.SetDuration(maybe.NewJust(5 * time.Minute))

	duration, ok := maybe.JustValid(track.Duration())
	require.True(t, ok)
	assert.Equal(t, 5*time.Minute, duration)
}

func TestTrack_SetTempo(t *testing.T) {
	track := newTestTrack(t)

	track.SetTempo(maybe.NewJust(140))

	tempo, ok := maybe.JustValid(track.Tempo())
	require.True(t, ok)
	assert.Equal(t, 140, tempo)
}

func TestTrack_SetKey(t *testing.T) {
	track := newTestTrack(t)

	key, err := valuetypes.MakeKey("Am")
	require.NoError(t, err)
	track.SetKey(maybe.NewJust(key))

	gotKey, ok := maybe.JustValid(track.Key())
	require.True(t, ok)
	assert.Equal(t, key, gotKey)
}

func TestTrack_SetNotes(t *testing.T) {
	track := newTestTrack(t)

	track.SetNotes(maybe.NewJust("updated notes"))

	notes, ok := maybe.JustValid(track.Notes())
	require.True(t, ok)
	assert.Equal(t, "updated notes", notes)
}
