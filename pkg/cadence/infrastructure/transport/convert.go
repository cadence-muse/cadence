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
		OwnerId:    googleuuid.UUID(band.OwnerID),
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
		ID:              googleuuid.UUID(setlist.ID),
		Name:            setlist.Name,
		TracksCount:     setlist.TracksCount,
		DurationSeconds: durationToIntSeconds(setlist.Duration),
		EventDate:       optDateFromMaybe(setlist.EventDate),
	}
}

func convertQuerySetlistDataToAPI(setlist query.SetlistData) publicapi.BandSetlist {
	return publicapi.BandSetlist{
		ID:              googleuuid.UUID(setlist.ID),
		Name:            setlist.Name,
		DurationSeconds: durationToIntSeconds(setlist.Duration),
		EventLocation:   optStringFromMaybe(setlist.EventLocation),
		EventDate:       optDateFromMaybe(setlist.EventDate),
		Tracks:          slices.Map(setlist.Tracks, convertQuerySetlistTrackItemToAPI),
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

func parseMusicalKey(value string) (maybe.Maybe[valuetypes.MusicalKey], error) {
	key, err := valuetypes.MakeKey(value)
	if err != nil {
		return maybe.NewNone[valuetypes.MusicalKey](), err
	}
	return maybe.NewJust(key), nil
}

func intSecondsToDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

func durationToIntSeconds(duration time.Duration) int {
	return int(duration / time.Second)
}
