package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/aura/x/aura/types/bridge/source"
)

//

func (k *Keeper) GetBridgeSourcePaused(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(source.PausedKey)
	if len(bz) == 1 && bz[0] == 1 {
		return true
	} else {
		return false
	}
}

func (k *Keeper) SetBridgeSourcePaused(ctx sdk.Context, paused bool) {
	store := ctx.KVStore(k.storeKey)
	if paused {
		store.Set(source.PausedKey, []byte{0x1})
	} else {
		store.Set(source.PausedKey, []byte{0x0})
	}
}

//

func (k *Keeper) GetBridgeSourceOwner(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(source.OwnerKey))
}

func (k *Keeper) SetBridgeSourceOwner(ctx sdk.Context, owner string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(source.OwnerKey, []byte(owner))
}

//

func (k *Keeper) GetBridgeSourceNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(source.NonceKey)
	return binary.BigEndian.Uint64(bz)
}

func (k *Keeper) IncrementBridgeSourceNonce(ctx sdk.Context) uint64 {
	nonce := k.GetBridgeSourceNonce(ctx)
	k.SetBridgeSourceNonce(ctx, nonce+1)
	return nonce
}

func (k *Keeper) SetBridgeSourceNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, nonce)
	store.Set(source.NonceKey, bz)
}

//

func (k *Keeper) GetBridgeChannel(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(source.ChannelKey))
}

func (k *Keeper) SetBridgeChannel(ctx sdk.Context, channel string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(source.ChannelKey, []byte(channel))
}

//

func (k *Keeper) GetBridgeDestination(ctx sdk.Context, chain string) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(source.DestinationKey(chain)))
}

func (k *Keeper) GetBridgeDestinations(ctx sdk.Context) map[string]string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), source.DestinationPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	destinations := make(map[string]string)
	for ; itr.Valid(); itr.Next() {
		destinations[string(itr.Key())] = string(itr.Value())
	}

	return destinations
}

func (k *Keeper) SetBridgeDestination(ctx sdk.Context, chain string, destination string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(source.DestinationKey(chain), []byte(destination))
}
