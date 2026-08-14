package domain

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

const (
	maxBandNameLength = 255

	inviteCodeAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	inviteCodeLength   = 6
)

var (
	ErrEmptyBandName      = errors.New("band name can not be empty")
	ErrBandNameTooLong    = fmt.Errorf("band name length should be less than or equal to %d", maxBandNameLength)
	ErrBandNotFound       = errors.New("band not found")
	ErrNotBandOwner       = errors.New("only band owner is allowed to perform this action")
	ErrBandMemberNotFound = errors.New("band member not found")
	ErrCannotRemoveOwner  = errors.New("band owner can not be removed from the band")
)

type Band struct {
	id         BandID
	name       string
	inviteCode string
	members    []BandMember
}

type BandRepository interface {
	NextID() BandID
	Store(*Band) error
	Get(BandID) (*Band, error)
	GetByInviteCode(inviteCode string) (*Band, error)
	Remove(BandID) error
}

func NewBand(
	id BandID,
	name string,
	ownerID UserID,
) (*Band, error) {
	err := validateBandNameLength(name)
	if err != nil {
		return nil, err
	}
	return &Band{
		id:         id,
		name:       name,
		inviteCode: generateInviteCode(),
		members:    []BandMember{{userID: ownerID, role: BandRoleOwner}},
	}, nil
}

func LoadBand(
	id BandID,
	name string,
	inviteCode string,
	members []BandMember,
) *Band {
	return &Band{
		id:         id,
		name:       name,
		inviteCode: inviteCode,
		members:    members,
	}
}

func (b *Band) ID() BandID {
	return b.id
}

func (b *Band) Name() string {
	return b.name
}

func (b *Band) SetName(name string) error {
	err := validateBandNameLength(name)
	if err != nil {
		return err
	}
	b.name = name
	return nil
}

func (b *Band) InviteCode() string {
	return b.inviteCode
}

func (b *Band) Members() []BandMember {
	return b.members
}

func (b *Band) AddMember(userID UserID, role BandRole) {
	for _, member := range b.members {
		if member.userID == userID {
			return
		}
	}
	b.members = append(b.members, BandMember{userID: userID, role: role})
}

func (b *Band) RemoveMember(userID UserID) {
	for i, member := range b.members {
		if member.userID == userID {
			b.members = append(b.members[:i], b.members[i+1:]...)
			return
		}
	}
}

func (b *Band) HasMember(userID UserID) bool {
	for _, member := range b.members {
		if member.userID == userID {
			return true
		}
	}
	return false
}

func (b *Band) IsOwner(userID UserID) bool {
	for _, member := range b.members {
		if member.userID == userID {
			return member.role == BandRoleOwner
		}
	}
	return false
}

func validateBandNameLength(name string) error {
	return checkStringLimits(name, maxBandNameLength, ErrEmptyBandName, ErrBandNameTooLong)
}

func generateInviteCode() string {
	code := make([]byte, inviteCodeLength)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(inviteCodeAlphabet))))
		if err != nil {
			panic(err)
		}
		code[i] = inviteCodeAlphabet[n.Int64()]
	}
	return string(code)
}
