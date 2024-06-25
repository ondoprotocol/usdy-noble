package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/aura/utils"
	"github.com/ondoprotocol/aura/utils/mocks"
	"github.com/ondoprotocol/aura/x/aura/keeper"
	"github.com/ondoprotocol/aura/x/aura/types"
	"github.com/stretchr/testify/require"
)

var ONE = sdk.NewInt(1_000_000_000_000_000_000)

func TestBurn(t *testing.T) {
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.AuraKeeperWithBank(t, bank)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set burner in state, with enough allowance for a single burn.
	burner := utils.TestAccount()
	k.SetBurner(ctx, burner.Address, ONE)

	// ACT: Attempt to burn with invalid signer.
	_, err := server.Burn(goCtx, &types.MsgBurn{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidBurner.Error())

	// ACT: Attempt to burn with invalid account address.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Signer: burner.Address,
		From:   "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to invalid account address.
	require.ErrorContains(t, err, "unable to decode account address")

	// ARRANGE: Generate a user account.
	user := utils.TestAccount()

	// ACT: Attempt to burn invalid amount.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Signer: burner.Address,
		From:   user.Address,
		Amount: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to invalid amount.
	require.ErrorContains(t, err, "amount must be positive")

	// ACT: Attempt to burn from user with insufficient funds.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Signer: burner.Address,
		From:   user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to insufficient funds.
	require.ErrorContains(t, err, "unable to transfer from user to module")

	// ARRANGE: Give user 1 $USDY.
	bank.Balances[user.Address] = sdk.NewCoins(sdk.NewCoin(k.Denom, ONE))

	// ACT: Attempt to burn.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Signer: burner.Address,
		From:   user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, bank.Balances[user.Address].IsZero())
	require.True(t, bank.Balances[types.ModuleName].IsZero())
	require.True(t, k.GetBurner(ctx, burner.Address).IsZero())

	// ACT: Attempt another burn with insufficient allowance.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Signer: burner.Address,
		From:   user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to insufficient allowance.
	require.ErrorContains(t, err, types.ErrInsufficientAllowance.Error())
}

func TestMint(t *testing.T) {
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.AuraKeeperWithBank(t, bank)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set minter in state, with enough allowance for a single mint.
	minter := utils.TestAccount()
	k.SetMinter(ctx, minter.Address, ONE)

	// ACT: Attempt to mint with invalid signer.
	_, err := server.Mint(goCtx, &types.MsgMint{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidMinter.Error())

	// ACT: Attempt to mint with invalid account address.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Signer: minter.Address,
		To:     "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to invalid account address.
	require.ErrorContains(t, err, "unable to decode account address")

	// ARRANGE: Generate a user account and add to blocklist.
	user := utils.TestAccount()
	k.SetBlockedAddress(ctx, user.Bytes)

	// ACT: Attempt to mint to blocked address.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Signer: minter.Address,
		To:     user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to blocked address.
	require.ErrorContains(t, err, "blocked from receiving")

	// ARRANGE: Unblock user account.
	k.DeleteBlockedAddress(ctx, user.Bytes)

	// ACT: Attempt to mint invalid amount.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Signer: minter.Address,
		To:     user.Address,
		Amount: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to invalid amount.
	require.ErrorContains(t, err, "amount must be positive")

	// ACT: Attempt to mint.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Signer: minter.Address,
		To:     user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, ONE, bank.Balances[user.Address].AmountOf(k.Denom))
	require.True(t, bank.Balances[types.ModuleName].IsZero())
	require.True(t, k.GetMinter(ctx, minter.Address).IsZero())

	// ACT: Attempt another mint with insufficient allowance.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Signer: minter.Address,
		To:     user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to insufficient allowance.
	require.ErrorContains(t, err, types.ErrInsufficientAllowance.Error())
}

