package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/usdy-noble/x/aura/types/blocklist"
)

//

func (k *Keeper) GetBlocklistOwner(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(blocklist.OwnerKey))
}

func (k *Keeper) SetBlocklistOwner(ctx sdk.Context, owner string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(blocklist.OwnerKey, []byte(owner))
}

//

func (k *Keeper) DeleteBlocklistPendingOwner(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(blocklist.PendingOwnerKey)
}

func (k *Keeper) GetBlocklistPendingOwner(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(blocklist.PendingOwnerKey))
}

func (k *Keeper) SetBlocklistPendingOwner(ctx sdk.Context, pendingOwner string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(blocklist.PendingOwnerKey, []byte(pendingOwner))
}

//

func (k *Keeper) DeleteBlockedAddress(ctx sdk.Context, address []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(blocklist.BlockedAddressKey(address))
}

func (k *Keeper) GetBlockedAddresses(ctx sdk.Context) (addresses []string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), blocklist.BlockedAddressPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		addresses = append(addresses, sdk.AccAddress(itr.Key()).String())
	}

	return
}

func (k *Keeper) HasBlockedAddress(ctx sdk.Context, address []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(blocklist.BlockedAddressKey(address))
}

func (k *Keeper) SetBlockedAddress(ctx sdk.Context, address []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(blocklist.BlockedAddressKey(address), []byte{})
}
