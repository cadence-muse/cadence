//go:build e2e

package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"cadence/api/server/publicapi"
)

func TestUserJourney(t *testing.T) {
	env := setupTestEnv(t)
	ctx := context.Background()

	const (
		ownerUsername  = "alice"
		ownerPassword  = "correct-horse-battery-staple"
		memberUsername = "bob"
		memberPassword = "another-horse-battery-staple"
	)

	var (
		ownerID                uuid.UUID
		memberID               uuid.UUID
		ownerToken             string
		memberToken            string
		bandID                 uuid.UUID
		inviteCode             string
		track1, track2, track3 uuid.UUID
		track4, track5         uuid.UUID
		setlistID              uuid.UUID
	)

	t.Run("register owner", func(t *testing.T) {
		res, err := env.client.Register(ctx, &publicapi.RegisterRequestBody{
			Username: ownerUsername,
			Password: ownerPassword,
		})
		body := requireResponse[publicapi.RegisterResponseBody](t, res, err)
		require.NotEqual(t, uuid.Nil, body.ID)
		ownerID = body.ID
	})

	t.Run("login owner", func(t *testing.T) {
		res, err := env.client.Login(ctx, &publicapi.LoginRequestBody{
			Username: ownerUsername,
			Password: ownerPassword,
		})
		body := requireResponse[publicapi.LoginResponseBody](t, res, err)
		require.NotEmpty(t, body.Token)
		ownerToken = body.Token
		env.sec.token = ownerToken
	})

	t.Run("get owner profile", func(t *testing.T) {
		res, err := env.client.GetUserProfile(ctx)
		body := requireResponse[publicapi.UserProfile](t, res, err)
		require.Equal(t, ownerID, body.ID)
		require.Equal(t, ownerUsername, body.Username)
	})

	t.Run("homepage data before any band", func(t *testing.T) {
		res, err := env.client.GetHomepageData(ctx)
		body := requireResponse[publicapi.HomepageData](t, res, err)
		require.Equal(t, ownerUsername, body.Username)
		require.Equal(t, 0, body.BandsCount)
	})

	t.Run("create band", func(t *testing.T) {
		res, err := env.client.CreateBand(ctx, &publicapi.CreateBandRequestBody{
			Name: "The Wanderers",
		})
		body := requireResponse[publicapi.CreateBandResponseBody](t, res, err)
		require.NotEqual(t, uuid.Nil, body.ID)
		bandID = body.ID
	})

	t.Run("get band after creation", func(t *testing.T) {
		res, err := env.client.GetBand(ctx, publicapi.GetBandParams{BandId: bandID})
		body := requireResponse[publicapi.Band](t, res, err)
		require.Equal(t, bandID, body.ID)
		require.Equal(t, "The Wanderers", body.Name)
		require.Equal(t, ownerID, body.OwnerId)
		require.NotEmpty(t, body.InviteCode)
		inviteCode = body.InviteCode
	})

	t.Run("list bands contains the new band", func(t *testing.T) {
		res, err := env.client.ListBands(ctx)
		body := requireResponse[publicapi.ListBandsResponseBody](t, res, err)
		require.Len(t, body.Items, 1)
		require.Equal(t, bandID, body.Items[0].ID)
		require.Equal(t, "The Wanderers", body.Items[0].Name)
	})

	t.Run("homepage data reflects the new band", func(t *testing.T) {
		res, err := env.client.GetHomepageData(ctx)
		body := requireResponse[publicapi.HomepageData](t, res, err)
		require.Equal(t, 1, body.BandsCount)
	})

	t.Run("update band name", func(t *testing.T) {
		res, err := env.client.UpdateBand(ctx, &publicapi.UpdateBandRequestBody{
			Name: publicapi.NewOptString("The Wanderers Reborn"),
		}, publicapi.UpdateBandParams{BandId: bandID})
		requireResponse[publicapi.UpdateBandOK](t, res, err)

		getRes, getErr := env.client.GetBand(ctx, publicapi.GetBandParams{BandId: bandID})
		body := requireResponse[publicapi.Band](t, getRes, getErr)
		require.Equal(t, "The Wanderers Reborn", body.Name)
		require.Equal(t, ownerID, body.OwnerId)
	})

	t.Run("register and login member", func(t *testing.T) {
		registerRes, registerErr := env.client.Register(ctx, &publicapi.RegisterRequestBody{
			Username: memberUsername,
			Password: memberPassword,
		})
		registerBody := requireResponse[publicapi.RegisterResponseBody](t, registerRes, registerErr)
		memberID = registerBody.ID

		loginRes, loginErr := env.client.Login(ctx, &publicapi.LoginRequestBody{
			Username: memberUsername,
			Password: memberPassword,
		})
		loginBody := requireResponse[publicapi.LoginResponseBody](t, loginRes, loginErr)
		memberToken = loginBody.Token
	})

	t.Run("member joins band by invite code", func(t *testing.T) {
		env.sec.token = memberToken
		defer func() { env.sec.token = ownerToken }()

		joinRes, joinErr := env.client.JoinBand(ctx, &publicapi.JoinBandRequestBody{
			InviteCode: inviteCode,
		})
		requireResponse[publicapi.JoinBandOK](t, joinRes, joinErr)

		listRes, listErr := env.client.ListBands(ctx)
		body := requireResponse[publicapi.ListBandsResponseBody](t, listRes, listErr)
		require.Len(t, body.Items, 1)
		require.Equal(t, bandID, body.Items[0].ID)
	})

	t.Run("create band tracks", func(t *testing.T) {
		res1, err1 := env.client.CreateBandTrack(ctx, &publicapi.CreateBandTrackRequestBody{
			Title:           "Track One",
			Artist:          "Artist One",
			DurationSeconds: publicapi.NewOptInt(180),
			Tempo:           publicapi.NewOptInt(120),
			Key:             publicapi.NewOptString("C"),
			Notes:           publicapi.NewOptString("first track"),
		}, publicapi.CreateBandTrackParams{BandId: bandID})
		track1Body := requireResponse[publicapi.CreateBandTrackResponseBody](t, res1, err1)
		track1 = track1Body.ID

		res2, err2 := env.client.CreateBandTrack(ctx, &publicapi.CreateBandTrackRequestBody{
			Title:           "Track Two",
			Artist:          "Artist Two",
			DurationSeconds: publicapi.NewOptInt(200),
		}, publicapi.CreateBandTrackParams{BandId: bandID})
		track2Body := requireResponse[publicapi.CreateBandTrackResponseBody](t, res2, err2)
		track2 = track2Body.ID
	})

	t.Run("list and get band tracks", func(t *testing.T) {
		listRes, listErr := env.client.ListBandTracks(ctx, publicapi.ListBandTracksParams{BandId: bandID})
		listBody := requireResponse[publicapi.ListBandTracksResponseBody](t, listRes, listErr)
		require.Len(t, listBody.Items, 2)

		getRes, getErr := env.client.GetBandTrack(ctx, publicapi.GetBandTrackParams{BandId: bandID, TrackId: track1})
		getBody := requireResponse[publicapi.BandTrack](t, getRes, getErr)
		require.Equal(t, "Track One", getBody.Title)
		require.Equal(t, "Artist One", getBody.Artist)
		seconds, ok := getBody.DurationSeconds.Get()
		require.True(t, ok)
		require.Equal(t, 180, seconds)
	})

	t.Run("update a band track", func(t *testing.T) {
		res, err := env.client.UpdateBandTrack(ctx, &publicapi.UpdateBandTrackRequestBody{
			Title: publicapi.NewOptString("Track One (Remastered)"),
		}, publicapi.UpdateBandTrackParams{BandId: bandID, TrackId: track1})
		requireResponse[publicapi.UpdateBandTrackOK](t, res, err)

		getRes, getErr := env.client.GetBandTrack(ctx, publicapi.GetBandTrackParams{BandId: bandID, TrackId: track1})
		body := requireResponse[publicapi.BandTrack](t, getRes, getErr)
		require.Equal(t, "Track One (Remastered)", body.Title)
	})

	t.Run("create band setlist with initial tracks", func(t *testing.T) {
		res, err := env.client.CreateBandSetlist(ctx, &publicapi.CreateBandSetlistsRequestBody{
			Name:     "Show One",
			TrackIds: []uuid.UUID{track1, track2},
		}, publicapi.CreateBandSetlistParams{BandId: bandID})
		body := requireResponse[publicapi.CreateBandSetlistsResponseBody](t, res, err)
		require.NotEqual(t, uuid.Nil, body.ID)
		setlistID = body.ID
	})

	t.Run("get and list band setlists", func(t *testing.T) {
		getRes, getErr := env.client.GetBandSetlist(ctx, publicapi.GetBandSetlistParams{BandId: bandID, SetlistId: setlistID})
		getBody := requireResponse[publicapi.BandSetlist](t, getRes, getErr)
		require.Equal(t, "Show One", getBody.Name)
		require.Len(t, getBody.Tracks, 2)
		require.Equal(t, track1, getBody.Tracks[0].TrackId)
		require.Equal(t, track2, getBody.Tracks[1].TrackId)

		listRes, listErr := env.client.ListBandSetlists(ctx, publicapi.ListBandSetlistsParams{BandId: bandID})
		listBody := requireResponse[publicapi.ListBandSetlistsResponseBody](t, listRes, listErr)
		require.Len(t, listBody.Items, 1)
		require.Equal(t, setlistID, listBody.Items[0].ID)
	})

	t.Run("add a third track to the setlist", func(t *testing.T) {
		trackRes, trackErr := env.client.CreateBandTrack(ctx, &publicapi.CreateBandTrackRequestBody{
			Title:  "Track Three",
			Artist: "Artist Three",
		}, publicapi.CreateBandTrackParams{BandId: bandID})
		trackBody := requireResponse[publicapi.CreateBandTrackResponseBody](t, trackRes, trackErr)
		track3 = trackBody.ID

		addRes, addErr := env.client.AddSetlistTrack(ctx, &publicapi.AddSetlistTrackRequestBody{
			TrackId: track3,
		}, publicapi.AddSetlistTrackParams{BandId: bandID, SetlistId: setlistID})
		requireResponse[publicapi.AddSetlistTrackNoContent](t, addRes, addErr)

		getRes, getErr := env.client.GetBandSetlist(ctx, publicapi.GetBandSetlistParams{BandId: bandID, SetlistId: setlistID})
		body := requireResponse[publicapi.BandSetlist](t, getRes, getErr)
		require.Len(t, body.Tracks, 3)
		require.Equal(t, track3, body.Tracks[2].TrackId)
	})

	t.Run("remove a track from the setlist", func(t *testing.T) {
		removeRes, removeErr := env.client.RemoveSetlistTrack(ctx, publicapi.RemoveSetlistTrackParams{
			BandId:    bandID,
			SetlistId: setlistID,
			TrackId:   track2,
		})
		requireResponse[publicapi.RemoveSetlistTrackNoContent](t, removeRes, removeErr)

		getRes, getErr := env.client.GetBandSetlist(ctx, publicapi.GetBandSetlistParams{BandId: bandID, SetlistId: setlistID})
		body := requireResponse[publicapi.BandSetlist](t, getRes, getErr)
		require.Len(t, body.Tracks, 2)
		require.Equal(t, track1, body.Tracks[0].TrackId)
		require.Equal(t, track3, body.Tracks[1].TrackId)
	})

	t.Run("reorder setlist tracks", func(t *testing.T) {
		reorderRes, reorderErr := env.client.ReorderSetlistTracks(ctx, &publicapi.ReorderSetlistTracksRequestBody{
			TrackIds: []uuid.UUID{track3, track1},
		}, publicapi.ReorderSetlistTracksParams{BandId: bandID, SetlistId: setlistID})
		requireResponse[publicapi.ReorderSetlistTracksNoContent](t, reorderRes, reorderErr)

		getRes, getErr := env.client.GetBandSetlist(ctx, publicapi.GetBandSetlistParams{BandId: bandID, SetlistId: setlistID})
		body := requireResponse[publicapi.BandSetlist](t, getRes, getErr)
		require.Len(t, body.Tracks, 2)
		require.Equal(t, track3, body.Tracks[0].TrackId)
		require.Equal(t, 0, body.Tracks[0].Position)
		require.Equal(t, track1, body.Tracks[1].TrackId)
		require.Equal(t, 1, body.Tracks[1].Position)
	})

	t.Run("update the setlist", func(t *testing.T) {
		eventDate := time.Date(2026, time.September, 12, 0, 0, 0, 0, time.UTC)
		updateRes, updateErr := env.client.UpdateBandSetlist(ctx, &publicapi.UpdateBandSetlistRequestBody{
			Name:          publicapi.NewOptString("Show One (Final)"),
			EventLocation: publicapi.NewOptNilString("The Venue"),
			EventDate:     publicapi.NewOptNilDate(eventDate),
		}, publicapi.UpdateBandSetlistParams{BandId: bandID, SetlistId: setlistID})
		requireResponse[publicapi.UpdateBandSetlistOK](t, updateRes, updateErr)

		getRes, getErr := env.client.GetBandSetlist(ctx, publicapi.GetBandSetlistParams{BandId: bandID, SetlistId: setlistID})
		body := requireResponse[publicapi.BandSetlist](t, getRes, getErr)
		require.Equal(t, "Show One (Final)", body.Name)
		location, ok := body.EventLocation.Get()
		require.True(t, ok)
		require.Equal(t, "The Venue", location)
	})

	t.Run("bulk add tracks to setlist", func(t *testing.T) {
		track4Res, track4Err := env.client.CreateBandTrack(ctx, &publicapi.CreateBandTrackRequestBody{
			Title:  "Track Four",
			Artist: "Artist Four",
		}, publicapi.CreateBandTrackParams{BandId: bandID})
		track4Body := requireResponse[publicapi.CreateBandTrackResponseBody](t, track4Res, track4Err)
		track4 = track4Body.ID

		track5Res, track5Err := env.client.CreateBandTrack(ctx, &publicapi.CreateBandTrackRequestBody{
			Title:  "Track Five",
			Artist: "Artist Five",
		}, publicapi.CreateBandTrackParams{BandId: bandID})
		track5Body := requireResponse[publicapi.CreateBandTrackResponseBody](t, track5Res, track5Err)
		track5 = track5Body.ID

		addRes, addErr := env.client.AddSetlistTracks(ctx, &publicapi.AddSetlistTracksRequestBody{
			TrackIds: []uuid.UUID{track4, track5},
		}, publicapi.AddSetlistTracksParams{BandId: bandID, SetlistId: setlistID})
		requireResponse[publicapi.AddSetlistTracksNoContent](t, addRes, addErr)

		getRes, getErr := env.client.GetBandSetlist(ctx, publicapi.GetBandSetlistParams{BandId: bandID, SetlistId: setlistID})
		body := requireResponse[publicapi.BandSetlist](t, getRes, getErr)
		require.Len(t, body.Tracks, 4)
		require.Equal(t, track3, body.Tracks[0].TrackId)
		require.Equal(t, track1, body.Tracks[1].TrackId)
		require.Equal(t, track4, body.Tracks[2].TrackId)
		require.Equal(t, track5, body.Tracks[3].TrackId)
	})

	t.Run("list user tracks across bands", func(t *testing.T) {
		res, err := env.client.ListUserTracks(ctx, publicapi.ListUserTracksParams{})
		body := requireResponse[publicapi.ListUserTracksResponseBody](t, res, err)
		require.Len(t, body.Items, 5)
		for _, item := range body.Items {
			require.Equal(t, bandID, item.BandId)
			require.Equal(t, "The Wanderers Reborn", item.BandName)
		}

		filteredRes, filteredErr := env.client.ListUserTracks(ctx, publicapi.ListUserTracksParams{BandId: publicapi.NewOptUUID(bandID)})
		filteredBody := requireResponse[publicapi.ListUserTracksResponseBody](t, filteredRes, filteredErr)
		require.Len(t, filteredBody.Items, 5)

		otherBandRes, otherBandErr := env.client.ListUserTracks(ctx, publicapi.ListUserTracksParams{BandId: publicapi.NewOptUUID(uuid.New())})
		otherBandBody := requireResponse[publicapi.ListUserTracksResponseBody](t, otherBandRes, otherBandErr)
		require.Empty(t, otherBandBody.Items)
	})

	t.Run("list user setlists across bands", func(t *testing.T) {
		res, err := env.client.ListUserSetlists(ctx, publicapi.ListUserSetlistsParams{BandId: publicapi.NewOptUUID(bandID)})
		body := requireResponse[publicapi.ListUserSetlistsResponseBody](t, res, err)
		require.Len(t, body.Items, 1)
		require.Equal(t, setlistID, body.Items[0].ID)
		require.Equal(t, "Show One (Final)", body.Items[0].Name)
		require.Equal(t, 4, body.Items[0].TracksCount)
		require.Equal(t, bandID, body.Items[0].BandId)
		require.Equal(t, "The Wanderers Reborn", body.Items[0].BandName)
	})

	t.Run("member also sees the band's tracks and setlists", func(t *testing.T) {
		env.sec.token = memberToken
		defer func() { env.sec.token = ownerToken }()

		tracksRes, tracksErr := env.client.ListUserTracks(ctx, publicapi.ListUserTracksParams{})
		tracksBody := requireResponse[publicapi.ListUserTracksResponseBody](t, tracksRes, tracksErr)
		require.Len(t, tracksBody.Items, 5)

		setlistsRes, setlistsErr := env.client.ListUserSetlists(ctx, publicapi.ListUserSetlistsParams{})
		setlistsBody := requireResponse[publicapi.ListUserSetlistsResponseBody](t, setlistsRes, setlistsErr)
		require.Len(t, setlistsBody.Items, 1)
		require.Equal(t, setlistID, setlistsBody.Items[0].ID)
	})

	t.Run("owner removes the member from the band", func(t *testing.T) {
		res, err := env.client.RemoveBandMember(ctx, publicapi.RemoveBandMemberParams{
			BandId: bandID,
			UserId: memberID,
		})
		requireResponse[publicapi.RemoveBandMemberNoContent](t, res, err)
	})

	t.Run("clean up setlist, tracks and band", func(t *testing.T) {
		setlistRes, setlistErr := env.client.RemoveBandSetlist(ctx, publicapi.RemoveBandSetlistParams{
			BandId:    bandID,
			SetlistId: setlistID,
		})
		requireResponse[publicapi.RemoveBandSetlistNoContent](t, setlistRes, setlistErr)

		for _, trackID := range []uuid.UUID{track1, track2, track3, track4, track5} {
			trackRes, trackErr := env.client.RemoveBandTrack(ctx, publicapi.RemoveBandTrackParams{
				BandId:  bandID,
				TrackId: trackID,
			})
			requireResponse[publicapi.RemoveBandTrackNoContent](t, trackRes, trackErr)
		}

		bandRes, bandErr := env.client.RemoveBand(ctx, publicapi.RemoveBandParams{BandId: bandID})
		requireResponse[publicapi.RemoveBandNoContent](t, bandRes, bandErr)

		homepageRes, homepageErr := env.client.GetHomepageData(ctx)
		body := requireResponse[publicapi.HomepageData](t, homepageRes, homepageErr)
		require.Equal(t, 0, body.BandsCount)
	})
}

func requireResponse[T any](t *testing.T, res any, err error) *T {
	t.Helper()
	require.NoError(t, err)
	body, ok := res.(*T)
	require.Truef(t, ok, "unexpected response type %T", res)
	return body
}
