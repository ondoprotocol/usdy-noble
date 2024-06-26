package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/aura/x/aura/types"
)

//

func (k *Keeper) GetPaused(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PausedKey)
	if len(bz) == 1 && bz[0] == 1 {
		return true
	} else {
		return false
	}
}

func (k *Keeper) SetPaused(ctx sdk.Context, paused bool) {
	store := ctx.KVStore(k.storeKey)
	if paused {
		store.Set(types.PausedKey, []byte{0x1})
	} else {
		store.Set(types.PausedKey, []byte{0x0})
	}
}

//

func (k *Keeper) GetOwner(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.OwnerKey))
}

func (k *Keeper) SetOwner(ctx sdk.Context, owner string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OwnerKey, []byte(owner))
}

//

func (k *Keeper) DeletePendingOwner(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PendingOwnerKey)
}

func (k *Keeper) GetPendingOwner(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PendingOwnerKey))
}

func (k *Keeper) SetPendingOwner(ctx sdk.Context, pendingOwner string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PendingOwnerKey, []byte(pendingOwner))
}

//

func (k *Keeper) DeleteBurner(ctx sdk.Context, burner string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.BurnerKey(burner))
}

func (k *Keeper) GetBurner(ctx sdk.Context, burner string) (allowance sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BurnerKey(burner))

	_ = allowance.Unmarshal(bz)
	return
}

func (k *Keeper) GetBurners(ctx sdk.Context) (burners []types.Burner) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BurnerPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		var allowance sdk.Int
		_ = allowance.Unmarshal(itr.Value())

		burners = append(burners, types.Burner{
			Address:   string(itr.Key()),
			Allowance: allowance,
		})
	}

	return
}

func (k *Keeper) HasBurner(ctx sdk.Context, burner string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.BurnerKey(burner))
}

func (k *Keeper) SetBurner(ctx sdk.Context, burner string, allowance sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := allowance.Marshal()
	store.Set(types.BurnerKey(burner), bz)
}

//

func (k *Keeper) DeleteMinter(ctx sdk.Context, minter string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.MinterKey(minter))
}

func (k *Keeper) GetMinter(ctx sdk.Context, minter string) (allowance sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MinterKey(minter))

	_ = allowance.Unmarshal(bz)
	return
}

func (k *Keeper) GetMinters(ctx sdk.Context) (minters []types.Minter) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MinterPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		var allowance sdk.Int
		_ = allowance.Unmarshal(itr.Value())

		minters = append(minters, types.Minter{
			Address:   string(itr.Key()),
			Allowance: allowance,
		})
	}

	return
}

func (k *Keeper) HasMinter(ctx sdk.Context, minter string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.MinterKey(minter))
}

func (k *Keeper) SetMinter(ctx sdk.Context, minter string, allowance sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := allowance.Marshal()
	store.Set(types.MinterKey(minter), bz)
}

//

func (k *Keeper) DeletePauser(ctx sdk.Context, pauser string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PauserKey(pauser))
}

func (k *Keeper) GetPausers(ctx sdk.Context) (pausers []string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PauserPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		pausers = append(pausers, string(itr.Key()))
	}

	return
}

func (k *Keeper) HasPauser(ctx sdk.Context, pauser string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.PauserKey(pauser))
}

func (k *Keeper) SetPauser(ctx sdk.Context, pauser string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PauserKey(pauser), []byte{})
}

//

func (k *Keeper) DeleteBlockedChannel(ctx sdk.Context, channel string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.BlockedChannelKey(channel))
}

func (k *Keeper) GetBlockedChannels(ctx sdk.Context) (channels []string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BlockedChannelPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		channels = append(channels, string(itr.Key()))
	}

	return
}

func (k *Keeper) HasBlockedChannel(ctx sdk.Context, channel string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.BlockedChannelKey(channel))
}

func (k *Keeper) SetBlockedChannel(ctx sdk.Context, channel string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.BlockedChannelKey(channel), []byte{})
}
