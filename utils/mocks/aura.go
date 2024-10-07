package mocks

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/usdy-noble/v2/keeper"
	"github.com/ondoprotocol/usdy-noble/v2/types"
)

func AuraKeeper() (*keeper.Keeper, sdk.Context) {
	return AuraKeeperWithBank(BankKeeper{
		Restriction: NoOpSendRestrictionFn,
	})
}

func AuraKeeperWithBank(bank BankKeeper) (*keeper.Keeper, sdk.Context) {
	key := storetypes.NewKVStoreKey(types.ModuleName)
	tkey := storetypes.NewTransientStoreKey("transient_aura")

	k := keeper.NewKeeper(
		"ausdy",
		runtime.NewKVStoreService(key),
		runtime.ProvideEventService(),
		address.NewBech32Codec("noble"),
		nil,
	)

	bank = bank.WithSendCoinsRestriction(k.SendRestrictionFn)
	k.SetBankKeeper(bank)

	return k, testutil.DefaultContext(key, tkey)
}
