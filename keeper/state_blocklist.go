package keeper

import (
	"context"
)

//

func (k *Keeper) GetBlocklistOwner(ctx context.Context) string {
	owner, _ := k.BlocklistOwner.Get(ctx)
	return owner
}

func (k *Keeper) SetBlocklistOwner(ctx context.Context, owner string) error {
	return k.BlocklistOwner.Set(ctx, owner)
}

//

func (k *Keeper) DeleteBlocklistPendingOwner(ctx context.Context) error {
	return k.BlocklistPendingOwner.Remove(ctx)
}

func (k *Keeper) GetBlocklistPendingOwner(ctx context.Context) string {
	pendingOwner, _ := k.BlocklistPendingOwner.Get(ctx)
	return pendingOwner
}

func (k *Keeper) SetBlocklistPendingOwner(ctx context.Context, pendingOwner string) error {
	return k.BlocklistPendingOwner.Set(ctx, pendingOwner)
}

//

func (k *Keeper) DeleteBlockedAddress(ctx context.Context, address []byte) error {
	return k.BlockedAddresses.Remove(ctx, address)
}

func (k *Keeper) GetBlockedAddresses(ctx context.Context) (addresses []string) {
	_ = k.BlockedAddresses.Walk(ctx, nil, func(rawAddress []byte, _ []byte) (stop bool, err error) {
		address, _ := k.addressCodec.BytesToString(rawAddress)
		addresses = append(addresses, address)

		return false, nil
	})

	return
}

func (k *Keeper) HasBlockedAddress(ctx context.Context, address []byte) bool {
	has, _ := k.BlockedAddresses.Has(ctx, address)
	return has
}

func (k *Keeper) SetBlockedAddress(ctx context.Context, address []byte) error {
	return k.BlockedAddresses.Set(ctx, address, []byte{})
}
