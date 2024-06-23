package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/noble-assets/aura/utils"
	"github.com/noble-assets/aura/utils/mocks"
	"github.com/noble-assets/aura/x/aura/keeper"
	"github.com/noble-assets/aura/x/aura/types"
	"github.com/stretchr/testify/require"
)

func TestDenomQuery(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query denom with invalid request.
	_, err := server.Denom(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query denom.
	res, err := server.Denom(goCtx, &types.QueryDenom{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Equal(t, "ausdy", res.Denom)
}

func TestPausedQuery(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query paused state with invalid request.
	_, err := server.Paused(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query paused state with no state.
	res, err := server.Paused(goCtx, &types.QueryPaused{})
	// ASSERT: The query should've succeeded, and returned false.
	require.NoError(t, err)
	require.False(t, res.Paused)

	// ARRANGE: Set paused state to true.
	k.SetPaused(ctx, true)

	// ACT: Attempt to query paused state with state.
	res, err = server.Paused(goCtx, &types.QueryPaused{})
	// ASSERT: The query should've succeeded, and returned true.
	require.NoError(t, err)
	require.True(t, res.Paused)
}

func TestOwnerQuery(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query owner with invalid request.
	_, err := server.Owner(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to query owner with state.
	res, err := server.Owner(goCtx, &types.QueryOwner{})
	// ASSERT: The query should've succeeded, with empty pending owner.
	require.NoError(t, err)
	require.Equal(t, owner.Address, res.Owner)
	require.Empty(t, res.PendingOwner)

	// ARRANGE: Set pending owner in state.
	pendingOwner := utils.TestAccount()
	k.SetPendingOwner(ctx, pendingOwner.Address)

	// ACT: Attempt to query owner with state.
	res, err = server.Owner(goCtx, &types.QueryOwner{})
	// ASSERT: The query should've succeeded, with pending owner.
	require.NoError(t, err)
	require.Equal(t, owner.Address, res.Owner)
	require.Equal(t, pendingOwner.Address, res.PendingOwner)
}

func TestBurnersQuery(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query burners with invalid request.
	_, err := server.Burners(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query burners with no state.
	res, err := server.Burners(goCtx, &types.QueryBurners{})
	// ASSERT: The query should've succeeded, and returned no burners.
	require.NoError(t, err)
	require.Empty(t, res.Burners)

	// ARRANGE: Set burners in state.
	burner1, burner2 := utils.TestAccount(), utils.TestAccount()
	k.SetBurner(ctx, burner1.Address, ONE)
	k.SetBurner(ctx, burner2.Address, ONE.MulRaw(2))

	// ACT: Attempt to query burners with state.
	res, err = server.Burners(goCtx, &types.QueryBurners{})
	// ASSERT: The query should've succeeded, and returned burners.
	require.NoError(t, err)
	require.Len(t, res.Burners, 2)
	require.Contains(t, res.Burners, types.Burner{
		Address:   burner1.Address,
		Allowance: ONE,
	})
	require.Contains(t, res.Burners, types.Burner{
		Address:   burner2.Address,
		Allowance: ONE.MulRaw(2),
	})
}

func TestMintersQuery(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query minters with invalid request.
	_, err := server.Minters(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query minters with no state.
	res, err := server.Minters(goCtx, &types.QueryMinters{})
	// ASSERT: The query should've succeeded, and returned no minters.
	require.NoError(t, err)
	require.Empty(t, res.Minters)

	// ARRANGE: Set minters in state.
	minter1, minter2 := utils.TestAccount(), utils.TestAccount()
	k.SetMinter(ctx, minter1.Address, ONE)
	k.SetMinter(ctx, minter2.Address, ONE.MulRaw(2))

	// ACT: Attempt to query minters with state.
	res, err = server.Minters(goCtx, &types.QueryMinters{})
	// ASSERT: The query should've succeeded, and returned minters.
	require.NoError(t, err)
	require.Len(t, res.Minters, 2)
	require.Contains(t, res.Minters, types.Minter{
		Address:   minter1.Address,
		Allowance: ONE,
	})
	require.Contains(t, res.Minters, types.Minter{
		Address:   minter2.Address,
		Allowance: ONE.MulRaw(2),
	})
}

func TestPausersQuery(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query pausers with invalid request.
	_, err := server.Pausers(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query pausers with no state.
	res, err := server.Pausers(goCtx, &types.QueryPausers{})
	// ASSERT: The query should've succeeded, and returned no pausers.
	require.NoError(t, err)
	require.Empty(t, res.Pausers)

	// ARRANGE: Set pausers in state.
	pauser1, pauser2 := utils.TestAccount(), utils.TestAccount()
	k.SetPauser(ctx, pauser1.Address)
	k.SetPauser(ctx, pauser2.Address)

	// ACT: Attempt to query pausers with state.
	res, err = server.Pausers(goCtx, &types.QueryPausers{})
	// ASSERT: The query should've succeeded, and returned pausers.
	require.NoError(t, err)
	require.Len(t, res.Pausers, 2)
	require.Contains(t, res.Pausers, pauser1.Address)
	require.Contains(t, res.Pausers, pauser2.Address)
}

func TestChannelsQuery(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query channels with invalid request.
	_, err := server.Channels(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query channels with no state.
	res, err := server.Channels(goCtx, &types.QueryChannels{})
	// ASSERT: The query should've succeeded, and returned no channels.
	require.NoError(t, err)
	require.Empty(t, res.Channels)

	// ARRANGE: Set channels in state.
	channel1, channel2 := "channel-0", "channel-1"
	k.SetChannel(ctx, channel1)
	k.SetChannel(ctx, channel2)

	// ACT: Attempt to query channels with state.
	res, err = server.Channels(goCtx, &types.QueryChannels{})
	// ASSERT: The query should've succeeded, and returned channels.
	require.NoError(t, err)
	require.Len(t, res.Channels, 2)
	require.Contains(t, res.Channels, channel1)
	require.Contains(t, res.Channels, channel2)
}
