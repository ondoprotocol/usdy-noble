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
	"github.com/noble-assets/ondo/x/usdy/keeper"
	"github.com/noble-assets/ondo/x/usdy/types"
)

func USDYKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	logger := log.NewNopLogger()

	key := storetypes.NewKVStoreKey(types.ModuleName)
	tkey := storetypes.NewTransientStoreKey("transient_usdy")
	wrapper := testutil.DefaultContextWithDB(t, key, tkey)

	return keeper.NewKeeper(
		codec.NewProtoCodec(codectypes.NewInterfaceRegistry()),
		logger,
		runtime.NewKVStoreService(key),
		runtime.ProvideEventService(),
		"ausdy",
		AccountKeeper{},
		BankKeeper{},
	), wrapper.Ctx
}
