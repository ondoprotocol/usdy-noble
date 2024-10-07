package keeper

import (
	"context"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/ondoprotocol/usdy-noble/v2/types"
)

//

func (k *Keeper) GetPaused(ctx context.Context) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.PausedKey)
	if len(bz) == 1 && bz[0] == 1 {
		return true
	} else {
		return false
	}
}

func (k *Keeper) SetPaused(ctx context.Context, paused bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	if paused {
		store.Set(types.PausedKey, []byte{0x1})
	} else {
		store.Set(types.PausedKey, []byte{0x0})
	}
}

//

func (k *Keeper) GetOwner(ctx context.Context) string {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return string(store.Get(types.OwnerKey))
}

func (k *Keeper) SetOwner(ctx context.Context, owner string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.OwnerKey, []byte(owner))
}

//

func (k *Keeper) DeletePendingOwner(ctx context.Context) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.PendingOwnerKey)
}

func (k *Keeper) GetPendingOwner(ctx context.Context) string {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return string(store.Get(types.PendingOwnerKey))
}

func (k *Keeper) SetPendingOwner(ctx context.Context, pendingOwner string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.PendingOwnerKey, []byte(pendingOwner))
}

//

func (k *Keeper) DeleteBurner(ctx context.Context, burner string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.BurnerKey(burner))
}

func (k *Keeper) GetBurner(ctx context.Context, burner string) (allowance math.Int) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.BurnerKey(burner))

	_ = allowance.Unmarshal(bz)
	return
}

func (k *Keeper) GetBurners(ctx context.Context) (burners []types.Burner) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.BurnerPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		var allowance math.Int
		_ = allowance.Unmarshal(itr.Value())

		burners = append(burners, types.Burner{
			Address:   string(itr.Key()),
			Allowance: allowance,
		})
	}

	return
}

func (k *Keeper) HasBurner(ctx context.Context, burner string) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Has(types.BurnerKey(burner))
}

func (k *Keeper) SetBurner(ctx context.Context, burner string, allowance math.Int) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz, _ := allowance.Marshal()
	store.Set(types.BurnerKey(burner), bz)
}

//

func (k *Keeper) DeleteMinter(ctx context.Context, minter string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.MinterKey(minter))
}

func (k *Keeper) GetMinter(ctx context.Context, minter string) (allowance math.Int) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.MinterKey(minter))

	_ = allowance.Unmarshal(bz)
	return
}

func (k *Keeper) GetMinters(ctx context.Context) (minters []types.Minter) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.MinterPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		var allowance math.Int
		_ = allowance.Unmarshal(itr.Value())

		minters = append(minters, types.Minter{
			Address:   string(itr.Key()),
			Allowance: allowance,
		})
	}

	return
}

func (k *Keeper) HasMinter(ctx context.Context, minter string) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Has(types.MinterKey(minter))
}

func (k *Keeper) SetMinter(ctx context.Context, minter string, allowance math.Int) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz, _ := allowance.Marshal()
	store.Set(types.MinterKey(minter), bz)
}

//

func (k *Keeper) DeletePauser(ctx context.Context, pauser string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.PauserKey(pauser))
}

func (k *Keeper) GetPausers(ctx context.Context) (pausers []string) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.PauserPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		pausers = append(pausers, string(itr.Key()))
	}

	return
}

func (k *Keeper) HasPauser(ctx context.Context, pauser string) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Has(types.PauserKey(pauser))
}

func (k *Keeper) SetPauser(ctx context.Context, pauser string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.PauserKey(pauser), []byte{})
}

//

func (k *Keeper) DeleteBlockedChannel(ctx context.Context, channel string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.BlockedChannelKey(channel))
}

func (k *Keeper) GetBlockedChannels(ctx context.Context) (channels []string) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.BlockedChannelPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		channels = append(channels, string(itr.Key()))
	}

	return
}

func (k *Keeper) HasBlockedChannel(ctx context.Context, channel string) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Has(types.BlockedChannelKey(channel))
}

func (k *Keeper) SetBlockedChannel(ctx context.Context, channel string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.BlockedChannelKey(channel), []byte{})
}