func TestPause(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set pauser in state.
	pauser := utils.TestAccount()
	k.SetPauser(ctx, pauser.Address)

	// ACT: Attempt to pause with invalid signer.
	_, err := server.Pause(goCtx, &types.MsgPause{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidPauser.Error())
	require.False(t, k.GetPaused(ctx))

	// ACT: Attempt to pause.
	_, err = server.Pause(goCtx, &types.MsgPause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, k.GetPaused(ctx))

	// ACT: Attempt to pause again.
	_, err = server.Pause(goCtx, &types.MsgPause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've failed due to module being paused already.
	require.ErrorContains(t, err, "module is already paused")
	require.True(t, k.GetPaused(ctx))
}

func TestUnpause(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set paused state to true.
	k.SetPaused(ctx, true)

	// ACT: Attempt to unpause with no owner set.
	_, err := server.Unpause(goCtx, &types.MsgUnpause{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")
	require.True(t, k.GetPaused(ctx))

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to unpause with invalid signer.
	_, err = server.Unpause(goCtx, &types.MsgUnpause{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())
	require.True(t, k.GetPaused(ctx))

	// ACT: Attempt to unpause.
	_, err = server.Unpause(goCtx, &types.MsgUnpause{
		Signer: owner.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.False(t, k.GetPaused(ctx))

	// ACT: Attempt to unpause again.
	_, err = server.Unpause(goCtx, &types.MsgUnpause{
		Signer: owner.Address,
	})
	// ASSERT: The action should've failed due to module being unpaused already.
	require.ErrorContains(t, err, "module is already unpaused")
	require.False(t, k.GetPaused(ctx))
}

func TestTransferOwnership(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to transfer ownership with no owner set.
	_, err := server.TransferOwnership(goCtx, &types.MsgTransferOwnership{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to transfer ownership with invalid signer.
	_, err = server.TransferOwnership(goCtx, &types.MsgTransferOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ACT: Attempt to transfer ownership to same owner.
	_, err = server.TransferOwnership(goCtx, &types.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: owner.Address,
	})
	// ASSERT: The action should've failed due to same owner.
	require.ErrorContains(t, err, types.ErrSameOwner.Error())

	// ARRANGE: Generate a pending owner account.
	pendingOwner := utils.TestAccount()

	// ACT: Attempt to transfer ownership.
	_, err = server.TransferOwnership(goCtx, &types.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded, and set a pending owner in state.
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, k.GetPendingOwner(ctx))
}

func TestAcceptOwnership(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to accept ownership with no pending owner set.
	_, err := server.AcceptOwnership(goCtx, &types.MsgAcceptOwnership{})
	// ASSERT: The action should've failed due to no pending owner set.
	require.ErrorContains(t, err, "there is no pending owner")

	// ARRANGE: Set pending owner in state.
	pendingOwner := utils.TestAccount()
	k.SetPendingOwner(ctx, pendingOwner.Address)

	// ACT: Attempt to accept ownership with invalid signer.
	_, err = server.AcceptOwnership(goCtx, &types.MsgAcceptOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidPendingOwner.Error())

	// ACT: Attempt to accept ownership.
	_, err = server.AcceptOwnership(goCtx, &types.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded, and updated the owner in state.
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, k.GetOwner(ctx))
	require.Empty(t, k.GetPendingOwner(ctx))
}

func TestAddBurner(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add burner with no owner set.
	_, err := server.AddBurner(goCtx, &types.MsgAddBurner{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to add burner with invalid signer.
	_, err = server.AddBurner(goCtx, &types.MsgAddBurner{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate two burner accounts, add one to state.
	burner1, burner2 := utils.TestAccount(), utils.TestAccount()
	k.SetBurner(ctx, burner2.Address, ONE)

	// ACT: Attempt to add burner that already exists.
	_, err = server.AddBurner(goCtx, &types.MsgAddBurner{
		Signer:    owner.Address,
		Burner:    burner2.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to existing burner.
	require.ErrorContains(t, err, "is already a burner")

	// ACT: Attempt to add burner with invalid allowance.
	_, err = server.AddBurner(goCtx, &types.MsgAddBurner{
		Signer:    owner.Address,
		Burner:    burner1.Address,
		Allowance: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to negative allowance.
	require.ErrorContains(t, err, "allowance cannot be negative")

	// ACT: Attempt to add burner.
	_, err = server.AddBurner(goCtx, &types.MsgAddBurner{
		Signer:    owner.Address,
		Burner:    burner1.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set burner in state.
	require.NoError(t, err)
	require.Equal(t, ONE, k.GetBurner(ctx, burner1.Address))
}

func TestRemoveBurner(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove burner with no owner set.
	_, err := server.RemoveBurner(goCtx, &types.MsgRemoveBurner{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to remove burner with invalid signer.
	_, err = server.RemoveBurner(goCtx, &types.MsgRemoveBurner{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a burner account.
	burner := utils.TestAccount()

	// ACT: Attempt to remove burner that does not exist.
	_, err = server.RemoveBurner(goCtx, &types.MsgRemoveBurner{
		Signer: owner.Address,
		Burner: burner.Address,
	})
	// ASSERT: The action should've failed due to non existent burner.
	require.ErrorContains(t, err, "is not a burner")

	// ARRANGE: Set burner in state.
	k.SetBurner(ctx, burner.Address, ONE)

	// ACT: Attempt to remove burner.
	_, err = server.RemoveBurner(goCtx, &types.MsgRemoveBurner{
		Signer: owner.Address,
		Burner: burner.Address,
	})
	// ASSERT: The action should've succeeded, and removed burner in state.
	require.NoError(t, err)
	require.False(t, k.HasBurner(ctx, burner.Address))
}

func TestSetBurnerAllowance(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to set burner allowance with no owner set.
	_, err := server.SetBurnerAllowance(goCtx, &types.MsgSetBurnerAllowance{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to set burner allowance with invalid signer.
	_, err = server.SetBurnerAllowance(goCtx, &types.MsgSetBurnerAllowance{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a burner account.
	burner := utils.TestAccount()

	// ACT: Attempt to set burner allowance that does not exist.
	_, err = server.SetBurnerAllowance(goCtx, &types.MsgSetBurnerAllowance{
		Signer:    owner.Address,
		Burner:    burner.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to non existent burner.
	require.ErrorContains(t, err, "is not a burner")

	// ARRANGE: Set burner in state.
	k.SetBurner(ctx, burner.Address, sdk.ZeroInt())

	// ACT: Attempt to set burner allowance with invalid allowance.
	_, err = server.SetBurnerAllowance(goCtx, &types.MsgSetBurnerAllowance{
		Signer:    owner.Address,
		Burner:    burner.Address,
		Allowance: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to negative allowance.
	require.ErrorContains(t, err, "allowance cannot be negative")

	// ACT: Attempt to set burner allowance.
	_, err = server.SetBurnerAllowance(goCtx, &types.MsgSetBurnerAllowance{
		Signer:    owner.Address,
		Burner:    burner.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set burner allowance in state.
	require.NoError(t, err)
	require.Equal(t, ONE, k.GetBurner(ctx, burner.Address))
}

func TestAddMinter(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add minter with no owner set.
	_, err := server.AddMinter(goCtx, &types.MsgAddMinter{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to add minter with invalid signer.
	_, err = server.AddMinter(goCtx, &types.MsgAddMinter{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate two minter accounts, add one to state.
	minter1, minter2 := utils.TestAccount(), utils.TestAccount()
	k.SetMinter(ctx, minter2.Address, ONE)

	// ACT: Attempt to add minter that already exists.
	_, err = server.AddMinter(goCtx, &types.MsgAddMinter{
		Signer:    owner.Address,
		Minter:    minter2.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to existing minter.
	require.ErrorContains(t, err, "is already a minter")

	// ACT: Attempt to add minter with invalid allowance.
	_, err = server.AddMinter(goCtx, &types.MsgAddMinter{
		Signer:    owner.Address,
		Minter:    minter1.Address,
		Allowance: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to negative allowance.
	require.ErrorContains(t, err, "allowance cannot be negative")

	// ACT: Attempt to add minter.
	_, err = server.AddMinter(goCtx, &types.MsgAddMinter{
		Signer:    owner.Address,
		Minter:    minter1.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set minter in state.
	require.NoError(t, err)
	require.Equal(t, ONE, k.GetMinter(ctx, minter1.Address))
}

func TestRemoveMinter(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove minter with no owner set.
	_, err := server.RemoveMinter(goCtx, &types.MsgRemoveMinter{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to remove minter with invalid signer.
	_, err = server.RemoveMinter(goCtx, &types.MsgRemoveMinter{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a minter account.
	minter := utils.TestAccount()

	// ACT: Attempt to remove minter that does not exist.
	_, err = server.RemoveMinter(goCtx, &types.MsgRemoveMinter{
		Signer: owner.Address,
		Minter: minter.Address,
	})
	// ASSERT: The action should've failed due to non-existent minter.
	require.ErrorContains(t, err, "is not a minter")

	// ARRANGE: Set minter in state.
	k.SetMinter(ctx, minter.Address, ONE)

	// ACT: Attempt to remove minter.
	_, err = server.RemoveMinter(goCtx, &types.MsgRemoveMinter{
		Signer: owner.Address,
		Minter: minter.Address,
	})
	// ASSERT: The action should've succeeded, and removed minter in state.
	require.NoError(t, err)
	require.False(t, k.HasMinter(ctx, minter.Address))
}

func TestSetMinterAllowance(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to set minter allowance with no owner set.
	_, err := server.SetMinterAllowance(goCtx, &types.MsgSetMinterAllowance{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to set minter allowance with invalid signer.
	_, err = server.SetMinterAllowance(goCtx, &types.MsgSetMinterAllowance{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a minter account.
	minter := utils.TestAccount()

	// ACT: Attempt to set minter allowance that does not exist.
	_, err = server.SetMinterAllowance(goCtx, &types.MsgSetMinterAllowance{
		Signer:    owner.Address,
		Minter:    minter.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to non-existent minter.
	require.ErrorContains(t, err, "is not a minter")

	// ARRANGE: Set minters in state.
	k.SetMinter(ctx, minter.Address, sdk.ZeroInt())

	// ACT: Attempt to set minter allowance with invalid allowance.
	_, err = server.SetMinterAllowance(goCtx, &types.MsgSetMinterAllowance{
		Signer:    owner.Address,
		Minter:    minter.Address,
		Allowance: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to negative allowance.
	require.ErrorContains(t, err, "allowance cannot be negative")

	// ACT: Attempt to set minter allowance.
	_, err = server.SetMinterAllowance(goCtx, &types.MsgSetMinterAllowance{
		Signer:    owner.Address,
		Minter:    minter.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set minter allowance in state.
	require.NoError(t, err)
	require.Equal(t, ONE, k.GetMinter(ctx, minter.Address))
}

func TestAddPauser(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add pauser with no owner set.
	_, err := server.AddPauser(goCtx, &types.MsgAddPauser{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to add pauser with invalid signer.
	_, err = server.AddPauser(goCtx, &types.MsgAddPauser{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate two pauser accounts, add one to state.
	pauser1, pauser2 := utils.TestAccount(), utils.TestAccount()
	k.SetPauser(ctx, pauser2.Address)

	// ACT: Attempt to add pauser that already exists.
	_, err = server.AddPauser(goCtx, &types.MsgAddPauser{
		Signer: owner.Address,
		Pauser: pauser2.Address,
	})
	// ASSERT: The action should've failed due to existing pauser.
	require.ErrorContains(t, err, "is already a pauser")

	// ACT: Attempt to add pauser.
	_, err = server.AddPauser(goCtx, &types.MsgAddPauser{
		Signer: owner.Address,
		Pauser: pauser1.Address,
	})
	// ASSERT: The action should've succeeded, and set pauser in state.
	require.NoError(t, err)
	require.True(t, k.HasPauser(ctx, pauser1.Address))
}

func TestRemovePauser(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove pauser with no owner set.
	_, err := server.RemovePauser(goCtx, &types.MsgRemovePauser{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to remove pauser with invalid signer.
	_, err = server.RemovePauser(goCtx, &types.MsgRemovePauser{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a pauser account.
	pauser := utils.TestAccount()

	// ACT: Attempt to remove pauser that does not exist.
	_, err = server.RemovePauser(goCtx, &types.MsgRemovePauser{
		Signer: owner.Address,
		Pauser: pauser.Address,
	})
	// ASSERT: The action should've failed due to non-existent pauser.
	require.ErrorContains(t, err, "is not a pauser")

	// ARRANGE: Set pauser in state.
	k.SetPauser(ctx, pauser.Address)

	// ACT: Attempt to remove pauser.
	_, err = server.RemovePauser(goCtx, &types.MsgRemovePauser{
		Signer: owner.Address,
		Pauser: pauser.Address,
	})
	// ASSERT: The action should've succeeded, and removed pauser in state.
	require.NoError(t, err)
	require.False(t, k.HasPauser(ctx, pauser.Address))
}
