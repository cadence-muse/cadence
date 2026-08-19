package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cadence/pkg/cadence/domain"
	"cadence/pkg/common/maybe"
	"cadence/pkg/common/uuid"
)

func TestBandService_Create(t *testing.T) {
	t.Run("creates band with owner as sole member", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)
		ownerID := uuid.Generate()

		bandID, err := svc.Create(context.Background(), CreateBandParams{OwnerID: ownerID, Name: "The Testers"})
		require.NoError(t, err)

		band, err := executor.repoProvider().BandRepository().Get(bandID)
		require.NoError(t, err)
		assert.Equal(t, "The Testers", band.Name())
		assert.True(t, band.IsOwner(ownerID))
	})
}

func TestBandService_Update(t *testing.T) {
	t.Run("renames an existing band", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)
		bandID := seedBand(t, executor, uuid.Generate())

		err := svc.Update(context.Background(), UpdateBandParams{BandID: bandID, Name: maybe.NewJust("New Name")})
		require.NoError(t, err)

		band, err := executor.repoProvider().BandRepository().Get(bandID)
		require.NoError(t, err)
		assert.Equal(t, "New Name", band.Name())
	})

	t.Run("band not found", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)

		err := svc.Update(context.Background(), UpdateBandParams{BandID: uuid.Generate(), Name: maybe.NewJust("New Name")})
		require.ErrorIs(t, err, domain.ErrBandNotFound)
	})
}

func TestBandService_JoinByInviteCode(t *testing.T) {
	t.Run("adds the user as a member", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)
		ownerID := uuid.Generate()
		bandID := seedBand(t, executor, ownerID)
		band, err := executor.repoProvider().BandRepository().Get(bandID)
		require.NoError(t, err)
		inviteCode := band.InviteCode()

		joinerID := uuid.Generate()
		err = svc.JoinByInviteCode(context.Background(), joinerID, inviteCode)
		require.NoError(t, err)

		band, err = executor.repoProvider().BandRepository().Get(bandID)
		require.NoError(t, err)
		assert.True(t, band.HasMember(joinerID))
		assert.False(t, band.IsOwner(joinerID))
	})

	t.Run("invalid invite code is rejected", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)

		err := svc.JoinByInviteCode(context.Background(), uuid.Generate(), "NOPE00")
		require.ErrorIs(t, err, domain.ErrBandNotFound)
	})
}

func TestBandService_Remove(t *testing.T) {
	t.Run("owner can remove the band", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)
		ownerID := uuid.Generate()
		bandID := seedBand(t, executor, ownerID)

		err := svc.Remove(context.Background(), bandID, ownerID)
		require.NoError(t, err)

		_, err = executor.repoProvider().BandRepository().Get(bandID)
		require.ErrorIs(t, err, domain.ErrBandNotFound)
	})

	t.Run("non-owner can not remove the band", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)
		ownerID := uuid.Generate()
		bandID := seedBand(t, executor, ownerID)

		err := svc.Remove(context.Background(), bandID, uuid.Generate())
		require.ErrorIs(t, err, domain.ErrNotBandOwner)
	})

	t.Run("band not found", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)

		err := svc.Remove(context.Background(), uuid.Generate(), uuid.Generate())
		require.ErrorIs(t, err, domain.ErrBandNotFound)
	})
}

func TestBandService_RemoveMember(t *testing.T) {
	t.Run("member can remove themselves", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)
		ownerID := uuid.Generate()
		bandID := seedBand(t, executor, ownerID)
		memberID := addBandMember(t, executor, bandID)

		err := svc.RemoveMember(context.Background(), bandID, memberID, memberID)
		require.NoError(t, err)

		band, err := executor.repoProvider().BandRepository().Get(bandID)
		require.NoError(t, err)
		assert.False(t, band.HasMember(memberID))
	})

	t.Run("owner can remove another member", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)
		ownerID := uuid.Generate()
		bandID := seedBand(t, executor, ownerID)
		memberID := addBandMember(t, executor, bandID)

		err := svc.RemoveMember(context.Background(), bandID, memberID, ownerID)
		require.NoError(t, err)

		band, err := executor.repoProvider().BandRepository().Get(bandID)
		require.NoError(t, err)
		assert.False(t, band.HasMember(memberID))
	})

	t.Run("non-owner can not remove another member", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)
		ownerID := uuid.Generate()
		bandID := seedBand(t, executor, ownerID)
		memberA := addBandMember(t, executor, bandID)
		memberB := addBandMember(t, executor, bandID)

		err := svc.RemoveMember(context.Background(), bandID, memberB, memberA)
		require.ErrorIs(t, err, domain.ErrNotBandOwner)
	})

	t.Run("removing a non-member is rejected", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)
		ownerID := uuid.Generate()
		bandID := seedBand(t, executor, ownerID)

		err := svc.RemoveMember(context.Background(), bandID, uuid.Generate(), ownerID)
		require.ErrorIs(t, err, domain.ErrBandMemberNotFound)
	})

	t.Run("removing the owner is rejected", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)
		ownerID := uuid.Generate()
		bandID := seedBand(t, executor, ownerID)

		err := svc.RemoveMember(context.Background(), bandID, ownerID, ownerID)
		require.ErrorIs(t, err, domain.ErrCannotRemoveOwner)
	})
}

func TestBandService_TransferOwnership(t *testing.T) {
	t.Run("owner transfers ownership to another member", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)
		ownerID := uuid.Generate()
		bandID := seedBand(t, executor, ownerID)
		memberID := addBandMember(t, executor, bandID)

		err := svc.TransferOwnership(context.Background(), bandID, ownerID, memberID)
		require.NoError(t, err)

		band, err := executor.repoProvider().BandRepository().Get(bandID)
		require.NoError(t, err)
		assert.True(t, band.IsOwner(memberID))
		assert.False(t, band.IsOwner(ownerID))
	})

	t.Run("non-owner can not transfer ownership", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)
		ownerID := uuid.Generate()
		bandID := seedBand(t, executor, ownerID)
		memberID := addBandMember(t, executor, bandID)

		err := svc.TransferOwnership(context.Background(), bandID, memberID, memberID)
		require.ErrorIs(t, err, domain.ErrNotBandOwner)
	})

	t.Run("transferring to a non-member is rejected", func(t *testing.T) {
		executor := newFakeExecutor()
		svc := NewBandService(executor)
		ownerID := uuid.Generate()
		bandID := seedBand(t, executor, ownerID)

		err := svc.TransferOwnership(context.Background(), bandID, ownerID, uuid.Generate())
		require.ErrorIs(t, err, domain.ErrBandMemberNotFound)
	})
}

func seedBand(t *testing.T, executor *fakeExecutor, ownerID uuid.UUID) uuid.UUID {
	t.Helper()

	repo := executor.repoProvider().BandRepository()
	band, err := domain.NewBand(repo.NextID(), "Band", ownerID)
	require.NoError(t, err)
	require.NoError(t, repo.Store(band))
	return band.ID()
}

func addBandMember(t *testing.T, executor *fakeExecutor, bandID uuid.UUID) uuid.UUID {
	t.Helper()

	repo := executor.repoProvider().BandRepository()
	band, err := repo.Get(bandID)
	require.NoError(t, err)

	memberID := uuid.Generate()
	band.AddMember(memberID, domain.BandRoleMember)
	require.NoError(t, repo.Store(band))
	return memberID
}
