package maybe

// Maybe is a monad that provides clear and explicit semantic of null value
type Maybe[T any] struct {
	v    T
	just bool
}

func NewJust[T any](v T) Maybe[T] {
	return Maybe[T]{
		v:    v,
		just: true,
	}
}

func NewNone[T any]() Maybe[T] {
	return Maybe[T]{}
}

func Valid[T any](maybe Maybe[T]) bool {
	return maybe.just
}

func Just[T any](maybe Maybe[T]) T {
	if !Valid(maybe) {
		panic("violated usage of maybe: Just on non Valid Maybe")
	}
	return maybe.v
}

func JustValid[T any](maybe Maybe[T]) (v T, ok bool) {
	if !Valid(maybe) {
		ok = false
		return
	}

	return Just(maybe), true
}
