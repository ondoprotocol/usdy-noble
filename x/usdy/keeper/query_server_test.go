package keeper_test

import (
	"testing"

	"cosmossdk.io/collections"
	"github.com/noble-assets/ondo/utils"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/noble-assets/ondo/utils/mocks"
	"github.com/noble-assets/ondo/x/usdy/keeper"
	"github.com/noble-assets/ondo/x/usdy/types"
	"github.com/stretchr/testify/require"
)

func TestDenomQuery(t *testing.T) {
	k, ctx := mocks.USDYKeeper(t)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query denom with invalid request.
	_, err := server.Denom(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query denom.
	res, err := server.Denom(ctx, &types.QueryDenom{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Equal(t, "ausdy", res.Denom)
}

func TestPausedQuery(t *testing.T) {
	k, ctx := mocks.USDYKeeper(t)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query paused state with invalid request.
	_, err := server.Paused(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query paused state with no state.
	req, err := server.Paused(ctx, &types.QueryPaused{})
	// ASSERT: The query should've succeeded, and returned false.
	require.NoError(t, err)
	require.False(t, req.Paused)

	// ARRANGE: Set paused state to true.
	require.NoError(t, k.Paused.Set(ctx, true))

	// ACT: Attempt to query paused state with state.
	req, err = server.Paused(ctx, &types.QueryPaused{})
	// ASSERT: The query should've succeeded, and returned true.
	require.NoError(t, err)
	require.True(t, req.Paused)
}

func TestPauserQuery(t *testing.T) {
	k, ctx := mocks.USDYKeeper(t)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query pauser with invalid request.
	_, err := server.Pauser(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query pauser with no state.
	_, err = server.Pauser(ctx, &types.QueryPauser{})
	// ASSERT: The query should've failed.
	require.ErrorContains(t, err, collections.ErrNotFound.Error())

	// ARRANGE: Set pauser in state.
	pauser := utils.TestAccount()
	require.NoError(t, k.Pauser.Set(ctx, pauser.Address))

	// ACT: Attempt to query pauser with state.
	req, err := server.Pauser(ctx, &types.QueryPauser{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Equal(t, pauser.Address, req.Pauser)
}
