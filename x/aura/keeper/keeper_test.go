package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/aura/utils"
	"github.com/noble-assets/aura/utils/mocks"
	"github.com/noble-assets/aura/x/aura/types"
	"github.com/stretchr/testify/require"
)

func TestSendRestriction(t *testing.T) {
	keeper, ctx := mocks.AuraKeeper(t)
	coins := sdk.NewCoins(sdk.NewCoin(
		keeper.Denom, ONE,
	))

	// ACT: Attempt to send different token.
	_, err := keeper.SendRestrictionFn(ctx, utils.TestAccount().Bytes, utils.TestAccount().Bytes, sdk.NewCoins(sdk.NewCoin(
		"uusdc", math.NewInt(1_000_000),
	)))
	// ASSERT: The action should've succeeded due to different denom.
	require.NoError(t, err)

	// ACT: Attempt to send.
	_, err = keeper.SendRestrictionFn(ctx, utils.TestAccount().Bytes, utils.TestAccount().Bytes, coins)
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)

	// ARRANGE: Set paused state to true.
	require.NoError(t, keeper.Paused.Set(ctx, true))

	// ACT: Attempt to send when paused.
	_, err = keeper.SendRestrictionFn(ctx, utils.TestAccount().Bytes, utils.TestAccount().Bytes, coins)
	// ASSERT: The action should've failed due to module being paused.
	require.ErrorContains(t, err, "ausdy transfers are paused")

	// ARRANGE: Set paused state to false.
	require.NoError(t, keeper.Paused.Set(ctx, false))
	// ARRANGE: Generate a user account.
	user := utils.TestAccount()

	// ACT: Attempt to burn from non-blocklisted address.
	_, err = keeper.SendRestrictionFn(ctx, user.Bytes, types.ModuleAddress, coins)
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)

	// ARRANGE: Block user address.
	require.NoError(t, keeper.BlockedAddresses.Set(ctx, user.Bytes, true))

	// ACT: Attempt to burn from blocklisted address.
	_, err = keeper.SendRestrictionFn(ctx, user.Bytes, types.ModuleAddress, coins)
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)

	// ACT: Attempt to mint to blocklisted address.
	_, err = keeper.SendRestrictionFn(ctx, types.ModuleAddress, user.Bytes, coins)
	// ASSERT: The action shoudl've failed due to blocked recipient.
	require.ErrorContains(t, err, "blocked from receiving")

	// ARRANGE: Unblock user address.
	require.NoError(t, keeper.BlockedAddresses.Remove(ctx, user.Bytes))

	// ACT: Attempt to mint to non-blocklisted address.
	_, err = keeper.SendRestrictionFn(ctx, types.ModuleAddress, user.Bytes, coins)
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)

	// ARRANGE: Block user address.
	require.NoError(t, keeper.BlockedAddresses.Set(ctx, user.Bytes, true))

	// ACT: Attempt to send from blocklisted address.
	_, err = keeper.SendRestrictionFn(ctx, user.Bytes, utils.TestAccount().Bytes, coins)
	// ASSERT: The action should've failed due to blocked sender.
	require.ErrorContains(t, err, "blocked from sending")

	// ACT: Attempt to send to blocklisted address.
	_, err = keeper.SendRestrictionFn(ctx, utils.TestAccount().Bytes, user.Bytes, coins)
	// ASSERT: The action should've failed due to blocked recipient.
	require.ErrorContains(t, err, "blocked from receiving")
}
