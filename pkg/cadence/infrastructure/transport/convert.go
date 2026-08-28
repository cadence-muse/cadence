package transport

import (
	"time"

	googleuuid "github.com/google/uuid"
	"github.com/nightnoryu/go-kita/maybe"
	"github.com/nightnoryu/go-kita/slices"

	"cadence/api/server/publicapi"
	"cadence/pkg/cadence/app/query"
	"cadence/pkg/common/valuetypes"
)

func convertQueryBandListItemToAPI(band query.BandListItem) publicapi.BandListItem {
	return publicapi.BandListItem{
		ID:           googleuuid.UUID(band.ID),
		Name:         band.Name,
		MembersCount: band.MembersCount,
		OwnerId:      publicapi.NewOptUUID(googleuuid.UUID(band.OwnerID)),
	}
}

func convertQueryBandDataToAPI(band query.BandData) publicapi.Band {
	return publicapi.Band{
		ID:         googleuuid.UUID(band.ID),
		Name:       band.Name,
		OwnerId:    googleuuid.UUID(band.OwnerID),
		InviteCode: band.InviteCode,
		Members:    slices.Map(band.Members, convertQueryBandMemberDataToAPI),
	}
}

func convertQueryBandMemberDataToAPI(member query.BandMemberData) publicapi.BandMember {
	return publicapi.BandMember{
		ID:       googleuuid.UUID(member.ID),
		Username: member.Username,
		Role:     publicapi.BandMemberRole(member.Role),
	}
}

func convertQueryTrackListItemToAPI(track query.TrackListItem) publicapi.TrackListItem {
	return publicapi.TrackListItem{
		ID:              googleuuid.UUID(track.ID),
		Title:           track.Title,
		Artist:          track.Artist,
		DurationSeconds: optIntFromDuration(track.Duration),
		Key:             optStringFromMaybe(track.Key),
	}
}

func convertQueryUserTrackListItemToAPI(track query.UserTrackListItem) publicapi.UserTrackListItem {
	return publicapi.UserTrackListItem{
		ID:              googleuuid.UUID(track.ID),
		Title:           track.Title,
		Artist:          track.Artist,
		DurationSeconds: optIntFromDuration(track.Duration),
		BandId:          googleuuid.UUID(track.BandID),
		BandName:        track.BandName,
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
		EventLocation:   optStringFromMaybe(setlist.EventLocation),
	}
}

func convertQueryUserSetlistListItemToAPI(setlist query.UserSetlistListItem) publicapi.UserSetlistListItem {
	return publicapi.UserSetlistListItem{
		ID:              googleuuid.UUID(setlist.ID),
		Name:            setlist.Name,
		TracksCount:     setlist.TracksCount,
		DurationSeconds: durationToIntSeconds(setlist.Duration),
		EventDate:       optDateFromMaybe(setlist.EventDate),
		BandId:          googleuuid.UUID(setlist.BandID),
		BandName:        setlist.BandName,
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

func parseTempo(value int) (maybe.Maybe[valuetypes.Tempo], error) {
	tempo, err := valuetypes.MakeTempo(value)
	if err != nil {
		return maybe.NewNone[valuetypes.Tempo](), err
	}
	return maybe.NewJust(tempo), nil
}

func intSecondsToDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

func durationToIntSeconds(duration time.Duration) int {
	return int(duration / time.Second)
}
