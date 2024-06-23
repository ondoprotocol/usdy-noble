package e2e

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	"github.com/strangelove-ventures/interchaintest/v4/ibc"
	"github.com/strangelove-ventures/interchaintest/v4/testutil"
	"github.com/stretchr/testify/require"
)

var ONE = sdk.NewInt(1_000_000_000_000_000_000)

func TestMintBurn(t *testing.T) {
	t.Parallel()

	var wrapper Wrapper
	ctx := Suite(t, &wrapper, false)
	validator := wrapper.chain.Validators[0]

	// ASSERT: Minter has an allowance of 1 $USDY.
	EnsureMinter(t, wrapper, ctx, wrapper.minter.FormattedAddress(), ONE)

	// ACT: Mint 1 $USDY to Alice.
	_, err := validator.ExecTx(
		ctx, wrapper.minter.KeyName(),
		"aura", "mint", wrapper.alice.FormattedAddress(), ONE.String(),
	)
	require.NoError(t, err)

	// ASSERT: Alice has 1 $USDY.
	balance, err := wrapper.chain.GetBalance(ctx, wrapper.alice.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Equal(t, ONE.Int64(), balance)
	// ASSERT: Minter has no allowance.
	EnsureMinter(t, wrapper, ctx, wrapper.minter.FormattedAddress(), sdk.ZeroInt())

	// ASSERT: Burner has an allowance of 1 $USDY.
	EnsureBurner(t, wrapper, ctx, wrapper.burner.FormattedAddress(), ONE)

	// ACT: Burn 1 $USDY from Alice.
	_, err = validator.ExecTx(
		ctx, wrapper.burner.KeyName(),
		"aura", "burn", wrapper.alice.FormattedAddress(), ONE.String(),
	)
	require.NoError(t, err)

	// ASSERT: Alice has no balance.
	balance, err = wrapper.chain.GetBalance(ctx, wrapper.alice.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Zero(t, balance)
	// ASSERT: Burner has no allowance.
	EnsureBurner(t, wrapper, ctx, wrapper.burner.FormattedAddress(), sdk.ZeroInt())
}

func TestMintTransferBurn(t *testing.T) {
	t.Parallel()

	var wrapper Wrapper
	ctx := Suite(t, &wrapper, false)
	validator := wrapper.chain.Validators[0]

	// ASSERT: Minter has an allowance of 1 $USDY.
	EnsureMinter(t, wrapper, ctx, wrapper.minter.FormattedAddress(), ONE)

	// ACT: Mint 1 $USDY to Alice.
	_, err := validator.ExecTx(
		ctx, wrapper.minter.KeyName(),
		"aura", "mint", wrapper.alice.FormattedAddress(), ONE.String(),
	)
	require.NoError(t, err)

	// ASSERT: Alice has 1 $USDY.
	balance, err := wrapper.chain.GetBalance(ctx, wrapper.alice.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Equal(t, ONE.Int64(), balance)
	// ASSERT: Minter has no allowance.
	EnsureMinter(t, wrapper, ctx, wrapper.minter.FormattedAddress(), sdk.ZeroInt())

	// ACT: Transfer 1 $USDY from Alice to Bob.
	err = validator.SendFunds(ctx, wrapper.alice.KeyName(), ibc.WalletAmount{
		Address: wrapper.bob.FormattedAddress(),
		Denom:   "ausdy",
		Amount:  ONE.Int64(),
	})
	require.NoError(t, err)

	// ASSERT: Alice has no balance.
	balance, err = wrapper.chain.GetBalance(ctx, wrapper.alice.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Zero(t, balance)
	// ASSERT: Bob has 1 $USDY.
	balance, err = wrapper.chain.GetBalance(ctx, wrapper.bob.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Equal(t, ONE.Int64(), balance)

	// ASSERT: Burner has an allowance of 1 $USDY.
	EnsureBurner(t, wrapper, ctx, wrapper.burner.FormattedAddress(), ONE)

	// ACT: Burn 1 $USDY from Bob.
	_, err = validator.ExecTx(
		ctx, wrapper.burner.KeyName(),
		"aura", "burn", wrapper.bob.FormattedAddress(), ONE.String(),
	)
	require.NoError(t, err)

	// ASSERT: Bob has no balance.
	balance, err = wrapper.chain.GetBalance(ctx, wrapper.bob.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Zero(t, balance)
	// ASSERT: Burner has no allowance.
	EnsureBurner(t, wrapper, ctx, wrapper.burner.FormattedAddress(), sdk.ZeroInt())
}

func TestMintTransferBlockBurn(t *testing.T) {
	t.Parallel()

	var wrapper Wrapper
	ctx := Suite(t, &wrapper, false)
	validator := wrapper.chain.Validators[0]

	// ASSERT: Minter has an allowance of 1 $USDY.
	EnsureMinter(t, wrapper, ctx, wrapper.minter.FormattedAddress(), ONE)

	// ACT: Mint 1 $USDY to Alice.
	_, err := validator.ExecTx(
		ctx, wrapper.minter.KeyName(),
		"aura", "mint", wrapper.alice.FormattedAddress(), ONE.String(),
	)
	require.NoError(t, err)

	// ASSERT: Alice has 1 $USDY.
	balance, err := wrapper.chain.GetBalance(ctx, wrapper.alice.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Equal(t, ONE.Int64(), balance)
	// ASSERT: Minter has no allowance.
	EnsureMinter(t, wrapper, ctx, wrapper.minter.FormattedAddress(), sdk.ZeroInt())

	// ACT: Transfer 1 $USDY from Alice to Bob.
	err = validator.SendFunds(ctx, wrapper.alice.KeyName(), ibc.WalletAmount{
		Address: wrapper.bob.FormattedAddress(),
		Denom:   "ausdy",
		Amount:  ONE.Int64(),
	})
	require.NoError(t, err)

	// ASSERT: Alice has no balance.
	balance, err = wrapper.chain.GetBalance(ctx, wrapper.alice.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Zero(t, balance)
	// ASSERT: Bob has 1 $USDY.
	balance, err = wrapper.chain.GetBalance(ctx, wrapper.bob.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Equal(t, ONE.Int64(), balance)

	// ACT: Add Bob to the blocklist.
	_, err = validator.ExecTx(
		ctx, wrapper.owner.KeyName(),
		"aura", "blocklist", "add-to-blocklist", wrapper.bob.FormattedAddress(),
	)
	require.NoError(t, err)

	// ASSERT: Bob is blocked.
	EnsureBlocked(t, wrapper, ctx, wrapper.bob.FormattedAddress())

	// ACT: Attempt to transfer 1 $USDY from Bob to Alice.
	err = validator.SendFunds(ctx, wrapper.bob.KeyName(), ibc.WalletAmount{
		Address: wrapper.alice.FormattedAddress(),
		Denom:   "ausdy",
		Amount:  ONE.Int64(),
	})
	require.ErrorContains(t, err, "blocked from sending")

	// ASSERT: Bob has 1 $USDY.
	balance, err = wrapper.chain.GetBalance(ctx, wrapper.bob.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Equal(t, ONE.Int64(), balance)
	// ASSERT: Burner has an allowance of 1 $USDY.
	EnsureBurner(t, wrapper, ctx, wrapper.burner.FormattedAddress(), ONE)

	// ACT: Burn 1 $USDY from Bob.
	_, err = validator.ExecTx(
		ctx, wrapper.burner.KeyName(),
		"aura", "burn", wrapper.bob.FormattedAddress(), ONE.String(),
	)
	require.NoError(t, err)

	// ASSERT: Bob has no balance.
	balance, err = wrapper.chain.GetBalance(ctx, wrapper.bob.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Zero(t, balance)
	// ASSERT: Burner has no allowance.
	EnsureBurner(t, wrapper, ctx, wrapper.burner.FormattedAddress(), sdk.ZeroInt())
}

func TestIBCTransfer(t *testing.T) {
	t.Parallel()

	var wrapper Wrapper
	ctx := Suite(t, &wrapper, true)
	validator := wrapper.chain.Validators[0]

	// ARRANGE: Mint 1 $USDY to Alice.
	_, err := validator.ExecTx(
		ctx, wrapper.minter.KeyName(),
		"aura", "mint", wrapper.alice.FormattedAddress(), ONE.String(),
	)
	require.NoError(t, err)
	// ARRANGE: Determine Bob's external address.
	recipient := sdk.MustBech32ifyAddressBytes(wrapper.gaia.Config().Bech32Prefix, wrapper.bob.Address())

	// ACT: Attempt to transfer out of Noble.
	_, err = validator.SendIBCTransfer(ctx, "channel-0", wrapper.alice.KeyName(), ibc.WalletAmount{
		Address: recipient,
		Denom:   "ausdy",
		Amount:  ONE.Int64(),
	}, ibc.TransferOptions{})
	// ASSERT: The action should've failed due to not allowed channel.
	require.ErrorContains(t, err, "ausdy cannot be transferred over channel-0")

	// ACT: Allow transfers over channel-0.
	_, err = validator.ExecTx(
		ctx, wrapper.owner.KeyName(),
		"aura", "allow-channel", "channel-0",
	)
	require.NoError(t, err)
	// ASSERT: channel-0 has been allowed.
	EnsureChannel(t, wrapper, ctx, "channel-0")

	// ACT: Transfer 1 $USDY out of Noble.
	_, err = validator.SendIBCTransfer(ctx, "channel-0", wrapper.alice.KeyName(), ibc.WalletAmount{
		Address: recipient,
		Denom:   "ausdy",
		Amount:  ONE.Int64(),
	}, ibc.TransferOptions{})
	require.NoError(t, err)

	// ACT: Wait 10 blocks for packet to be relayed.
	require.NoError(t, testutil.WaitForBlocks(ctx, 10, wrapper.gaia))

	// ASSERT: Alice has no balance.
	balance, err := wrapper.chain.GetBalance(ctx, wrapper.alice.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Zero(t, balance)
	// ASSERT: Bob has 1 $USDY.
	balance, err = wrapper.gaia.GetBalance(ctx, recipient, transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: "ausdy",
	}.IBCDenom())
	require.NoError(t, err)
	require.Equal(t, ONE.Int64(), balance)
}
