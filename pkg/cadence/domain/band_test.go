package domain

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cadence/pkg/common/uuid"
)

var inviteCodePattern = regexp.MustCompile(`^[A-Z0-9]{6}$`)

func TestNewBand(t *testing.T) {
	t.Run("valid name creates band with owner as sole member", func(t *testing.T) {
		id := uuid.Generate()
		ownerID := uuid.Generate()

		band, err := NewBand(id, "The Testers", ownerID)
		require.NoError(t, err)
		assert.Equal(t, id, band.ID())
		assert.Equal(t, "The Testers", band.Name())
		assert.Regexp(t, inviteCodePattern, band.InviteCode())

		members := band.Members()
		require.Len(t, members, 1)
		assert.Equal(t, ownerID, members[0].UserID())
		assert.Equal(t, BandRoleOwner, members[0].Role())
	})

	t.Run("empty name is rejected", func(t *testing.T) {
		_, err := NewBand(uuid.Generate(), "", uuid.Generate())
		assert.ErrorIs(t, err, ErrEmptyBandName)
	})

	t.Run("name over the length limit is rejected", func(t *testing.T) {
		_, err := NewBand(uuid.Generate(), strings.Repeat("a", maxBandNameLength+1), uuid.Generate())
		assert.ErrorIs(t, err, ErrBandNameTooLong)
	})

	t.Run("name at the length limit is accepted", func(t *testing.T) {
		_, err := NewBand(uuid.Generate(), strings.Repeat("a", maxBandNameLength), uuid.Generate())
		assert.NoError(t, err)
	})

	t.Run("generates distinct invite codes", func(t *testing.T) {
		bandA, err := NewBand(uuid.Generate(), "Band A", uuid.Generate())
		require.NoError(t, err)
		bandB, err := NewBand(uuid.Generate(), "Band B", uuid.Generate())
		require.NoError(t, err)
		assert.NotEqual(t, bandA.InviteCode(), bandB.InviteCode())
	})
}

func TestLoadBand(t *testing.T) {
	id := uuid.Generate()
	ownerID := uuid.Generate()
	members := []BandMember{LoadBandMember(ownerID, BandRoleOwner)}

	band := LoadBand(id, "Loaded Band", "AB12CD", members)

	assert.Equal(t, id, band.ID())
	assert.Equal(t, "Loaded Band", band.Name())
	assert.Equal(t, "AB12CD", band.InviteCode())
	require.Len(t, band.Members(), 1)
	assert.Equal(t, ownerID, band.Members()[0].UserID())
}

func TestBand_AddMember(t *testing.T) {
	t.Run("adds a new member", func(t *testing.T) {
		ownerID := uuid.Generate()
		band, err := NewBand(uuid.Generate(), "Band", ownerID)
		require.NoError(t, err)

		newMemberID := uuid.Generate()
		band.AddMember(newMemberID, BandRoleMember)

		members := band.Members()
		require.Len(t, members, 2)
		assert.Equal(t, newMemberID, members[1].UserID())
		assert.Equal(t, BandRoleMember, members[1].Role())
	})

	t.Run("adding an existing member is a no-op", func(t *testing.T) {
		ownerID := uuid.Generate()
		band, err := NewBand(uuid.Generate(), "Band", ownerID)
		require.NoError(t, err)

		band.AddMember(ownerID, BandRoleMember)

		members := band.Members()
		require.Len(t, members, 1)
		assert.Equal(t, BandRoleOwner, members[0].Role())
	})
}

func TestBand_RemoveMember(t *testing.T) {
	t.Run("removes an existing member", func(t *testing.T) {
		ownerID := uuid.Generate()
		memberID := uuid.Generate()
		band, err := NewBand(uuid.Generate(), "Band", ownerID)
		require.NoError(t, err)
		band.AddMember(memberID, BandRoleMember)

		band.RemoveMember(memberID)

		members := band.Members()
		require.Len(t, members, 1)
		assert.Equal(t, ownerID, members[0].UserID())
	})

	t.Run("removing a non-member is a no-op", func(t *testing.T) {
		ownerID := uuid.Generate()
		band, err := NewBand(uuid.Generate(), "Band", ownerID)
		require.NoError(t, err)

		band.RemoveMember(uuid.Generate())

		assert.Len(t, band.Members(), 1)
	})
}
