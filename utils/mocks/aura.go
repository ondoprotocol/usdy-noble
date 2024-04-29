package mocks

import (
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/aura/x/aura/keeper"
	"github.com/noble-assets/aura/x/aura/types"
)

func AuraKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	return AuraKeeperWithBank(t, BankKeeper{})
}

func AuraKeeperWithBank(t testing.TB, bank types.BankKeeper) (*keeper.Keeper, sdk.Context) {
	logger := log.NewNopLogger()

	key := storetypes.NewKVStoreKey(types.ModuleName)
	tkey := storetypes.NewTransientStoreKey("transient_aura")
	wrapper := testutil.DefaultContextWithDB(t, key, tkey)

	return keeper.NewKeeper(
		codec.NewProtoCodec(codectypes.NewInterfaceRegistry()),
		logger,
		runtime.NewKVStoreService(key),
		runtime.ProvideEventService(),
		"ausdy",
		AccountKeeper{},
		bank,
	), wrapper.Ctx
}
