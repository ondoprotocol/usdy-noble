package keeper_test

import (
	"testing"

	"cosmossdk.io/collections"
	"github.com/ondoprotocol/usdy-noble/v2/keeper"
	"github.com/ondoprotocol/usdy-noble/v2/types"
	"github.com/ondoprotocol/usdy-noble/v2/types/blocklist"
	"github.com/ondoprotocol/usdy-noble/v2/utils"
	"github.com/ondoprotocol/usdy-noble/v2/utils/mocks"
	"github.com/stretchr/testify/require"
)

func TestBlocklistTransferOwnership(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewBlocklistMsgServer(k)

	// ACT: Attempt to transfer ownership with no owner set.
	_, err := server.TransferOwnership(ctx, &blocklist.MsgTransferOwnership{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no blocklist owner")

	// ARRANGE: Set blocklist owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetBlocklistOwner(ctx, owner.Address))

	// ACT: Attempt to transfer ownership with invalid signer.
	_, err = server.TransferOwnership(ctx, &blocklist.MsgTransferOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, blocklist.ErrInvalidOwner.Error())

	// ACT: Attempt to transfer ownership to same owner.
	_, err = server.TransferOwnership(ctx, &blocklist.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: owner.Address,
	})
	// ASSERT: The action should've failed due to same owner.
	require.ErrorContains(t, err, blocklist.ErrSameOwner.Error())

	// ARRANGE: Generate a pending owner account.
	pendingOwner := utils.TestAccount()

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.BlocklistPendingOwner
	k.BlocklistPendingOwner = collections.NewItem(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		blocklist.PendingOwnerKey, "blocklist_pending_owner", collections.StringValue,
	)

	// ACT: Attempt to transfer ownership with failing BlocklistPendingOwner collection store.
	_, err = server.TransferOwnership(ctx, &blocklist.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: pendingOwner.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.BlocklistPendingOwner = tmp

	// ACT: Attempt to transfer ownership.
	_, err = server.TransferOwnership(ctx, &blocklist.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded, and set a pending owner in state.
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, k.GetBlocklistPendingOwner(ctx))
}

func TestBlocklistAcceptOwnership(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewBlocklistMsgServer(k)

	// ACT: Attempt to accept ownership with no pending owner set.
	_, err := server.AcceptOwnership(ctx, &blocklist.MsgAcceptOwnership{})
	// ASSERT: The action should've failed due to no pending owner set.
	require.ErrorContains(t, err, "there is no pending blocklist owner")

	// ARRANGE: Set blocklist pending owner in state.
	pendingOwner := utils.TestAccount()
	require.NoError(t, k.SetBlocklistPendingOwner(ctx, pendingOwner.Address))

	// ACT: Attempt to accept ownership with invalid signer.
	_, err = server.AcceptOwnership(ctx, &blocklist.MsgAcceptOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, blocklist.ErrInvalidPendingOwner.Error())

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.BlocklistOwner
	k.BlocklistOwner = collections.NewItem(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		blocklist.OwnerKey, "blocklist_owner", collections.StringValue,
	)

	// ACT: Attempt to accept ownership with failing BlocklistOwner collection store.
	_, err = server.AcceptOwnership(ctx, &blocklist.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.BlocklistOwner = tmp

	// ARRANGE: Set up a failing collection store for the attribute delete.
	tmp = k.BlocklistPendingOwner
	k.BlocklistPendingOwner = collections.NewItem(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Delete, utils.GetKVStore(ctx, types.ModuleName))),
		blocklist.PendingOwnerKey, "blocklist_pending_owner", collections.StringValue,
	)

	// ACT: Attempt to accept ownership with failing BlocklistPendingOwner collection store.
	_, err = server.AcceptOwnership(ctx, &blocklist.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've failed due to collection store delete error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.BlocklistPendingOwner = tmp

	// ACT: Attempt to accept ownership.
	_, err = server.AcceptOwnership(ctx, &blocklist.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded, and updated the owner in state.
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, k.GetBlocklistOwner(ctx))
	require.Empty(t, k.GetBlocklistPendingOwner(ctx))
}

func TestAddToBlocklist(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewBlocklistMsgServer(k)

	// ACT: Attempt to add to blocklist with no owner set.
	_, err := server.AddToBlocklist(ctx, &blocklist.MsgAddToBlocklist{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no blocklist owner")

	// ARRANGE: Set blocklist owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetBlocklistOwner(ctx, owner.Address))

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

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.BlockedAddresses
	k.BlockedAddresses = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		blocklist.BlockedAddressPrefix, "blocked_addresses", collections.BytesKey, collections.BytesValue,
	)

	// ACT: Attempt to add to blocklist with failing BlockedAddresses collection store.
	_, err = server.AddToBlocklist(ctx, &blocklist.MsgAddToBlocklist{
		Signer:   owner.Address,
		Accounts: []string{user.Address},
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.BlockedAddresses = tmp

	// ACT: Attempt to add to blocklist.
	_, err = server.AddToBlocklist(ctx, &blocklist.MsgAddToBlocklist{
		Signer:   owner.Address,
		Accounts: []string{user.Address},
	})
	// ASSERT: The action should've succeeded, and blocked the user in state.
	require.NoError(t, err)
	require.True(t, k.HasBlockedAddress(ctx, user.Bytes))
}

func TestRemoveFromBlocklist(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewBlocklistMsgServer(k)

	// ACT: Attempt to remove from blocklist with no owner set.
	_, err := server.RemoveFromBlocklist(ctx, &blocklist.MsgRemoveFromBlocklist{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no blocklist owner")

	// ARRANGE: Set blocklist owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetBlocklistOwner(ctx, owner.Address))

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
	require.False(t, k.HasBlockedAddress(ctx, user.Bytes))

	// ARRANGE: Set user as blocked in state.
	require.NoError(t, k.SetBlockedAddress(ctx, user.Bytes))
	require.True(t, k.HasBlockedAddress(ctx, user.Bytes))

	// ARRANGE: Set up a failing collection store for the attribute delete.
	tmp := k.BlockedAddresses
	k.BlockedAddresses = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Delete, utils.GetKVStore(ctx, types.ModuleName))),
		blocklist.BlockedAddressPrefix, "blocked_addresses", collections.BytesKey, collections.BytesValue,
	)

	// ACT: Attempt to remove from blocklist with failing BlockedAddresses collection store.
	_, err = server.RemoveFromBlocklist(ctx, &blocklist.MsgRemoveFromBlocklist{
		Signer:   owner.Address,
		Accounts: []string{user.Address},
	})
	// ASSERT: The action should've failed due to collection store delete error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.BlockedAddresses = tmp

	// ACT: Attempt to remove from blocklist.
	_, err = server.RemoveFromBlocklist(ctx, &blocklist.MsgRemoveFromBlocklist{
		Signer:   owner.Address,
		Accounts: []string{user.Address},
	})
	// ASSERT: The action should've succeeded, and unblocked the user.
	require.NoError(t, err)
	require.False(t, k.HasBlockedAddress(ctx, user.Bytes))
}
