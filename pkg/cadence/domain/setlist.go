package domain

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"cadence/pkg/common/maybe"
)

const (
	maxSetlistNameLength          = 255
	maxSetlistEventLocationLength = 255
	maxSetlistTracks              = 100
)

var (
	ErrEmptySetlistName            = errors.New("setlist name can not be empty")
	ErrSetlistNameTooLong          = fmt.Errorf("setlist name length should be less than or equal to %d", maxSetlistNameLength)
	ErrSetlistEventLocationTooLong = fmt.Errorf("setlist event location length should be less than or equal to %d", maxSetlistEventLocationLength)

	ErrSetlistNotFound = errors.New("setlist not found")

	ErrTrackAlreadyInSetlist    = errors.New("track is already in the setlist")
	ErrTrackNotInSetlist        = errors.New("track is not in the setlist")
	ErrDuplicateTrackIDs        = errors.New("track ids must not contain duplicates")
	ErrInvalidSetlistTrackOrder = errors.New("track order must contain exactly the tracks currently in the setlist")
	ErrTooManySetlistTracks     = fmt.Errorf("setlist can not contain more than %d tracks", maxSetlistTracks)
)

type Setlist struct {
	id            SetlistID
	bandID        BandID
	name          string
	eventLocation maybe.Maybe[string]
	eventDate     maybe.Maybe[time.Time]
	trackIDs      []TrackID
}

type SetlistRepository interface {
	NextID() SetlistID
	Store(*Setlist) error
	Get(SetlistID) (*Setlist, error)
	Remove(SetlistID) error
}

func NewSetlist(
	id SetlistID,
	bandID BandID,
	name string,
	eventLocation maybe.Maybe[string],
	eventDate maybe.Maybe[time.Time],
	trackIDs []TrackID,
) (*Setlist, error) {
	if err := validateSetlistNameLength(name); err != nil {
		return nil, err
	}
	if err := validateEventLocationLength(eventLocation); err != nil {
		return nil, err
	}
	if containsDuplicate(trackIDs) {
		return nil, ErrDuplicateTrackIDs
	}
	if len(trackIDs) > maxSetlistTracks {
		return nil, ErrTooManySetlistTracks
	}
	return &Setlist{
		id:            id,
		bandID:        bandID,
		name:          name,
		eventLocation: eventLocation,
		eventDate:     eventDate,
		trackIDs:      slices.Clone(trackIDs),
	}, nil
}

func LoadSetlist(
	id SetlistID,
	bandID BandID,
	name string,
	eventLocation maybe.Maybe[string],
	eventDate maybe.Maybe[time.Time],
	trackIDs []TrackID,
) *Setlist {
	return &Setlist{
		id:            id,
		bandID:        bandID,
		name:          name,
		eventLocation: eventLocation,
		eventDate:     eventDate,
		trackIDs:      slices.Clone(trackIDs),
	}
}

func (s *Setlist) ID() SetlistID {
	return s.id
}

func (s *Setlist) BandID() BandID {
	return s.bandID
}

func (s *Setlist) Name() string {
	return s.name
}

func (s *Setlist) EventLocation() maybe.Maybe[string] {
	return s.eventLocation
}

func (s *Setlist) EventDate() maybe.Maybe[time.Time] {
	return s.eventDate
}

func (s *Setlist) TrackIDs() []TrackID {
	return slices.Clone(s.trackIDs)
}

func (s *Setlist) SetName(name string) error {
	if err := validateSetlistNameLength(name); err != nil {
		return err
	}
	s.name = name
	return nil
}

func (s *Setlist) SetEventLocation(eventLocation maybe.Maybe[string]) error {
	if err := validateEventLocationLength(eventLocation); err != nil {
		return err
	}
	s.eventLocation = eventLocation
	return nil
}

func (s *Setlist) SetEventDate(eventDate maybe.Maybe[time.Time]) {
	s.eventDate = eventDate
}

func (s *Setlist) AddTrack(trackID TrackID) error {
	if slices.Contains(s.trackIDs, trackID) {
		return ErrTrackAlreadyInSetlist
	}
	if len(s.trackIDs) >= maxSetlistTracks {
		return ErrTooManySetlistTracks
	}
	s.trackIDs = append(s.trackIDs, trackID)
	return nil
}

func (s *Setlist) RemoveTrack(trackID TrackID) error {
	for i, id := range s.trackIDs {
		if id == trackID {
			s.trackIDs = append(s.trackIDs[:i], s.trackIDs[i+1:]...)
			return nil
		}
	}
	return ErrTrackNotInSetlist
}

func (s *Setlist) Reorder(trackIDs []TrackID) error {
	if len(trackIDs) != len(s.trackIDs) || containsDuplicate(trackIDs) {
		return ErrInvalidSetlistTrackOrder
	}

	current := make(map[TrackID]struct{}, len(s.trackIDs))
	for _, id := range s.trackIDs {
		current[id] = struct{}{}
	}
	for _, id := range trackIDs {
		if _, ok := current[id]; !ok {
			return ErrInvalidSetlistTrackOrder
		}
	}

	s.trackIDs = slices.Clone(trackIDs)
	return nil
}

func validateSetlistNameLength(name string) error {
	return checkStringLimits(name, maxSetlistNameLength, ErrEmptySetlistName, ErrSetlistNameTooLong)
}

func validateEventLocationLength(eventLocation maybe.Maybe[string]) error {
	value, ok := maybe.JustValid(eventLocation)
	if !ok || len(value) <= maxSetlistEventLocationLength {
		return nil
	}
	return ErrSetlistEventLocationTooLong
}

func containsDuplicate(trackIDs []TrackID) bool {
	seen := make(map[TrackID]struct{}, len(trackIDs))
	for _, id := range trackIDs {
		if _, dup := seen[id]; dup {
			return true
		}
		seen[id] = struct{}{}
	}
	return false
}
