package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/usdy-noble/v2/x/aura/types/blocklist"
)

//

func (k *Keeper) GetBlocklistOwner(ctx context.Context) string {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return string(store.Get(blocklist.OwnerKey))
}

func (k *Keeper) SetBlocklistOwner(ctx context.Context, owner string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(blocklist.OwnerKey, []byte(owner))
}

//

func (k *Keeper) DeleteBlocklistPendingOwner(ctx context.Context) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(blocklist.PendingOwnerKey)
}

func (k *Keeper) GetBlocklistPendingOwner(ctx context.Context) string {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return string(store.Get(blocklist.PendingOwnerKey))
}

func (k *Keeper) SetBlocklistPendingOwner(ctx context.Context, pendingOwner string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(blocklist.PendingOwnerKey, []byte(pendingOwner))
}

//

func (k *Keeper) DeleteBlockedAddress(ctx context.Context, address []byte) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(blocklist.BlockedAddressKey(address))
}

func (k *Keeper) GetBlockedAddresses(ctx context.Context) (addresses []string) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, blocklist.BlockedAddressPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		addresses = append(addresses, sdk.AccAddress(itr.Key()).String())
	}

	return
}

func (k *Keeper) HasBlockedAddress(ctx context.Context, address []byte) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Has(blocklist.BlockedAddressKey(address))
}

func (k *Keeper) SetBlockedAddress(ctx context.Context, address []byte) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(blocklist.BlockedAddressKey(address), []byte{})
}
