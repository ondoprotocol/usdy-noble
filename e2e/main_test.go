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
	denom := transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: "ausdy",
	}.IBCDenom()

	// ARRANGE: Mint 2 $USDY to Alice.
	_, err := validator.ExecTx(
		ctx, wrapper.owner.KeyName(),
		"aura", "set-minter-allowance", wrapper.minter.FormattedAddress(), ONE.MulRaw(2).String(),
	)
	require.NoError(t, err)
	_, err = validator.ExecTx(
		ctx, wrapper.minter.KeyName(),
		"aura", "mint", wrapper.alice.FormattedAddress(), ONE.MulRaw(2).String(),
	)
	require.NoError(t, err)

	// ACT: Attempt to transfer out of Noble, channel is allowed.
	_, err = wrapper.chain.SendIBCTransfer(ctx, "channel-0", wrapper.alice.KeyName(), ibc.WalletAmount{
		Address: wrapper.charlie.FormattedAddress(),
		Denom:   "ausdy",
		Amount:  ONE.MulRaw(2).Int64(),
	}, ibc.TransferOptions{})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.NoError(t, testutil.WaitForBlocks(ctx, 5, wrapper.chain, wrapper.gaia))
	balance, err := wrapper.chain.GetBalance(ctx, wrapper.alice.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Zero(t, balance)
	balance, err = wrapper.gaia.GetBalance(ctx, wrapper.charlie.FormattedAddress(), denom)
	require.NoError(t, err)
	require.Equal(t, ONE.MulRaw(2).Int64(), balance)

	// ACT: Attempt to transfer back to Noble, channel is allowed.
	_, err = wrapper.gaia.SendIBCTransfer(ctx, "channel-0", wrapper.charlie.KeyName(), ibc.WalletAmount{
		Address: wrapper.alice.FormattedAddress(),
		Denom:   denom,
		Amount:  ONE.Int64(),
	}, ibc.TransferOptions{})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.NoError(t, testutil.WaitForBlocks(ctx, 5, wrapper.chain, wrapper.gaia))
	balance, err = wrapper.chain.GetBalance(ctx, wrapper.alice.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Equal(t, ONE.Int64(), balance)
	balance, err = wrapper.gaia.GetBalance(ctx, wrapper.charlie.FormattedAddress(), denom)
	require.NoError(t, err)
	require.Equal(t, ONE.Int64(), balance)

	// ACT: Block transfers over channel-0.
	_, err = validator.ExecTx(
		ctx, wrapper.owner.KeyName(),
		"aura", "add-blocked-channel", "channel-0",
	)
	require.NoError(t, err)
	// ASSERT: channel-0 has been blocked.
	EnsureBlockedChannel(t, wrapper, ctx, "channel-0")

	// ACT: Attempt to transfer out of Noble, channel is blocked.
	_, err = wrapper.chain.SendIBCTransfer(ctx, "channel-0", wrapper.alice.KeyName(), ibc.WalletAmount{
		Address: wrapper.charlie.FormattedAddress(),
		Denom:   "ausdy",
		Amount:  ONE.Int64(),
	}, ibc.TransferOptions{})
	// ASSERT: The action should've failed.
	require.ErrorContains(t, err, "ausdy transfers are blocked on channel-0")

	// ACT: Attempt to transfer back to Noble, channel is allowed.
	_, err = wrapper.gaia.SendIBCTransfer(ctx, "channel-0", wrapper.charlie.KeyName(), ibc.WalletAmount{
		Address: wrapper.alice.FormattedAddress(),
		Denom:   denom,
		Amount:  ONE.Int64(),
	}, ibc.TransferOptions{})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.NoError(t, testutil.WaitForBlocks(ctx, 5, wrapper.chain, wrapper.gaia))
	balance, err = wrapper.chain.GetBalance(ctx, wrapper.alice.FormattedAddress(), "ausdy")
	require.NoError(t, err)
	require.Equal(t, ONE.MulRaw(2).Int64(), balance)
	balance, err = wrapper.gaia.GetBalance(ctx, wrapper.charlie.FormattedAddress(), denom)
	require.NoError(t, err)
	require.Zero(t, balance)
}
