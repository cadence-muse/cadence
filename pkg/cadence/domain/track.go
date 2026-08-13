package domain

import (
	"errors"
	"fmt"
	"time"

	"cadence/pkg/common/maybe"
	"cadence/pkg/common/valuetypes"
)

const (
	maxTrackTitleLength  = 255
	maxTrackArtistLength = 255
)

var (
	ErrEmptyTrackTitle   = errors.New("track title can not be empty")
	ErrTrackTitleTooLong = fmt.Errorf("track title length should be less than or equal to %d", maxTrackTitleLength)

	ErrEmptyTrackArtist   = errors.New("track artist can not be empty")
	ErrTrackArtistTooLong = fmt.Errorf("track artist length should be less than or equal to %d", maxTrackArtistLength)

	ErrTrackNotFound = errors.New("track not found")
)

type Track struct {
	id       TrackID
	bandID   BandID
	title    string
	artist   string
	duration maybe.Maybe[time.Duration]
	tempo    maybe.Maybe[int]
	key      maybe.Maybe[valuetypes.MusicalKey]
	notes    maybe.Maybe[string]
}

type TrackRepository interface {
	NextID() TrackID
	Store(*Track) error
	Get(TrackID) (*Track, error)
}

func NewTrack(
	id TrackID,
	bandID BandID,
	title string,
	artist string,
	duration maybe.Maybe[time.Duration],
	tempo maybe.Maybe[int],
	key maybe.Maybe[valuetypes.MusicalKey],
	notes maybe.Maybe[string],
) (*Track, error) {
	err := validateTrackTitleLength(title)
	if err != nil {
		return nil, err
	}
	err = validateTrackArtistLength(artist)
	if err != nil {
		return nil, err
	}
	return &Track{
		id:       id,
		bandID:   bandID,
		title:    title,
		artist:   artist,
		duration: duration,
		tempo:    tempo,
		key:      key,
		notes:    notes,
	}, nil
}

func LoadTrack(
	id TrackID,
	bandID BandID,
	title string,
	artist string,
	duration maybe.Maybe[time.Duration],
	tempo maybe.Maybe[int],
	key maybe.Maybe[valuetypes.MusicalKey],
	notes maybe.Maybe[string],
) *Track {
	return &Track{
		id:       id,
		bandID:   bandID,
		title:    title,
		artist:   artist,
		duration: duration,
		tempo:    tempo,
		key:      key,
		notes:    notes,
	}
}

func (t *Track) ID() TrackID {
	return t.id
}

func (t *Track) BandID() BandID {
	return t.bandID
}

func (t *Track) Title() string {
	return t.title
}

func (t *Track) Artist() string {
	return t.artist
}

func (t *Track) Duration() maybe.Maybe[time.Duration] {
	return t.duration
}

func (t *Track) Tempo() maybe.Maybe[int] {
	return t.tempo
}

func (t *Track) Key() maybe.Maybe[valuetypes.MusicalKey] {
	return t.key
}

func (t *Track) Notes() maybe.Maybe[string] {
	return t.notes
}

func validateTrackTitleLength(title string) error {
	return checkStringLimits(title, maxTrackTitleLength, ErrEmptyTrackTitle, ErrTrackTitleTooLong)
}

func validateTrackArtistLength(artist string) error {
	return checkStringLimits(artist, maxTrackArtistLength, ErrEmptyTrackArtist, ErrTrackArtistTooLong)
}
