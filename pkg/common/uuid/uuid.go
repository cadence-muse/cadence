package uuid

import "github.com/google/uuid"

const uuidSize = 16

// UUID Version 7 time-ordered UUID
type UUID [uuidSize]byte

func Generate() UUID {
	impl, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return UUID(impl)
}
