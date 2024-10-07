package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ondoprotocol/usdy-noble/v2/utils"
	"github.com/ondoprotocol/usdy-noble/v2/utils/mocks"
	"github.com/ondoprotocol/usdy-noble/v2/x/aura/keeper"
	"github.com/ondoprotocol/usdy-noble/v2/x/aura/types/blocklist"
	"github.com/stretchr/testify/require"
)

func TestBlocklistOwnerQuery(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewBlocklistQueryServer(k)

	// ACT: Attempt to query blocklist owner with invalid request.
	_, err := server.Owner(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ARRANGE: Set blocklist owner in state.
	owner := utils.TestAccount()
	k.SetBlocklistOwner(ctx, owner.Address)

	// ACT: Attempt to query blocklist owner with state.
	res, err := server.Owner(ctx, &blocklist.QueryOwner{})
	// ASSERT: The query should've succeeded, with empty pending owner.
	require.NoError(t, err)
	require.Equal(t, owner.Address, res.Owner)
	require.Empty(t, res.PendingOwner)

	// ARRANGE: Set blocklist pending owner in state.
	pendingOwner := utils.TestAccount()
	k.SetBlocklistPendingOwner(ctx, pendingOwner.Address)

	// ACT: Attempt to query blocklist owner with state.
	res, err = server.Owner(ctx, &blocklist.QueryOwner{})
	// ASSERT: The query should've succeeded, with pending owner.
	require.NoError(t, err)
	require.Equal(t, owner.Address, res.Owner)
	require.Equal(t, pendingOwner.Address, res.PendingOwner)
}

func TestBlocklistAddressesQuery(t *testing.T) {
	// NOTE: Query pagination is assumed working, so isn't testing here.

	k, ctx := mocks.AuraKeeper()
	server := keeper.NewBlocklistQueryServer(k)

	// ACT: Attempt to query blocklist addresses with invalid request.
	_, err := server.Addresses(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query blocklist addresses with no state.
	res, err := server.Addresses(ctx, &blocklist.QueryAddresses{})
	// ASSERT: The query should've succeeded, with empty addresses.
	require.NoError(t, err)
	require.Empty(t, res.Addresses)

	// ARRANGE: Set blocklist addresses in state.
	user1, user2 := utils.TestAccount(), utils.TestAccount()
	k.SetBlockedAddress(ctx, user1.Bytes)
	k.SetBlockedAddress(ctx, user2.Bytes)

	// ACT: Attempt to query blocklist addresses with state.
	res, err = server.Addresses(ctx, &blocklist.QueryAddresses{})
	// ASSERT: The query should've succeeded, with addresses.
	require.NoError(t, err)
	require.Len(t, res.Addresses, 2)
	require.Contains(t, res.Addresses, user1.Address)
	require.Contains(t, res.Addresses, user2.Address)
}

func TestBlocklistAddressQuery(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewBlocklistQueryServer(k)

	// ACT: Attempt to query blocked state with invalid request.
	_, err := server.Address(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query blocked state with invalid address.
	_, err = server.Address(ctx, &blocklist.QueryAddress{
		Address: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
	})
	// ASSERT: The query should've failed due to invalid address.
	require.ErrorContains(t, err, "unable to decode address")

	// ACT: Attempt to query blocked state of unblocked address.
	res, err := server.Address(ctx, &blocklist.QueryAddress{
		Address: utils.TestAccount().Address,
	})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.False(t, res.Blocked)

	// ARRANGE: Set blocklist address in state.
	user := utils.TestAccount()
	k.SetBlockedAddress(ctx, user.Bytes)

	// ACT: Attempt to query blocked state of blocked address.
	res, err = server.Address(ctx, &blocklist.QueryAddress{
		Address: user.Address,
	})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.True(t, res.Blocked)
}
