package transport

import (
	googleuuid "github.com/google/uuid"

	"cadence/api/server/publicapi"
	"cadence/pkg/cadence/app/query"
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
