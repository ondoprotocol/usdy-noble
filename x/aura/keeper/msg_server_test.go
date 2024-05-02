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

	// ARRANGE: Set burner in state, with enough allowance for a single burn.
	burner := utils.TestAccount()
	require.NoError(t, k.Burners.Set(ctx, burner.Address, ONE))

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
		Amount: ONE,
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
	allowance, err := k.Burners.Get(ctx, burner.Address)
	require.NoError(t, err)
	require.True(t, allowance.IsZero())

	// ACT: Attempt another burn with insufficient allowance.
	_, err = server.Burn(ctx, &types.MsgBurn{
		Signer: burner.Address,
		From:   user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to insufficient allowance.
	require.ErrorContains(t, err, types.ErrInsufficientAllowance.Error())
}

func TestMint(t *testing.T) {
	bank := mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}
	k, ctx := mocks.AuraKeeperWithBank(t, bank)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set minter in state, with enough allowance for a single mint.
	minter := utils.TestAccount()
	require.NoError(t, k.Minters.Set(ctx, minter.Address, ONE))

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
		Amount: ONE,
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
	allowance, err := k.Minters.Get(ctx, minter.Address)
	require.NoError(t, err)
	require.True(t, allowance.IsZero())

	// ACT: Attempt another mint with insufficient allowance.
	_, err = server.Mint(ctx, &types.MsgMint{
		Signer: minter.Address,
		To:     user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to insufficient allowance.
	require.ErrorContains(t, err, types.ErrInsufficientAllowance.Error())
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

	// ACT: Attempt to unpause with no owner set.
	_, err := server.Unpause(ctx, &types.MsgUnpause{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "unable to retrieve owner from state")
	paused, _ := k.Paused.Get(ctx)
	require.True(t, paused)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.Owner.Set(ctx, owner.Address))

	// ACT: Attempt to unpause with invalid signer.
	_, err = server.Unpause(ctx, &types.MsgUnpause{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())
	paused, _ = k.Paused.Get(ctx)
	require.True(t, paused)

	// ACT: Attempt to unpause.
	_, err = server.Unpause(ctx, &types.MsgUnpause{
		Signer: owner.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	paused, _ = k.Paused.Get(ctx)
	require.False(t, paused)

	// ACT: Attempt to unpause again.
	_, err = server.Unpause(ctx, &types.MsgUnpause{
		Signer: owner.Address,
	})
	// ASSERT: The action should've failed due to module being unpaused already.
	require.ErrorContains(t, err, "module is already unpaused")
	paused, _ = k.Paused.Get(ctx)
	require.False(t, paused)
}

func TestTransferOwnership(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to transfer ownership with no owner set.
	_, err := server.TransferOwnership(ctx, &types.MsgTransferOwnership{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "unable to retrieve owner from state")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.Owner.Set(ctx, owner.Address))

	// ACT: Attempt to transfer ownership with invalid signer.
	_, err = server.TransferOwnership(ctx, &types.MsgTransferOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a pending owner account.
	pendingOwner := utils.TestAccount()

	// ACT: Attempt to transfer ownership.
	_, err = server.TransferOwnership(ctx, &types.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded, and set a pending owner in state.
	require.NoError(t, err)
	res, err := k.PendingOwner.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, res)
}

func TestAcceptOwnership(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to accept ownership with no pending owner set.
	_, err := server.AcceptOwnership(ctx, &types.MsgAcceptOwnership{})
	// ASSERT: The action should've failed due to no pending owner set.
	require.ErrorContains(t, err, "there is no pending owner")

	// ARRANGE: Set pending owner in state.
	pendingOwner := utils.TestAccount()
	require.NoError(t, k.PendingOwner.Set(ctx, pendingOwner.Address))

	// ACT: Attempt to accept ownership with invalid signer.
	_, err = server.AcceptOwnership(ctx, &types.MsgAcceptOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidPendingOwner.Error())

	// ACT: Attempt to accept ownership.
	_, err = server.AcceptOwnership(ctx, &types.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded, and updated the owner in state.
	require.NoError(t, err)
	res, err := k.Owner.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, res)
	has, err := k.PendingOwner.Has(ctx)
	require.NoError(t, err)
	require.False(t, has)
}

func TestAddBurner(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add burner with no owner set.
	_, err := server.AddBurner(ctx, &types.MsgAddBurner{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "unable to retrieve owner from state")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.Owner.Set(ctx, owner.Address))

	// ACT: Attempt to add burner with invalid signer.
	_, err = server.AddBurner(ctx, &types.MsgAddBurner{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate two burner accounts, add one to state.
	burner1, burner2 := utils.TestAccount(), utils.TestAccount()
	require.NoError(t, k.Burners.Set(ctx, burner2.Address, ONE))

	// ACT: Attempt to add burner that already exists.
	_, err = server.AddBurner(ctx, &types.MsgAddBurner{
		Signer:    owner.Address,
		Burner:    burner2.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to existing burner.
	require.ErrorContains(t, err, "is already a burner")

	// ACT: Attempt to add burner.
	_, err = server.AddBurner(ctx, &types.MsgAddBurner{
		Signer:    owner.Address,
		Burner:    burner1.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set burner in state.
	require.NoError(t, err)
	allowance, err := k.Burners.Get(ctx, burner1.Address)
	require.NoError(t, err)
	require.Equal(t, ONE, allowance)
}

func TestRemoveBurner(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove burner with no owner set.
	_, err := server.RemoveBurner(ctx, &types.MsgRemoveBurner{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "unable to retrieve owner from state")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.Owner.Set(ctx, owner.Address))

	// ACT: Attempt to remove burner with invalid signer.
	_, err = server.RemoveBurner(ctx, &types.MsgRemoveBurner{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a burner account.
	burner := utils.TestAccount()

	// ACT: Attempt to remove burner that does not exist.
	_, err = server.RemoveBurner(ctx, &types.MsgRemoveBurner{
		Signer: owner.Address,
		Burner: burner.Address,
	})
	// ASSERT: The action should've failed due to non existent burner.
	require.ErrorContains(t, err, "is not a burner")

	// ARRANGE: Set burner in state.
	require.NoError(t, k.Burners.Set(ctx, burner.Address, ONE))

	// ACT: Attempt to remove burner.
	_, err = server.RemoveBurner(ctx, &types.MsgRemoveBurner{
		Signer: owner.Address,
		Burner: burner.Address,
	})
	// ASSERT: The action should've succeeded, and removed burner in state.
	require.NoError(t, err)
	has, err := k.Burners.Has(ctx, burner.Address)
	require.NoError(t, err)
	require.False(t, has)
}

func TestSetBurnerAllowance(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to set burner allowance with no owner set.
	_, err := server.SetBurnerAllowance(ctx, &types.MsgSetBurnerAllowance{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "unable to retrieve owner from state")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.Owner.Set(ctx, owner.Address))

	// ACT: Attempt to set burner allowance with invalid signer.
	_, err = server.SetBurnerAllowance(ctx, &types.MsgSetBurnerAllowance{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a burner account.
	burner := utils.TestAccount()

	// ACT: Attempt to set burner allowance that does not exist.
	_, err = server.SetBurnerAllowance(ctx, &types.MsgSetBurnerAllowance{
		Signer:    owner.Address,
		Burner:    burner.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to non existent burner.
	require.ErrorContains(t, err, "is not a burner")

	// ARRANGE: Set burner in state.
	require.NoError(t, k.Burners.Set(ctx, burner.Address, math.ZeroInt()))

	// ACT: Attempt to set burner allowance.
	_, err = server.SetBurnerAllowance(ctx, &types.MsgSetBurnerAllowance{
		Signer:    owner.Address,
		Burner:    burner.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set burner allowance in state.
	require.NoError(t, err)
	allowance, err := k.Burners.Get(ctx, burner.Address)
	require.NoError(t, err)
	require.Equal(t, ONE, allowance)
}

func TestAddMinter(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add minter with no owner set.
	_, err := server.AddMinter(ctx, &types.MsgAddMinter{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "unable to retrieve owner from state")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.Owner.Set(ctx, owner.Address))

	// ACT: Attempt to add minter with invalid signer.
	_, err = server.AddMinter(ctx, &types.MsgAddMinter{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate two minter accounts, add one to state.
	minter1, minter2 := utils.TestAccount(), utils.TestAccount()
	require.NoError(t, k.Minters.Set(ctx, minter2.Address, ONE))

	// ACT: Attempt to add minter that already exists.
	_, err = server.AddMinter(ctx, &types.MsgAddMinter{
		Signer:    owner.Address,
		Minter:    minter2.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to existing minter.
	require.ErrorContains(t, err, "is already a minter")

	// ACT: Attempt to add minter.
	_, err = server.AddMinter(ctx, &types.MsgAddMinter{
		Signer:    owner.Address,
		Minter:    minter1.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set minter in state.
	require.NoError(t, err)
	allowance, err := k.Minters.Get(ctx, minter1.Address)
	require.NoError(t, err)
	require.Equal(t, ONE, allowance)
}

func TestRemoveMinter(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove minter with no owner set.
	_, err := server.RemoveMinter(ctx, &types.MsgRemoveMinter{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "unable to retrieve owner from state")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.Owner.Set(ctx, owner.Address))

	// ACT: Attempt to remove minter with invalid signer.
	_, err = server.RemoveMinter(ctx, &types.MsgRemoveMinter{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a minter account.
	minter := utils.TestAccount()

	// ACT: Attempt to remove minter that does not exist.
	_, err = server.RemoveMinter(ctx, &types.MsgRemoveMinter{
		Signer: owner.Address,
		Minter: minter.Address,
	})
	// ASSERT: The action should've failed due to non-existent minter.
	require.ErrorContains(t, err, "is not a minter")

	// ARRANGE: Set minter in state.
	require.NoError(t, k.Minters.Set(ctx, minter.Address, ONE))

	// ACT: Attempt to remove minter.
	_, err = server.RemoveMinter(ctx, &types.MsgRemoveMinter{
		Signer: owner.Address,
		Minter: minter.Address,
	})
	// ASSERT: The action should've succeeded, and removed minter in state.
	require.NoError(t, err)
	has, err := k.Minters.Has(ctx, minter.Address)
	require.NoError(t, err)
	require.False(t, has)
}

func TestSetMinterAllowance(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to set minter allowance with no owner set.
	_, err := server.SetMinterAllowance(ctx, &types.MsgSetMinterAllowance{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "unable to retrieve owner from state")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.Owner.Set(ctx, owner.Address))

	// ACT: Attempt to set minter allowance with invalid signer.
	_, err = server.SetMinterAllowance(ctx, &types.MsgSetMinterAllowance{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a minter account.
	minter := utils.TestAccount()

	// ACT: Attempt to set minter allowance that does not exist.
	_, err = server.SetMinterAllowance(ctx, &types.MsgSetMinterAllowance{
		Signer:    owner.Address,
		Minter:    minter.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to non-existent minter.
	require.ErrorContains(t, err, "is not a minter")

	// ARRANGE: Set minters in state.
	require.NoError(t, k.Minters.Set(ctx, minter.Address, math.ZeroInt()))

	// ACT: Attempt to set minter allowance.
	_, err = server.SetMinterAllowance(ctx, &types.MsgSetMinterAllowance{
		Signer:    owner.Address,
		Minter:    minter.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set minter allowance in state.
	require.NoError(t, err)
	allowance, err := k.Minters.Get(ctx, minter.Address)
	require.NoError(t, err)
	require.Equal(t, ONE, allowance)
}

func TestAddPauser(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add pauser with no owner set.
	_, err := server.AddPauser(ctx, &types.MsgAddPauser{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "unable to retrieve owner from state")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.Owner.Set(ctx, owner.Address))

	// ACT: Attempt to add pauser with invalid signer.
	_, err = server.AddPauser(ctx, &types.MsgAddPauser{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate two pauser accounts, add one to state.
	pauser1, pauser2 := utils.TestAccount(), utils.TestAccount()
	require.NoError(t, k.Pausers.Set(ctx, pauser2.Address))

	// ACT: Attempt to add pauser that already exists.
	_, err = server.AddPauser(ctx, &types.MsgAddPauser{
		Signer: owner.Address,
		Pauser: pauser2.Address,
	})
	// ASSERT: The action should've failed due to existing pauser.
	require.ErrorContains(t, err, "is already a pauser")

	// ACT: Attempt to add pauser.
	_, err = server.AddPauser(ctx, &types.MsgAddPauser{
		Signer: owner.Address,
		Pauser: pauser1.Address,
	})
	// ASSERT: The action should've succeeded, and set pauser in state.
	require.NoError(t, err)
	has, err := k.Pausers.Has(ctx, pauser1.Address)
	require.NoError(t, err)
	require.True(t, has)
}

func TestRemovePauser(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove pauser with no owner set.
	_, err := server.RemovePauser(ctx, &types.MsgRemovePauser{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "unable to retrieve owner from state")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.Owner.Set(ctx, owner.Address))

	// ACT: Attempt to remove pauser with invalid signer.
	_, err = server.RemovePauser(ctx, &types.MsgRemovePauser{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a pauser account.
	pauser := utils.TestAccount()

	// ACT: Attempt to remove pauser that does not exist.
	_, err = server.RemovePauser(ctx, &types.MsgRemovePauser{
		Signer: owner.Address,
		Pauser: pauser.Address,
	})
	// ASSERT: The action should've failed due to non-existent pauser.
	require.ErrorContains(t, err, "is not a pauser")

	// ARRANGE: Set pauser in state.
	require.NoError(t, k.Pausers.Set(ctx, pauser.Address))

	// ACT: Attempt to remove pauser.
	_, err = server.RemovePauser(ctx, &types.MsgRemovePauser{
		Signer: owner.Address,
		Pauser: pauser.Address,
	})
	// ASSERT: The action should've succeeded, and removed pauser in state.
	require.NoError(t, err)
	has, err := k.Pausers.Has(ctx, pauser.Address)
	require.NoError(t, err)
	require.False(t, has)
}
