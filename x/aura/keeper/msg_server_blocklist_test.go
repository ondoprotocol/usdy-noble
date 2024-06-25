package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/aura/utils"
	"github.com/ondoprotocol/aura/utils/mocks"
	"github.com/ondoprotocol/aura/x/aura/keeper"
	"github.com/ondoprotocol/aura/x/aura/types/blocklist"
	"github.com/stretchr/testify/require"
)

func TestBlocklistTransferOwnership(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewBlocklistMsgServer(k)

	// ACT: Attempt to transfer ownership with no owner set.
	_, err := server.TransferOwnership(goCtx, &blocklist.MsgTransferOwnership{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no blocklist owner")

	// ARRANGE: Set blocklist owner in state.
	owner := utils.TestAccount()
	k.SetBlocklistOwner(ctx, owner.Address)

	// ACT: Attempt to transfer ownership with invalid signer.
	_, err = server.TransferOwnership(goCtx, &blocklist.MsgTransferOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, blocklist.ErrInvalidOwner.Error())

	// ACT: Attempt to transfer ownership to same owner.
	_, err = server.TransferOwnership(goCtx, &blocklist.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: owner.Address,
	})
	// ASSERT: The action should've failed due to same owner.
	require.ErrorContains(t, err, blocklist.ErrSameOwner.Error())

	// ARRANGE: Generate a pending owner account.
	pendingOwner := utils.TestAccount()

	// ACT: Attempt to transfer ownership.
	_, err = server.TransferOwnership(goCtx, &blocklist.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded, and set a pending owner in state.
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, k.GetBlocklistPendingOwner(ctx))
}

func TestBlocklistAcceptOwnership(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewBlocklistMsgServer(k)

	// ACT: Attempt to accept ownership with no pending owner set.
	_, err := server.AcceptOwnership(goCtx, &blocklist.MsgAcceptOwnership{})
	// ASSERT: The action should've failed due to no pending owner set.
	require.ErrorContains(t, err, "there is no pending blocklist owner")

	// ARRANGE: Set blocklist pending owner in state.
	pendingOwner := utils.TestAccount()
	k.SetBlocklistPendingOwner(ctx, pendingOwner.Address)

	// ACT: Attempt to accept ownership with invalid signer.
	_, err = server.AcceptOwnership(goCtx, &blocklist.MsgAcceptOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, blocklist.ErrInvalidPendingOwner.Error())

	// ACT: Attempt to accept ownership.
	_, err = server.AcceptOwnership(goCtx, &blocklist.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded, and updated the owner in state.
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, k.GetBlocklistOwner(ctx))
	require.Empty(t, k.GetBlocklistPendingOwner(ctx))
}

func TestAddToBlocklist(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewBlocklistMsgServer(k)

	// ACT: Attempt to add to blocklist with no owner set.
	_, err := server.AddToBlocklist(goCtx, &blocklist.MsgAddToBlocklist{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no blocklist owner")

	// ARRANGE: Set blocklist owner in state.
	owner := utils.TestAccount()
	k.SetBlocklistOwner(ctx, owner.Address)

	// ACT: Attempt to add to blocklist with invalid signer.
	_, err = server.AddToBlocklist(goCtx, &blocklist.MsgAddToBlocklist{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, blocklist.ErrInvalidOwner.Error())

	// ACT: Attempt to add to blocklist with invalid account address.
	_, err = server.AddToBlocklist(goCtx, &blocklist.MsgAddToBlocklist{
		Signer:   owner.Address,
		Accounts: []string{"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"},
	})
	// ASSERT: The action should've failed due to invalid account address.
	require.ErrorContains(t, err, "unable to decode account address")

	// ARRANGE: Generate a user account.
	user := utils.TestAccount()

	// ACT: Attempt to add to blocklist.
	_, err = server.AddToBlocklist(goCtx, &blocklist.MsgAddToBlocklist{
		Signer:   owner.Address,
		Accounts: []string{user.Address},
	})
	// ASSERT: The action should've succeeded, and blocked the user in state.
	require.NoError(t, err)
	require.True(t, k.HasBlockedAddress(ctx, user.Bytes))
}

func TestRemoveFromBlocklist(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewBlocklistMsgServer(k)

	// ACT: Attempt to remove from blocklist with no owner set.
	_, err := server.RemoveFromBlocklist(goCtx, &blocklist.MsgRemoveFromBlocklist{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no blocklist owner")

	// ARRANGE: Set blocklist owner in state.
	owner := utils.TestAccount()
	k.SetBlocklistOwner(ctx, owner.Address)

	// ACT: Attempt to remove from blocklist with invalid signer.
	_, err = server.RemoveFromBlocklist(goCtx, &blocklist.MsgRemoveFromBlocklist{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, blocklist.ErrInvalidOwner.Error())

	// ACT: Attempt to remove from blocklist with invalid account address.
	_, err = server.RemoveFromBlocklist(goCtx, &blocklist.MsgRemoveFromBlocklist{
		Signer:   owner.Address,
		Accounts: []string{"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"},
	})
	// ASSERT: The action should've failed due to invalid account address.
	require.ErrorContains(t, err, "unable to decode account address")

	// ARRANGE: Generate a user account.
	user := utils.TestAccount()

	// ACT: Attempt to remove from blocklist with unblocked account.
	_, err = server.RemoveFromBlocklist(goCtx, &blocklist.MsgRemoveFromBlocklist{
		Signer:   owner.Address,
		Accounts: []string{user.Address},
	})
	// ASSERT: The action should've succeeded, and the user shouldn't be blocked.
	require.NoError(t, err)
	require.False(t, k.HasBlockedAddress(ctx, user.Bytes))

	// ARRANGE: Set user as blocked in state.
	k.SetBlockedAddress(ctx, user.Bytes)
	require.True(t, k.HasBlockedAddress(ctx, user.Bytes))

	// ACT: Attempt to remove from blocklist.
	_, err = server.RemoveFromBlocklist(goCtx, &blocklist.MsgRemoveFromBlocklist{
		Signer:   owner.Address,
		Accounts: []string{user.Address},
	})
	// ASSERT: The action should've succeeded, and unblocked the user.
	require.NoError(t, err)
	require.False(t, k.HasBlockedAddress(ctx, user.Bytes))
}
