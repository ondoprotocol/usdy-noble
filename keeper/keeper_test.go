package keeper_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ondoprotocol/usdy-noble/v2/keeper"
	"github.com/ondoprotocol/usdy-noble/v2/types"
	"github.com/ondoprotocol/usdy-noble/v2/utils"
	"github.com/ondoprotocol/usdy-noble/v2/utils/mocks"
	"github.com/stretchr/testify/require"
)

func TestSendRestrictionBurn(t *testing.T) {
	user := utils.TestAccount()
	k, ctx := mocks.AuraKeeper()
	coins := sdk.NewCoins(sdk.NewCoin(k.Denom, ONE))

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
			require.NoError(t, k.SetPaused(ctx, testCase.paused))
			// ARRANGE: Set blocked state.
			if testCase.blocked {
				require.NoError(t, k.SetBlockedAddress(ctx, user.Bytes))
			} else {
				require.NoError(t, k.DeleteBlockedAddress(ctx, user.Bytes))
			}

			// ACT: Attempt to burn.
			_, err := k.SendRestrictionFn(ctx, user.Bytes, types.ModuleAddress, coins)

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
	k, ctx := mocks.AuraKeeper()
	coins := sdk.NewCoins(sdk.NewCoin(k.Denom, ONE))

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
			require.NoError(t, k.SetPaused(ctx, testCase.paused))
			// ARRANGE: Set blocked state.
			if testCase.blocked {
				require.NoError(t, k.SetBlockedAddress(ctx, user.Bytes))
			} else {
				require.NoError(t, k.DeleteBlockedAddress(ctx, user.Bytes))
			}

			// ACT: Attempt to mint.
			_, err := k.SendRestrictionFn(ctx, types.ModuleAddress, user.Bytes, coins)

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
	k, ctx := mocks.AuraKeeper()
	coins := sdk.NewCoins(sdk.NewCoin(k.Denom, ONE))

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
			coins:            sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1_000_000))),
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
			require.NoError(t, k.SetPaused(ctx, testCase.paused))
			// ARRANGE: Set sender blocked state.
			if testCase.senderBlocked {
				require.NoError(t, k.SetBlockedAddress(ctx, alice.Bytes))
			} else {
				require.NoError(t, k.DeleteBlockedAddress(ctx, alice.Bytes))
			}
			// ARRANGE: Set recipient blocked state.
			if testCase.recipientBlocked {
				require.NoError(t, k.SetBlockedAddress(ctx, bob.Bytes))
			} else {
				require.NoError(t, k.DeleteBlockedAddress(ctx, bob.Bytes))
			}

			// ACT: Attempt to transfer.
			_, err := k.SendRestrictionFn(ctx, alice.Bytes, bob.Bytes, testCase.coins)

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
	k, ctx := mocks.AuraKeeper()
	coins := sdk.NewCoins(sdk.NewCoin(k.Denom, ONE))

	// ARRANGE: Set a blocked channel in state.
	require.NoError(t, k.SetBlockedChannel(ctx, "channel-0"))
	escrow := transfertypes.GetEscrowAddress(transfertypes.PortID, "channel-0")

	// ACT: Attempt to transfer from user to escrow account.
	// This is to mimic the underlying transfer that occurs when using IBC.
	_, err := k.SendRestrictionFn(ctx, user.Bytes, escrow, coins)

	// ASSERT: The action should've failed due to blocked channel.
	require.ErrorContains(t, err, "transfers are blocked")
}

func TestNewKeeper(t *testing.T) {
	// ARRANGE: Set the PausedKey to an already existing key
	types.PausedKey = types.OwnerKey

	// ACT: Attempt to create a new Keeper with overlapping prefixes
	require.Panics(t, func() {
		keeper.NewKeeper(
			"ausdy",
			mocks.FailingStore(mocks.Set, nil),
			runtime.ProvideEventService(),
			address.NewBech32Codec("noble"),
			mocks.BankKeeper{},
		)
	})
	// ASSERT: The function should've panicked.

	// ARRANGE: Restore the original PausedKey
	types.PausedKey = []byte("paused")
}
