package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/aura/utils"
	"github.com/noble-assets/aura/utils/mocks"
	"github.com/stretchr/testify/require"
)

func TestSendRestriction(t *testing.T) {
	keeper, ctx := mocks.AuraKeeper(t)
	from, to := utils.TestAccount(), utils.TestAccount()

	// ACT: Attempt to send 1 USDC.
	_, err := keeper.SendRestrictionFn(ctx, from.Bytes, to.Bytes, sdk.NewCoins(sdk.NewCoin(
		"uusdc", math.NewInt(1_000_000),
	)))
	// ASSERT: The action should've succeeded due to different denom.
	require.NoError(t, err)

	// ACT: Attempt to send 1 USDY.
	_, err = keeper.SendRestrictionFn(ctx, from.Bytes, to.Bytes, sdk.NewCoins(sdk.NewCoin(
		keeper.Denom, ONE,
	)))
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)

	// ARRANGE: Set paused state to true.
	require.NoError(t, keeper.Paused.Set(ctx, true))

	// ACT: Attempt to send 1 USDY when paused.
	_, err = keeper.SendRestrictionFn(ctx, from.Bytes, to.Bytes, sdk.NewCoins(sdk.NewCoin(
		keeper.Denom, ONE,
	)))
	// ASSERT: The action should've failed due to module being paused.
	require.ErrorContains(t, err, "ausdy transfers are paused")

	// ARRANGE: Set paused state to false.
	require.NoError(t, keeper.Paused.Set(ctx, false))
	// ARRANGE: Block from address.
	require.NoError(t, keeper.BlockedAddresses.Set(ctx, from.Bytes, true))

	// ACT: Attempt to send 1 USDY from blocklisted address.
	_, err = keeper.SendRestrictionFn(ctx, from.Bytes, to.Bytes, sdk.NewCoins(sdk.NewCoin(
		keeper.Denom, ONE,
	)))
	// ASSERT: The action should've failed due to blocked sender.
	require.ErrorContains(t, err, "blocked from sending")

	// ARRANGE: Unblock from address.
	require.NoError(t, keeper.BlockedAddresses.Remove(ctx, from.Bytes))
	// ARRANGE: Block to address.
	require.NoError(t, keeper.BlockedAddresses.Set(ctx, to.Bytes, true))

	// ACT: Attempt to send 1 USDY to blocklisted address.
	_, err = keeper.SendRestrictionFn(ctx, from.Bytes, to.Bytes, sdk.NewCoins(sdk.NewCoin(
		keeper.Denom, ONE,
	)))
	// ASSERT: The action should've failed due to blocked recipient.
	require.ErrorContains(t, err, "blocked from receiving")
}
