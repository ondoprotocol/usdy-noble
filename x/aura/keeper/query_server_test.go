package keeper_test

import (
	"testing"

	"cosmossdk.io/collections"
	"github.com/noble-assets/aura/utils"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/noble-assets/aura/utils/mocks"
	"github.com/noble-assets/aura/x/aura/keeper"
	"github.com/noble-assets/aura/x/aura/types"
	"github.com/stretchr/testify/require"
)

func TestDenomQuery(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
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
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query paused state with invalid request.
	_, err := server.Paused(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query paused state with no state.
	res, err := server.Paused(ctx, &types.QueryPaused{})
	// ASSERT: The query should've succeeded, and returned false.
	require.NoError(t, err)
	require.False(t, res.Paused)

	// ARRANGE: Set paused state to true.
	require.NoError(t, k.Paused.Set(ctx, true))

	// ACT: Attempt to query paused state with state.
	res, err = server.Paused(ctx, &types.QueryPaused{})
	// ASSERT: The query should've succeeded, and returned true.
	require.NoError(t, err)
	require.True(t, res.Paused)
}

func TestBurnersQuery(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query burners with invalid request.
	_, err := server.Burners(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query burners with no state.
	res, err := server.Burners(ctx, &types.QueryBurners{})
	// ASSERT: The query should've succeeded, and returned no burners.
	require.NoError(t, err)
	require.Empty(t, res.Burners)

	// ARRANGE: Set burners in state.
	burner1, burner2 := utils.TestAccount(), utils.TestAccount()
	require.NoError(t, k.Burners.Set(ctx, burner1.Address))
	require.NoError(t, k.Burners.Set(ctx, burner2.Address))

	// ACT: Attempt to query burners with state.
	res, err = server.Burners(ctx, &types.QueryBurners{})
	// ASSERT: The query should've succeeded, and returned burners.
	require.NoError(t, err)
	require.Len(t, res.Burners, 2)
	require.Contains(t, res.Burners, burner1.Address)
	require.Contains(t, res.Burners, burner2.Address)
}

func TestMintersQuery(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query minters with invalid request.
	_, err := server.Minters(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query minters with no state.
	res, err := server.Minters(ctx, &types.QueryMinters{})
	// ASSERT: The query should've succeeded, and returned no minters.
	require.NoError(t, err)
	require.Empty(t, res.Minters)

	// ARRANGE: Set minters in state.
	minter1, minter2 := utils.TestAccount(), utils.TestAccount()
	require.NoError(t, k.Minters.Set(ctx, minter1.Address))
	require.NoError(t, k.Minters.Set(ctx, minter2.Address))

	// ACT: Attempt to query minters with state.
	res, err = server.Minters(ctx, &types.QueryMinters{})
	// ASSERT: The query should've succeeded, and returned minters.
	require.NoError(t, err)
	require.Len(t, res.Minters, 2)
	require.Contains(t, res.Minters, minter1.Address)
	require.Contains(t, res.Minters, minter2.Address)
}

func TestPauserQuery(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
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
	res, err := server.Pauser(ctx, &types.QueryPauser{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Equal(t, pauser.Address, res.Pauser)
}
