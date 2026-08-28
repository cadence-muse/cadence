package transport

import (
	"time"

	"github.com/nightnoryu/go-kita/maybe"

	"cadence/api/server/publicapi"
	"cadence/pkg/common/valuetypes"
)

type OptNilValue[T any] interface {
	OptValue[T]
	IsNull() bool
}

func maybeValueFromOptNil[T any](opt OptNilValue[T]) maybe.Maybe[T] {
	if !opt.IsSet() {
		return maybe.NewAbsent[T]()
	}
	if opt.IsNull() {
		return maybe.NewNone[T]()
	}
	value, _ := opt.Get()
	return maybe.NewJust(value)
}

func maybeDurationFromOptNil(opt publicapi.OptNilInt) maybe.Maybe[time.Duration] {
	if !opt.Set {
		return maybe.NewAbsent[time.Duration]()
	}
	if opt.Null {
		return maybe.NewNone[time.Duration]()
	}
	return maybe.NewJust(intSecondsToDuration(opt.Value))
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

func maybeTempoFromOptNil(opt publicapi.OptNilInt) (maybe.Maybe[valuetypes.Tempo], error) {
	if !opt.Set {
		return maybe.NewAbsent[valuetypes.Tempo](), nil
	}
	if opt.Null {
		return maybe.NewNone[valuetypes.Tempo](), nil
	}
	return parseTempo(opt.Value)
}
