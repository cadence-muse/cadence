package domain

import (
	"errors"
	"fmt"
)

const (
	maxBandNameLength = 255
)

var (
	ErrEmptyBandName   = errors.New("band name can not be empty")
	ErrBandNameTooLong = fmt.Errorf("band name length should be less than or equal to %d", maxTrackTitleLength)
	ErrBandNotFound    = errors.New("band not found")
)

type Band struct {
	id   BandID
	name string
}

type BandRepository interface {
	NextID() BandID
	Store(*Band) error
	Get(BandID) (*Band, error)
}

func NewBand(
	id BandID,
	name string,
) (*Band, error) {
	err := validateBandNameLength(name)
	if err != nil {
		return nil, err
	}
	return &Band{
		id:   id,
		name: name,
	}, nil
}

func LoadBand(
	id BandID,
	name string,
) *Band {
	return &Band{
		id:   id,
		name: name,
	}
}

func (b *Band) ID() BandID {
	return b.id
}

func (b *Band) Name() string {
	return b.name
}

func validateBandNameLength(name string) error {
	return checkStringLimits(name, maxBandNameLength, ErrEmptyBandName, ErrBandNameTooLong)
}
