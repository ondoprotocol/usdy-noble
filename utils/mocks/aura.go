package mocks

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/aura/x/aura/keeper"
	"github.com/noble-assets/aura/x/aura/types"
)

func AuraKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	return AuraKeeperWithBank(t, BankKeeper{
		Restriction: NoOpSendRestrictionFn,
	})
}

func AuraKeeperWithBank(_ testing.TB, bank BankKeeper) (*keeper.Keeper, sdk.Context) {
	key := storetypes.NewKVStoreKey(types.ModuleName)
	tkey := storetypes.NewTransientStoreKey("transient_aura")

	k := keeper.NewKeeper(
		codec.NewProtoCodec(codectypes.NewInterfaceRegistry()),
		key,
		"ausdy",
		nil,
	)

	bank = bank.WithSendCoinsRestriction(k.SendRestrictionFn)
	k.SetBankKeeper(bank)

	return k, testutil.DefaultContext(key, tkey)
}
