package keeper_test

import (
	"testing"

	"github.com/noble-assets/aura/utils"
	"github.com/noble-assets/aura/utils/mocks"
	"github.com/noble-assets/aura/x/aura/keeper"
	"github.com/noble-assets/aura/x/aura/types/blocklist"
	"github.com/stretchr/testify/require"
)

func TestBlocklistTransferOwnership(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewBlocklistMsgServer(k)

	// ACT: Attempt to transfer ownership with no owner set.
	_, err := server.TransferOwnership(ctx, &blocklist.MsgTransferOwnership{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "unable to retrieve blocklist owner from state")

	// ARRANGE: Set blocklist owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.BlocklistOwner.Set(ctx, owner.Address))

	// ACT: Attempt to transfer ownership with invalid signer.
	_, err = server.TransferOwnership(ctx, &blocklist.MsgTransferOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, blocklist.ErrInvalidOwner.Error())

	// ARRANGE: Generate a pending owner account.
	pendingOwner := utils.TestAccount()

	// ACT: Attempt to transfer ownership.
	_, err = server.TransferOwnership(ctx, &blocklist.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded, and set a pending owner in state.
	require.NoError(t, err)
	res, err := k.BlocklistPendingOwner.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, res)
}

func TestBlocklistAcceptOwnership(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewBlocklistMsgServer(k)

	// ACT: Attempt to accept ownership with no pending owner set.
	_, err := server.AcceptOwnership(ctx, &blocklist.MsgAcceptOwnership{})
	// ASSERT: The action should've failed due to no pending owner set.
	require.ErrorContains(t, err, "there is no blocklist pending owner")

	// ARRANGE: Set blocklist pending owner in state.
	pendingOwner := utils.TestAccount()
	require.NoError(t, k.BlocklistPendingOwner.Set(ctx, pendingOwner.Address))

	// ACT: Attempt to accept ownership with invalid signer.
	_, err = server.AcceptOwnership(ctx, &blocklist.MsgAcceptOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, blocklist.ErrInvalidPendingOwner.Error())

	// ACT: Attempt to accept ownership.
	_, err = server.AcceptOwnership(ctx, &blocklist.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded, and updated the owner in state.
	require.NoError(t, err)
	res, err := k.BlocklistOwner.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, res)
	has, err := k.BlocklistPendingOwner.Has(ctx)
	require.NoError(t, err)
	require.False(t, has)
}

func TestAddToBlocklist(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewBlocklistMsgServer(k)

	// ACT: Attempt to add to blocklist with no owner set.
	_, err := server.AddToBlocklist(ctx, &blocklist.MsgAddToBlocklist{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "unable to retrieve blocklist owner from state")

	// ARRANGE: Set blocklist owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.BlocklistOwner.Set(ctx, owner.Address))

	// ACT: Attempt to add to blocklist with invalid signer.
	_, err = server.AddToBlocklist(ctx, &blocklist.MsgAddToBlocklist{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, blocklist.ErrInvalidOwner.Error())

	// ACT: Attempt to add to blocklist with invalid account address.
	_, err = server.AddToBlocklist(ctx, &blocklist.MsgAddToBlocklist{
		Signer:   owner.Address,
		Accounts: []string{"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"},
	})
	// ASSERT: The action should've failed due to invalid account address.
	require.ErrorContains(t, err, "unable to decode account address")

	// ARRANGE: Generate a user account.
	user := utils.TestAccount()

	// ACT: Attempt to add to blocklist.
	_, err = server.AddToBlocklist(ctx, &blocklist.MsgAddToBlocklist{
		Signer:   owner.Address,
		Accounts: []string{user.Address},
	})
	// ASSERT: The action should've succeeded, and blocked the user in state.
	require.NoError(t, err)
	res, err := k.BlockedAddresses.Get(ctx, user.Bytes)
	require.NoError(t, err)
	require.True(t, res)
}

func TestRemoveFromBlocklist(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewBlocklistMsgServer(k)

	// ACT: Attempt to remove from blocklist with no owner set.
	_, err := server.RemoveFromBlocklist(ctx, &blocklist.MsgRemoveFromBlocklist{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "unable to retrieve blocklist owner from state")

	// ARRANGE: Set blocklist owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.BlocklistOwner.Set(ctx, owner.Address))

	// ACT: Attempt to remove from blocklist with invalid signer.
	_, err = server.RemoveFromBlocklist(ctx, &blocklist.MsgRemoveFromBlocklist{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, blocklist.ErrInvalidOwner.Error())

	// ACT: Attempt to remove from blocklist with invalid account address.
	_, err = server.RemoveFromBlocklist(ctx, &blocklist.MsgRemoveFromBlocklist{
		Signer:   owner.Address,
		Accounts: []string{"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"},
	})
	// ASSERT: The action should've failed due to invalid account address.
	require.ErrorContains(t, err, "unable to decode account address")

	// ARRANGE: Generate a user account.
	user := utils.TestAccount()

	// ACT: Attempt to remove from blocklist with unblocked account.
	_, err = server.RemoveFromBlocklist(ctx, &blocklist.MsgRemoveFromBlocklist{
		Signer:   owner.Address,
		Accounts: []string{user.Address},
	})
	// ASSERT: The action should've succeeded, and the user shouldn't be blocked.
	require.NoError(t, err)
	has, err := k.BlockedAddresses.Has(ctx, user.Bytes)
	require.NoError(t, err)
	require.False(t, has)

	// ARRANGE: Set user as blocked in state.
	require.NoError(t, k.BlockedAddresses.Set(ctx, user.Bytes, true))

	// ACT: Attempt to remove from blocklist.
	_, err = server.RemoveFromBlocklist(ctx, &blocklist.MsgRemoveFromBlocklist{
		Signer:   owner.Address,
		Accounts: []string{user.Address},
	})
	// ASSERT: The action should've succeeded, and unblocked the user.
	require.NoError(t, err)
	has, err = k.BlockedAddresses.Has(ctx, user.Bytes)
	require.NoError(t, err)
	require.False(t, has)
}
