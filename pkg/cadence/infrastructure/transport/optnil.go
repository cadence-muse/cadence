package transport

import (
	"cadence/pkg/common/maybe"
)

type OptValue[T any] interface {
	IsSet() bool
	Get() (v T, ok bool)
}

type OptNilValue[T any] interface {
	OptValue[T]
	IsNull() bool
}

func maybeValueFromOpt[T any](opt OptValue[T]) maybe.Maybe[T] {
	if !opt.IsSet() {
		return maybe.NewAbsent[T]()
	}
	value, _ := opt.Get()
	return maybe.NewJust(value)
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
