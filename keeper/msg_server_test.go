package keeper_test

import (
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/usdy-noble/v2/keeper"
	"github.com/ondoprotocol/usdy-noble/v2/types"
	"github.com/ondoprotocol/usdy-noble/v2/utils"
	"github.com/ondoprotocol/usdy-noble/v2/utils/mocks"
	"github.com/stretchr/testify/require"
)

var ONE = math.NewInt(1_000_000_000_000_000_000)

func TestBurn(t *testing.T) {
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.AuraKeeperWithBank(bank)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set burner in state, with enough allowance for a single burn.
	burner := utils.TestAccount()
	require.NoError(t, k.SetBurner(ctx, burner.Address, ONE))

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

	// ACT: Attempt to burn invalid amount.
	_, err = server.Burn(ctx, &types.MsgBurn{
		Signer: burner.Address,
		From:   user.Address,
		Amount: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to invalid amount.
	require.ErrorContains(t, err, "amount must be positive")

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

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.Burners
	k.Burners = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.BurnerPrefix, "burners", collections.StringKey, sdk.IntValue,
	)

	// ACT: Attempt to burn with failing Burners collection store.
	_, err = server.Burn(ctx, &types.MsgBurn{
		Signer: burner.Address,
		From:   user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Burners = tmp

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
	require.True(t, k.GetBurner(ctx, burner.Address).IsZero())

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
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.AuraKeeperWithBank(bank)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set minter in state, with enough allowance for a single mint.
	minter := utils.TestAccount()
	require.NoError(t, k.SetMinter(ctx, minter.Address, ONE))

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

	// ARRANGE: Generate a user account and add to blocklist.
	user := utils.TestAccount()
	require.NoError(t, k.SetBlockedAddress(ctx, user.Bytes))

	// ACT: Attempt to mint to blocked address.
	_, err = server.Mint(ctx, &types.MsgMint{
		Signer: minter.Address,
		To:     user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to blocked address.
	require.ErrorContains(t, err, "blocked from receiving")

	// ARRANGE: Unblock user account.
	require.NoError(t, k.DeleteBlockedAddress(ctx, user.Bytes))

	// ACT: Attempt to mint invalid amount.
	_, err = server.Mint(ctx, &types.MsgMint{
		Signer: minter.Address,
		To:     user.Address,
		Amount: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to invalid amount.
	require.ErrorContains(t, err, "amount must be positive")

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.Minters
	k.Minters = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.MinterPrefix, "minters", collections.StringKey, sdk.IntValue,
	)

	// ACT: Attempt to mint with failing Minters collection store.
	_, err = server.Mint(ctx, &types.MsgMint{
		Signer: minter.Address,
		To:     user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Minters = tmp

	// ACT: Attempt to mint.
	_, err = server.Mint(ctx, &types.MsgMint{
		Signer: minter.Address,
		To:     user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, ONE.MulRaw(2), bank.Balances[user.Address].AmountOf(k.Denom))
	require.True(t, bank.Balances[types.ModuleName].IsZero())
	require.True(t, k.GetMinter(ctx, minter.Address).IsZero())

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
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set pauser in state.
	pauser := utils.TestAccount()
	require.NoError(t, k.SetPauser(ctx, pauser.Address))

	// ACT: Attempt to pause with invalid signer.
	_, err := server.Pause(ctx, &types.MsgPause{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidPauser.Error())
	require.False(t, k.GetPaused(ctx))

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.Paused
	k.Paused = collections.NewItem(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.PausedKey, "paused", collections.BoolValue,
	)

	// ACT: Attempt to pause with failing Paused collection store.
	_, err = server.Pause(ctx, &types.MsgPause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Paused = tmp

	// ACT: Attempt to pause.
	_, err = server.Pause(ctx, &types.MsgPause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, k.GetPaused(ctx))

	// ACT: Attempt to pause again.
	_, err = server.Pause(ctx, &types.MsgPause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've failed due to module being paused already.
	require.ErrorContains(t, err, "module is already paused")
	require.True(t, k.GetPaused(ctx))
}

func TestUnpause(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set paused state to true.
	require.NoError(t, k.SetPaused(ctx, true))

	// ACT: Attempt to unpause with no owner set.
	_, err := server.Unpause(ctx, &types.MsgUnpause{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")
	require.True(t, k.GetPaused(ctx))

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetOwner(ctx, owner.Address))

	// ACT: Attempt to unpause with invalid signer.
	_, err = server.Unpause(ctx, &types.MsgUnpause{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())
	require.True(t, k.GetPaused(ctx))

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.Paused
	k.Paused = collections.NewItem(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.PausedKey, "paused", collections.BoolValue,
	)

	// ACT: Attempt to unpause with failing Paused collection store.
	_, err = server.Unpause(ctx, &types.MsgUnpause{
		Signer: owner.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Paused = tmp

	// ACT: Attempt to unpause.
	_, err = server.Unpause(ctx, &types.MsgUnpause{
		Signer: owner.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.False(t, k.GetPaused(ctx))

	// ACT: Attempt to unpause again.
	_, err = server.Unpause(ctx, &types.MsgUnpause{
		Signer: owner.Address,
	})
	// ASSERT: The action should've failed due to module being unpaused already.
	require.ErrorContains(t, err, "module is already unpaused")
	require.False(t, k.GetPaused(ctx))
}

func TestTransferOwnership(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to transfer ownership with no owner set.
	_, err := server.TransferOwnership(ctx, &types.MsgTransferOwnership{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetOwner(ctx, owner.Address))

	// ACT: Attempt to transfer ownership with invalid signer.
	_, err = server.TransferOwnership(ctx, &types.MsgTransferOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ACT: Attempt to transfer ownership to same owner.
	_, err = server.TransferOwnership(ctx, &types.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: owner.Address,
	})
	// ASSERT: The action should've failed due to same owner.
	require.ErrorContains(t, err, types.ErrSameOwner.Error())

	// ARRANGE: Generate a pending owner account.
	pendingOwner := utils.TestAccount()

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.PendingOwner
	k.PendingOwner = collections.NewItem(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.PendingOwnerKey, "pending_owner", collections.StringValue,
	)

	// ACT: Attempt to transfer ownership with failing PendingOwner collection store.
	_, err = server.TransferOwnership(ctx, &types.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: pendingOwner.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.PendingOwner = tmp

	// ACT: Attempt to transfer ownership.
	_, err = server.TransferOwnership(ctx, &types.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded, and set a pending owner in state.
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, k.GetPendingOwner(ctx))
}

func TestAcceptOwnership(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to accept ownership with no pending owner set.
	_, err := server.AcceptOwnership(ctx, &types.MsgAcceptOwnership{})
	// ASSERT: The action should've failed due to no pending owner set.
	require.ErrorContains(t, err, "there is no pending owner")

	// ARRANGE: Set pending owner in state.
	pendingOwner := utils.TestAccount()
	require.NoError(t, k.SetPendingOwner(ctx, pendingOwner.Address))

	// ACT: Attempt to accept ownership with invalid signer.
	_, err = server.AcceptOwnership(ctx, &types.MsgAcceptOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidPendingOwner.Error())

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.Owner
	k.Owner = collections.NewItem(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.OwnerKey, "owner", collections.StringValue,
	)

	// ACT: Attempt to accept ownership with failing Owner collection store.
	_, err = server.AcceptOwnership(ctx, &types.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Owner = tmp

	// ARRANGE: Set up a failing collection store for the attribute delete.
	tmp = k.PendingOwner
	k.PendingOwner = collections.NewItem(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Delete, utils.GetKVStore(ctx, types.ModuleName))),
		types.PendingOwnerKey, "pending_owner", collections.StringValue,
	)

	// ACT: Attempt to accept ownership with failing PendingOwner collection store.
	_, err = server.AcceptOwnership(ctx, &types.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've failed due to collection store delete error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.PendingOwner = tmp

	// ACT: Attempt to accept ownership.
	_, err = server.AcceptOwnership(ctx, &types.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded, and updated the owner in state.
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, k.GetOwner(ctx))
	require.Empty(t, k.GetPendingOwner(ctx))
}

func TestAddBurner(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add burner with no owner set.
	_, err := server.AddBurner(ctx, &types.MsgAddBurner{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetOwner(ctx, owner.Address))

	// ACT: Attempt to add burner with invalid signer.
	_, err = server.AddBurner(ctx, &types.MsgAddBurner{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate two burner accounts, add one to state.
	burner1, burner2 := utils.TestAccount(), utils.TestAccount()
	require.NoError(t, k.SetBurner(ctx, burner2.Address, ONE))

	// ACT: Attempt to add burner that already exists.
	_, err = server.AddBurner(ctx, &types.MsgAddBurner{
		Signer:    owner.Address,
		Burner:    burner2.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to existing burner.
	require.ErrorContains(t, err, "is already a burner")

	// ACT: Attempt to add burner with invalid allowance.
	_, err = server.AddBurner(ctx, &types.MsgAddBurner{
		Signer:    owner.Address,
		Burner:    burner1.Address,
		Allowance: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to negative allowance.
	require.ErrorContains(t, err, "allowance cannot be negative")

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.Burners
	k.Burners = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.BurnerPrefix, "burners", collections.StringKey, sdk.IntValue,
	)

	// ACT: Attempt to add burner with failing Burners collection store.
	_, err = server.AddBurner(ctx, &types.MsgAddBurner{
		Signer:    owner.Address,
		Burner:    burner1.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Burners = tmp

	// ACT: Attempt to add burner.
	_, err = server.AddBurner(ctx, &types.MsgAddBurner{
		Signer:    owner.Address,
		Burner:    burner1.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set burner in state.
	require.NoError(t, err)
	require.Equal(t, ONE, k.GetBurner(ctx, burner1.Address))
}

func TestRemoveBurner(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove burner with no owner set.
	_, err := server.RemoveBurner(ctx, &types.MsgRemoveBurner{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetOwner(ctx, owner.Address))

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
	require.NoError(t, k.SetBurner(ctx, burner.Address, ONE))

	// ARRANGE: Set up a failing collection store for the attribute delete.
	tmp := k.Burners
	k.Burners = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Delete, utils.GetKVStore(ctx, types.ModuleName))),
		types.BurnerPrefix, "burners", collections.StringKey, sdk.IntValue,
	)

	// ACT: Attempt to remove burner with failing Burners collection store.
	_, err = server.RemoveBurner(ctx, &types.MsgRemoveBurner{
		Signer: owner.Address,
		Burner: burner.Address,
	})
	// ASSERT: The action should've failed due to collection store delete error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Burners = tmp

	// ACT: Attempt to remove burner.
	_, err = server.RemoveBurner(ctx, &types.MsgRemoveBurner{
		Signer: owner.Address,
		Burner: burner.Address,
	})
	// ASSERT: The action should've succeeded, and removed burner in state.
	require.NoError(t, err)
	require.False(t, k.HasBurner(ctx, burner.Address))
}

func TestSetBurnerAllowance(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to set burner allowance with no owner set.
	_, err := server.SetBurnerAllowance(ctx, &types.MsgSetBurnerAllowance{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetOwner(ctx, owner.Address))

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
	require.NoError(t, k.SetBurner(ctx, burner.Address, math.ZeroInt()))

	// ACT: Attempt to set burner allowance with invalid allowance.
	_, err = server.SetBurnerAllowance(ctx, &types.MsgSetBurnerAllowance{
		Signer:    owner.Address,
		Burner:    burner.Address,
		Allowance: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to negative allowance.
	require.ErrorContains(t, err, "allowance cannot be negative")

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.Burners
	k.Burners = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.BurnerPrefix, "burners", collections.StringKey, sdk.IntValue,
	)

	// ACT: Attempt to set burner allowance with failing Burners collection store.
	_, err = server.SetBurnerAllowance(ctx, &types.MsgSetBurnerAllowance{
		Signer:    owner.Address,
		Burner:    burner.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Burners = tmp

	// ARRANGE: Set invalid un-decodable burn allowance in state.
	key, _ := collections.EncodeKeyWithPrefix(types.BurnerPrefix, collections.StringKey, burner.Address)
	utils.GetKVStore(ctx, types.ModuleName).Set(key, []byte("invalid"))

	// ACT: Attempt to set burner allowance.
	_, err = server.SetBurnerAllowance(ctx, &types.MsgSetBurnerAllowance{
		Signer:    owner.Address,
		Burner:    burner.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set burner allowance in state.
	require.NoError(t, err)
	require.Equal(t, ONE, k.GetBurner(ctx, burner.Address))
}

func TestAddMinter(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add minter with no owner set.
	_, err := server.AddMinter(ctx, &types.MsgAddMinter{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetOwner(ctx, owner.Address))

	// ACT: Attempt to add minter with invalid signer.
	_, err = server.AddMinter(ctx, &types.MsgAddMinter{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate two minter accounts, add one to state.
	minter1, minter2 := utils.TestAccount(), utils.TestAccount()
	require.NoError(t, k.SetMinter(ctx, minter2.Address, ONE))

	// ACT: Attempt to add minter that already exists.
	_, err = server.AddMinter(ctx, &types.MsgAddMinter{
		Signer:    owner.Address,
		Minter:    minter2.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to existing minter.
	require.ErrorContains(t, err, "is already a minter")

	// ACT: Attempt to add minter with invalid allowance.
	_, err = server.AddMinter(ctx, &types.MsgAddMinter{
		Signer:    owner.Address,
		Minter:    minter1.Address,
		Allowance: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to negative allowance.
	require.ErrorContains(t, err, "allowance cannot be negative")

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.Minters
	k.Minters = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.MinterPrefix, "minters", collections.StringKey, sdk.IntValue,
	)

	// ACT: Attempt to add minter with failing Minters collection store.
	_, err = server.AddMinter(ctx, &types.MsgAddMinter{
		Signer:    owner.Address,
		Minter:    minter1.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Minters = tmp

	// ACT: Attempt to add minter.
	_, err = server.AddMinter(ctx, &types.MsgAddMinter{
		Signer:    owner.Address,
		Minter:    minter1.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set minter in state.
	require.NoError(t, err)
	require.Equal(t, ONE, k.GetMinter(ctx, minter1.Address))
}

func TestRemoveMinter(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove minter with no owner set.
	_, err := server.RemoveMinter(ctx, &types.MsgRemoveMinter{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetOwner(ctx, owner.Address))

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
	require.NoError(t, k.SetMinter(ctx, minter.Address, ONE))

	// ARRANGE: Set up a failing collection store for the attribute delete.
	tmp := k.Minters
	k.Minters = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Delete, utils.GetKVStore(ctx, types.ModuleName))),
		types.MinterPrefix, "minters", collections.StringKey, sdk.IntValue,
	)

	// ACT: Attempt to remove minter with failing Minters collection store.
	_, err = server.RemoveMinter(ctx, &types.MsgRemoveMinter{
		Signer: owner.Address,
		Minter: minter.Address,
	})
	// ASSERT: The action should've failed due to collection store delete error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Minters = tmp

	// ACT: Attempt to remove minter.
	_, err = server.RemoveMinter(ctx, &types.MsgRemoveMinter{
		Signer: owner.Address,
		Minter: minter.Address,
	})
	// ASSERT: The action should've succeeded, and removed minter in state.
	require.NoError(t, err)
	require.False(t, k.HasMinter(ctx, minter.Address))
}

func TestSetMinterAllowance(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to set minter allowance with no owner set.
	_, err := server.SetMinterAllowance(ctx, &types.MsgSetMinterAllowance{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetOwner(ctx, owner.Address))

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
	require.NoError(t, k.SetMinter(ctx, minter.Address, math.ZeroInt()))

	// ACT: Attempt to set minter allowance with invalid allowance.
	_, err = server.SetMinterAllowance(ctx, &types.MsgSetMinterAllowance{
		Signer:    owner.Address,
		Minter:    minter.Address,
		Allowance: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to negative allowance.
	require.ErrorContains(t, err, "allowance cannot be negative")

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.Minters
	k.Minters = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.MinterPrefix, "minters", collections.StringKey, sdk.IntValue,
	)

	// ACT: Attempt to set minter allowance with failing Minters collection store.
	_, err = server.SetMinterAllowance(ctx, &types.MsgSetMinterAllowance{
		Signer:    owner.Address,
		Minter:    minter.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Minters = tmp

	// ARRANGE: Set invalid un-decodable mint allowance in state.
	key, _ := collections.EncodeKeyWithPrefix(types.MinterPrefix, collections.StringKey, minter.Address)
	utils.GetKVStore(ctx, types.ModuleName).Set(key, []byte("invalid"))

	// ACT: Attempt to set minter allowance.
	_, err = server.SetMinterAllowance(ctx, &types.MsgSetMinterAllowance{
		Signer:    owner.Address,
		Minter:    minter.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set minter allowance in state.
	require.NoError(t, err)
	require.Equal(t, ONE, k.GetMinter(ctx, minter.Address))
}

func TestAddPauser(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add pauser with no owner set.
	_, err := server.AddPauser(ctx, &types.MsgAddPauser{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetOwner(ctx, owner.Address))

	// ACT: Attempt to add pauser with invalid signer.
	_, err = server.AddPauser(ctx, &types.MsgAddPauser{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate two pauser accounts, add one to state.
	pauser1, pauser2 := utils.TestAccount(), utils.TestAccount()
	require.NoError(t, k.SetPauser(ctx, pauser2.Address))

	// ACT: Attempt to add pauser that already exists.
	_, err = server.AddPauser(ctx, &types.MsgAddPauser{
		Signer: owner.Address,
		Pauser: pauser2.Address,
	})
	// ASSERT: The action should've failed due to existing pauser.
	require.ErrorContains(t, err, "is already a pauser")

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.Pausers
	k.Pausers = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.PauserPrefix, "pausers", collections.StringKey, collections.BytesValue,
	)

	// ACT: Attempt to add pauser with failing Pausers collection store.
	_, err = server.AddPauser(ctx, &types.MsgAddPauser{
		Signer: owner.Address,
		Pauser: pauser1.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Pausers = tmp

	// ACT: Attempt to add pauser.
	_, err = server.AddPauser(ctx, &types.MsgAddPauser{
		Signer: owner.Address,
		Pauser: pauser1.Address,
	})
	// ASSERT: The action should've succeeded, and set pauser in state.
	require.NoError(t, err)
	require.True(t, k.HasPauser(ctx, pauser1.Address))
}

func TestRemovePauser(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove pauser with no owner set.
	_, err := server.RemovePauser(ctx, &types.MsgRemovePauser{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetOwner(ctx, owner.Address))

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
	require.NoError(t, k.SetPauser(ctx, pauser.Address))

	// ARRANGE: Set up a failing collection store for the attribute delete.
	tmp := k.Pausers
	k.Pausers = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Delete, utils.GetKVStore(ctx, types.ModuleName))),
		types.PauserPrefix, "pausers", collections.StringKey, collections.BytesValue,
	)

	// ACT: Attempt to remove pauser with failing Pausers collection store.
	_, err = server.RemovePauser(ctx, &types.MsgRemovePauser{
		Signer: owner.Address,
		Pauser: pauser.Address,
	})
	// ASSERT: The action should've failed due to collection store delete error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Pausers = tmp

	// ACT: Attempt to remove pauser.
	_, err = server.RemovePauser(ctx, &types.MsgRemovePauser{
		Signer: owner.Address,
		Pauser: pauser.Address,
	})
	// ASSERT: The action should've succeeded, and removed pauser in state.
	require.NoError(t, err)
	require.False(t, k.HasPauser(ctx, pauser.Address))
}

func TestAddBlockedChannel(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add blocked channel with no owner set.
	_, err := server.AddBlockedChannel(ctx, &types.MsgAddBlockedChannel{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetOwner(ctx, owner.Address))

	// ACT: Attempt to add blocked channel with invalid signer.
	_, err = server.AddBlockedChannel(ctx, &types.MsgAddBlockedChannel{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate two channel, add one to state.
	channel1, channel2 := "channel-0", "channel-1"
	require.NoError(t, k.SetBlockedChannel(ctx, channel2))

	// ACT: Attempt to add blocked channel that is blocked.
	_, err = server.AddBlockedChannel(ctx, &types.MsgAddBlockedChannel{
		Signer:  owner.Address,
		Channel: channel2,
	})
	// ASSERT: The action should've failed due to blocked channel.
	require.ErrorContains(t, err, "is already blocked")

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.BlockedChannels
	k.BlockedChannels = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.BlockedChannelPrefix, "blocked_channels", collections.StringKey, collections.BytesValue,
	)

	// ACT: Attempt to add blocked channel with failing BlockedChannels collection store.
	_, err = server.AddBlockedChannel(ctx, &types.MsgAddBlockedChannel{
		Signer:  owner.Address,
		Channel: channel1,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.BlockedChannels = tmp

	// ACT: Attempt to add blocked channel.
	_, err = server.AddBlockedChannel(ctx, &types.MsgAddBlockedChannel{
		Signer:  owner.Address,
		Channel: channel1,
	})
	// ASSERT: The action should've succeeded, and set channel in state.
	require.NoError(t, err)
	require.True(t, k.HasBlockedChannel(ctx, channel1))
}

func TestRemoveBlockedChannel(t *testing.T) {
	k, ctx := mocks.AuraKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove blocked channel with no owner set.
	_, err := server.RemoveBlockedChannel(ctx, &types.MsgRemoveBlockedChannel{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.SetOwner(ctx, owner.Address))

	// ACT: Attempt to remove blocked channel with invalid signer.
	_, err = server.RemoveBlockedChannel(ctx, &types.MsgRemoveBlockedChannel{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a channel.
	channel := "channel-0"

	// ACT: Attempt to remove blocked channel that isn't blocked.
	_, err = server.RemoveBlockedChannel(ctx, &types.MsgRemoveBlockedChannel{
		Signer:  owner.Address,
		Channel: channel,
	})
	// ASSERT: The action should've failed due to allowed channel.
	require.ErrorContains(t, err, "is not blocked")

	// ARRANGE: Set channel in state.
	require.NoError(t, k.SetBlockedChannel(ctx, channel))

	// ARRANGE: Set up a failing collection store for the attribute delete.
	tmp := k.BlockedChannels
	k.BlockedChannels = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Delete, utils.GetKVStore(ctx, types.ModuleName))),
		types.BlockedChannelPrefix, "blocked_channels", collections.StringKey, collections.BytesValue,
	)

	// ACT: Attempt to remove blocked channel with failing BlockedChannels collection store.
	_, err = server.RemoveBlockedChannel(ctx, &types.MsgRemoveBlockedChannel{
		Signer:  owner.Address,
		Channel: channel,
	})
	// ASSERT: The action should've failed due to collection store delete error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.BlockedChannels = tmp

	// ACT: Attempt to remove blocked channel.
	_, err = server.RemoveBlockedChannel(ctx, &types.MsgRemoveBlockedChannel{
		Signer:  owner.Address,
		Channel: channel,
	})
	// ASSERT: The action should've succeeded, and removed channel in state.
	require.NoError(t, err)
	require.False(t, k.HasBlockedChannel(ctx, channel))
}
