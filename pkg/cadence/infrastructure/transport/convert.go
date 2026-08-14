package transport

import (
	"time"

	googleuuid "github.com/google/uuid"

	"cadence/api/server/publicapi"
	"cadence/pkg/cadence/app/query"
	"cadence/pkg/common/maybe"
	"cadence/pkg/common/slices"
	"cadence/pkg/common/valuetypes"
)

func convertQueryBandListItemToAPI(band query.BandListItem) publicapi.BandListItem {
	return publicapi.BandListItem{
		ID:   googleuuid.UUID(band.ID),
		Name: band.Name,
	}
}

func convertQueryBandDataToAPI(band query.BandData) publicapi.Band {
	return publicapi.Band{
		ID:         googleuuid.UUID(band.ID),
		Name:       band.Name,
		InviteCode: band.InviteCode,
	}
}

func convertQueryTrackListItemToAPI(track query.TrackListItem) publicapi.TrackListItem {
	return publicapi.TrackListItem{
		ID:              googleuuid.UUID(track.ID),
		Title:           track.Title,
		Artist:          track.Artist,
		DurationSeconds: optIntFromDuration(track.Duration),
	}
}

func convertQueryTrackDataToAPI(track query.TrackData) publicapi.BandTrack {
	return publicapi.BandTrack{
		ID:              googleuuid.UUID(track.ID),
		Title:           track.Title,
		Artist:          track.Artist,
		DurationSeconds: optIntFromDuration(track.Duration),
		Tempo:           optIntFromMaybe(track.Tempo),
		Key:             optStringFromMaybe(track.Key),
		Notes:           optStringFromMaybe(track.Notes),
	}
}

func convertQuerySetlistListItemToAPI(setlist query.SetlistListItem) publicapi.SetlistListItem {
	return publicapi.SetlistListItem{
		ID:        googleuuid.UUID(setlist.ID),
		Name:      setlist.Name,
		EventDate: optDateFromMaybe(setlist.EventDate),
	}
}

func convertQuerySetlistDataToAPI(setlist query.SetlistData) publicapi.BandSetlist {
	return publicapi.BandSetlist{
		ID:            googleuuid.UUID(setlist.ID),
		Name:          setlist.Name,
		EventLocation: optStringFromMaybe(setlist.EventLocation),
		EventDate:     optDateFromMaybe(setlist.EventDate),
		Tracks:        slices.Map(setlist.Tracks, convertQuerySetlistTrackItemToAPI),
	}
}

func convertQuerySetlistTrackItemToAPI(track query.SetlistTrackItem) publicapi.SetlistTrackItem {
	return publicapi.SetlistTrackItem{
		TrackId:         googleuuid.UUID(track.TrackID),
		Position:        track.Position,
		Title:           track.Title,
		Artist:          track.Artist,
		DurationSeconds: optIntFromDuration(track.Duration),
	}
}

func optDateFromMaybe(v maybe.Maybe[time.Time]) publicapi.OptDate {
	value, ok := maybe.JustValid(v)
	if !ok {
		return publicapi.OptDate{}
	}
	return publicapi.NewOptDate(value)
}

func optIntFromDuration(duration maybe.Maybe[time.Duration]) publicapi.OptInt {
	value, ok := maybe.JustValid(duration)
	if !ok {
		return publicapi.OptInt{}
	}
	return publicapi.NewOptInt(int(value / time.Second))
}

func optIntFromMaybe(v maybe.Maybe[int]) publicapi.OptInt {
	value, ok := maybe.JustValid(v)
	if !ok {
		return publicapi.OptInt{}
	}
	return publicapi.NewOptInt(value)
}

func optStringFromMaybe(v maybe.Maybe[string]) publicapi.OptString {
	value, ok := maybe.JustValid(v)
	if !ok {
		return publicapi.OptString{}
	}
	return publicapi.NewOptString(value)
}

func maybeDurationFromOpt(opt publicapi.OptInt) maybe.Maybe[time.Duration] {
	if !opt.Set {
		return maybe.NewAbsent[time.Duration]()
	}
	return maybe.NewJust(time.Duration(opt.Value) * time.Second)
}

func maybeKeyFromOpt(opt publicapi.OptString) (maybe.Maybe[valuetypes.MusicalKey], error) {
	if !opt.Set {
		return maybe.NewAbsent[valuetypes.MusicalKey](), nil
	}
	return parseMusicalKey(opt.Value)
}

func maybeDurationFromOptNil(opt publicapi.OptNilInt) maybe.Maybe[time.Duration] {
	if !opt.Set {
		return maybe.NewAbsent[time.Duration]()
	}
	if opt.Null {
		return maybe.NewNone[time.Duration]()
	}
	return maybe.NewJust(time.Duration(opt.Value) * time.Second)
}

func maybeKeyFromOptNil(opt publicapi.OptNilString) (maybe.Maybe[valuetypes.MusicalKey], error) {
	if !opt.Set {
		return maybe.NewAbsent[valuetypes.MusicalKey](), nil
	}
	if opt.Null {
		return maybe.NewNone[valuetypes.MusicalKey](), nil
	}
	return parseMusicalKey(opt.Value)
}

func parseMusicalKey(value string) (maybe.Maybe[valuetypes.MusicalKey], error) {
	key, err := valuetypes.MakeKey(value)
	if err != nil {
		return maybe.NewNone[valuetypes.MusicalKey](), err
	}
	return maybe.NewJust(key), nil
}
