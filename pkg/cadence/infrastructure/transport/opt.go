package transport

import (
	"time"

	"github.com/nightnoryu/go-kita/maybe"

	"cadence/api/server/publicapi"
	"cadence/pkg/common/uuid"
	"cadence/pkg/common/valuetypes"
)

type OptValue[T any] interface {
	IsSet() bool
	Get() (v T, ok bool)
}

func maybeValueFromOpt[T any](opt OptValue[T]) maybe.Maybe[T] {
	if !opt.IsSet() {
		return maybe.NewAbsent[T]()
	}
	value, _ := opt.Get()
	return maybe.NewJust(value)
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
	return publicapi.NewOptInt(durationToIntSeconds(value))
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
	return maybe.NewJust(intSecondsToDuration(opt.Value))
}

func maybeKeyFromOpt(opt publicapi.OptString) (maybe.Maybe[valuetypes.MusicalKey], error) {
	if !opt.Set {
		return maybe.NewAbsent[valuetypes.MusicalKey](), nil
	}
	return parseMusicalKey(opt.Value)
}

func maybeTempoFromOpt(opt publicapi.OptInt) (maybe.Maybe[valuetypes.Tempo], error) {
	if !opt.Set {
		return maybe.NewAbsent[valuetypes.Tempo](), nil
	}
	return parseTempo(opt.Value)
}

func maybeUUIDFromOpt(opt publicapi.OptUUID) maybe.Maybe[uuid.UUID] {
	value, ok := opt.Get()
	if !ok {
		return maybe.NewAbsent[uuid.UUID]()
	}
	return maybe.NewJust(uuid.UUID(value))
}
