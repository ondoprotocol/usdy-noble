package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	"github.com/noble-assets/aura/utils"
	"github.com/noble-assets/aura/utils/mocks"
	"github.com/noble-assets/aura/x/aura/types"
	"github.com/stretchr/testify/require"
)

func TestSendRestrictionBurn(t *testing.T) {
	user := utils.TestAccount()
	keeper, ctx := mocks.AuraKeeper(t)
	coins := sdk.NewCoins(sdk.NewCoin(keeper.Denom, ONE))

	testCases := []struct {
		name    string
		paused  bool
		blocked bool
		err     error
	}{
		{
			name:    "PausedAndBlocked",
			paused:  true,
			blocked: true,
			err:     nil,
		},
		{
			name:    "PausedAndUnblocked",
			paused:  true,
			blocked: false,
			err:     nil,
		},
		{
			name:    "UnpausedAndBlocked",
			paused:  false,
			blocked: true,
			err:     nil,
		},
		{
			name:    "UnpausedAndUnblocked",
			paused:  false,
			blocked: false,
			err:     nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// ARRANGE: Set paused state.
			keeper.SetPaused(ctx, testCase.paused)
			// ARRANGE: Set blocked state.
			if testCase.blocked {
				keeper.SetBlockedAddress(ctx, user.Bytes)
			} else {
				keeper.DeleteBlockedAddress(ctx, user.Bytes)
			}

			// ACT: Attempt to burn.
			_, err := keeper.SendRestrictionFn(ctx, user.Bytes, types.ModuleAddress, coins)

			// ASSERT: Send restriction correctly handled test case.
			if testCase.err != nil {
				require.ErrorContains(t, err, testCase.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSendRestrictionMint(t *testing.T) {
	user := utils.TestAccount()
	keeper, ctx := mocks.AuraKeeper(t)
	coins := sdk.NewCoins(sdk.NewCoin(keeper.Denom, ONE))

	testCases := []struct {
		name    string
		paused  bool
		blocked bool
		err     error
	}{
		{
			name:    "PausedAndBlocked",
			paused:  true,
			blocked: true,
			err:     errors.New("transfers are paused"),
		},
		{
			name:    "PausedAndUnblocked",
			paused:  true,
			blocked: false,
			err:     errors.New("transfers are paused"),
		},
		{
			name:    "UnpausedAndBlocked",
			paused:  false,
			blocked: true,
			err:     errors.New("blocked from receiving"),
		},
		{
			name:    "UnpausedAndUnblocked",
			paused:  false,
			blocked: false,
			err:     nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// ARRANGE: Set paused state.
			keeper.SetPaused(ctx, testCase.paused)
			// ARRANGE: Set blocked state.
			if testCase.blocked {
				keeper.SetBlockedAddress(ctx, user.Bytes)
			} else {
				keeper.DeleteBlockedAddress(ctx, user.Bytes)
			}

			// ACT: Attempt to mint.
			_, err := keeper.SendRestrictionFn(ctx, types.ModuleAddress, user.Bytes, coins)

			// ASSERT: Send restriction correctly handled test case.
			if testCase.err != nil {
				require.ErrorContains(t, err, testCase.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSendRestrictionTransfer(t *testing.T) {
	alice, bob := utils.TestAccount(), utils.TestAccount()
	keeper, ctx := mocks.AuraKeeper(t)
	coins := sdk.NewCoins(sdk.NewCoin(keeper.Denom, ONE))

	testCases := []struct {
		name             string
		paused           bool
		senderBlocked    bool
		recipientBlocked bool
		coins            sdk.Coins
		err              error
	}{
		{
			name:             "NonUSDYTransfer",
			paused:           true,
			senderBlocked:    true,
			recipientBlocked: true,
			coins:            sdk.NewCoins(sdk.NewCoin("uusdc", sdk.NewInt(1_000_000))),
			err:              nil,
		},
		{
			name:             "PausedAndSenderBlockedAndRecipientBlocked",
			paused:           true,
			senderBlocked:    true,
			recipientBlocked: true,
			coins:            coins,
			err:              errors.New("transfers are paused"),
		},
		{
			name:             "PausedAndSenderBlockedAndRecipientUnblocked",
			paused:           true,
			senderBlocked:    true,
			recipientBlocked: false,
			coins:            coins,
			err:              errors.New("transfers are paused"),
		},
		{
			name:             "PausedAndSenderUnblockedAndRecipientBlocked",
			paused:           true,
			senderBlocked:    false,
			recipientBlocked: true,
			coins:            coins,
			err:              errors.New("transfers are paused"),
		},
		{
			name:             "PausedAndSenderUnblockedAndRecipientUnblocked",
			paused:           true,
			senderBlocked:    false,
			recipientBlocked: false,
			coins:            coins,
			err:              errors.New("transfers are paused"),
		},
		{
			name:             "UnpausedAndSenderBlockedAndRecipientBlocked",
			paused:           false,
			senderBlocked:    true,
			recipientBlocked: true,
			coins:            coins,
			err:              errors.New("blocked from sending"),
		},
		{
			name:             "UnpausedAndSenderBlockedAndRecipientUnblocked",
			paused:           false,
			senderBlocked:    true,
			recipientBlocked: false,
			coins:            coins,
			err:              errors.New("blocked from sending"),
		},
		{
			name:             "UnpausedAndSenderUnblockedAndRecipientBlocked",
			paused:           false,
			senderBlocked:    false,
			recipientBlocked: true,
			coins:            coins,
			err:              errors.New("blocked from receiving"),
		},
		{
			name:             "UnpausedAndSenderUnblockedAndRecipientUnblocked",
			paused:           false,
			senderBlocked:    false,
			recipientBlocked: false,
			coins:            coins,
			err:              nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// ARRANGE: Set paused state.
			keeper.SetPaused(ctx, testCase.paused)
			// ARRANGE: Set sender blocked state.
			if testCase.senderBlocked {
				keeper.SetBlockedAddress(ctx, alice.Bytes)
			} else {
				keeper.DeleteBlockedAddress(ctx, alice.Bytes)
			}
			// ARRANGE: Set recipient blocked state.
			if testCase.recipientBlocked {
				keeper.SetBlockedAddress(ctx, bob.Bytes)
			} else {
				keeper.DeleteBlockedAddress(ctx, bob.Bytes)
			}

			// ACT: Attempt to transfer.
			_, err := keeper.SendRestrictionFn(ctx, alice.Bytes, bob.Bytes, testCase.coins)

			// ASSERT: Send restriction correctly handled test case.
			if testCase.err != nil {
				require.ErrorContains(t, err, testCase.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSendRestrictionIBCTransfer(t *testing.T) {
	user := utils.TestAccount()
	keeper, ctx := mocks.AuraKeeper(t)
	coins := sdk.NewCoins(sdk.NewCoin(keeper.Denom, ONE))

	// ARRANGE: Set a blocked channel in state.
	keeper.SetBlockedChannel(ctx, "channel-0")
	escrow := transfertypes.GetEscrowAddress(transfertypes.PortID, "channel-0")

	// ACT: Attempt to transfer from user to escrow account.
	// This is to mimic the underlying transfer that occurs when using IBC.
	_, err := keeper.SendRestrictionFn(ctx, user.Bytes, escrow, coins)

	// ASSERT: The action should've failed due to blocked channel.
	require.ErrorContains(t, err, "transfers are blocked")
}
