package utils

import (
	"cosmossdk.io/store/rootmulti"
	"cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetKVStore retrieves the KVStore for the specified module from the context.
func GetKVStore(ctx sdk.Context, moduleName string) types.KVStore {
	return ctx.KVStore(ctx.MultiStore().(*rootmulti.Store).StoreKeysByName()[moduleName])
}
