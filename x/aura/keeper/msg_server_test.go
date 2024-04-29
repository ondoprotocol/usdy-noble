package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/aura/utils"
	"github.com/noble-assets/aura/utils/mocks"
	"github.com/noble-assets/aura/x/aura/keeper"
	"github.com/noble-assets/aura/x/aura/types"
	"github.com/stretchr/testify/require"
)

var ONE = math.NewInt(1_000_000_000_000_000_000)

func TestBurn(t *testing.T) {
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.AuraKeeperWithBank(t, bank)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set burner in state.
	burner := utils.TestAccount()
	require.NoError(t, k.Burners.Set(ctx, burner.Address))

	// ACT: Attempt to burn with invalid signer.
	_, err := server.Burn(ctx, &types.MsgBurn{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidBurner.Error())

	// ACT: Attempt to burn with invalid account address.
	_, err = server.Burn(ctx, &types.MsgBurn{
		Signer: burner.Address,
		From:   "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
	})
	// ASSERT: The action should've failed due to invalid account address.
	require.ErrorContains(t, err, "unable to decode account address")

	// ARRANGE: Generate a user account.
	user := utils.TestAccount()

	// ACT: Attempt to burn from user with insufficient funds.
	_, err = server.Burn(ctx, &types.MsgBurn{
		Signer: burner.Address,
		From:   user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to insufficient funds.
	require.ErrorContains(t, err, "unable to transfer from user to module")

	// ARRANGE: Give user 1 $USDY.
	bank.Balances[user.Address] = sdk.NewCoins(sdk.NewCoin(k.Denom, ONE))

	// ACT: Attempt to burn.
	_, err = server.Burn(ctx, &types.MsgBurn{
		Signer: burner.Address,
		From:   user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, bank.Balances[user.Address].IsZero())
	require.True(t, bank.Balances[types.ModuleName].IsZero())
}

func TestMint(t *testing.T) {
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.AuraKeeperWithBank(t, bank)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set minter in state.
	minter := utils.TestAccount()
	require.NoError(t, k.Minters.Set(ctx, minter.Address))

	// ACT: Attempt to mint with invalid signer.
	_, err := server.Mint(ctx, &types.MsgMint{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidMinter.Error())

	// ACT: Attempt to mint with invalid account address.
	_, err = server.Mint(ctx, &types.MsgMint{
		Signer: minter.Address,
		To:     "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
	})
	// ASSERT: The action should've failed due to invalid account address.
	require.ErrorContains(t, err, "unable to decode account address")

	// ARRANGE: Generate a user account.
	user := utils.TestAccount()

	// ACT: Attempt to mint.
	_, err = server.Mint(ctx, &types.MsgMint{
		Signer: minter.Address,
		To:     user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, ONE, bank.Balances[user.Address].AmountOf(k.Denom))
	require.True(t, bank.Balances[types.ModuleName].IsZero())
}

func TestPause(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set pauser in state.
	pauser := utils.TestAccount()
	require.NoError(t, k.Pausers.Set(ctx, pauser.Address))

	// ACT: Attempt to pause with invalid signer.
	_, err := server.Pause(ctx, &types.MsgPause{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidPauser.Error())
	paused, _ := k.Paused.Get(ctx)
	require.False(t, paused)

	// ACT: Attempt to pause.
	_, err = server.Pause(ctx, &types.MsgPause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	paused, _ = k.Paused.Get(ctx)
	require.True(t, paused)

	// ACT: Attempt to pause again.
	_, err = server.Pause(ctx, &types.MsgPause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've failed due to module being paused already.
	require.ErrorContains(t, err, "module is already paused")
	paused, _ = k.Paused.Get(ctx)
	require.True(t, paused)
}

func TestUnpause(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set paused state to true.
	require.NoError(t, k.Paused.Set(ctx, true))
	// ARRANGE: Set pauser in state.
	pauser := utils.TestAccount()
	require.NoError(t, k.Pausers.Set(ctx, pauser.Address))

	// ACT: Attempt to unpause with invalid signer.
	_, err := server.Unpause(ctx, &types.MsgUnpause{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidPauser.Error())
	paused, _ := k.Paused.Get(ctx)
	require.True(t, paused)

	// ACT: Attempt to unpause.
	_, err = server.Unpause(ctx, &types.MsgUnpause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	paused, _ = k.Paused.Get(ctx)
	require.False(t, paused)

	// ACT: Attempt to unpause again.
	_, err = server.Unpause(ctx, &types.MsgUnpause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've failed due to module being unpaused already.
	require.ErrorContains(t, err, "module is already unpaused")
	paused, _ = k.Paused.Get(ctx)
	require.False(t, paused)
}
