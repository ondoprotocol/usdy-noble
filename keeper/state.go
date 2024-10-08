package keeper

import (
	"context"

	"cosmossdk.io/math"
	"github.com/ondoprotocol/usdy-noble/v2/types"
)

//

func (k *Keeper) GetPaused(ctx context.Context) bool {
	paused, _ := k.Paused.Get(ctx)
	return paused
}

func (k *Keeper) SetPaused(ctx context.Context, paused bool) error {
	return k.Paused.Set(ctx, paused)
}

//

func (k *Keeper) GetOwner(ctx context.Context) string {
	owner, _ := k.Owner.Get(ctx)
	return owner
}

func (k *Keeper) SetOwner(ctx context.Context, owner string) error {
	return k.Owner.Set(ctx, owner)
}

//

func (k *Keeper) DeletePendingOwner(ctx context.Context) error {
	return k.PendingOwner.Remove(ctx)
}

func (k *Keeper) GetPendingOwner(ctx context.Context) string {
	pendingOwner, _ := k.PendingOwner.Get(ctx)
	return pendingOwner
}

func (k *Keeper) SetPendingOwner(ctx context.Context, pendingOwner string) error {
	return k.PendingOwner.Set(ctx, pendingOwner)
}

//

func (k *Keeper) DeleteBurner(ctx context.Context, burner string) error {
	return k.Burners.Remove(ctx, burner)
}

func (k *Keeper) GetBurner(ctx context.Context, burner string) (allowance math.Int) {
	allowance, err := k.Burners.Get(ctx, burner)
	if err != nil {
		return math.ZeroInt()
	}

	return
}

func (k *Keeper) GetBurners(ctx context.Context) (burners []types.Burner) {
	_ = k.Burners.Walk(ctx, nil, func(burner string, allowance math.Int) (stop bool, err error) {
		burners = append(burners, types.Burner{
			Address:   burner,
			Allowance: allowance,
		})

		return false, nil
	})

	return
}

func (k *Keeper) HasBurner(ctx context.Context, burner string) bool {
	has, _ := k.Burners.Has(ctx, burner)
	return has
}

func (k *Keeper) SetBurner(ctx context.Context, burner string, allowance math.Int) error {
	return k.Burners.Set(ctx, burner, allowance)
}

//

func (k *Keeper) DeleteMinter(ctx context.Context, minter string) error {
	return k.Minters.Remove(ctx, minter)
}

func (k *Keeper) GetMinter(ctx context.Context, minter string) (allowance math.Int) {
	allowance, err := k.Minters.Get(ctx, minter)
	if err != nil {
		return math.ZeroInt()
	}

	return
}

func (k *Keeper) GetMinters(ctx context.Context) (minters []types.Minter) {
	_ = k.Minters.Walk(ctx, nil, func(minter string, allowance math.Int) (stop bool, err error) {
		minters = append(minters, types.Minter{
			Address:   minter,
			Allowance: allowance,
		})

		return false, nil
	})

	return
}

func (k *Keeper) HasMinter(ctx context.Context, minter string) bool {
	has, _ := k.Minters.Has(ctx, minter)
	return has
}

func (k *Keeper) SetMinter(ctx context.Context, minter string, allowance math.Int) error {
	return k.Minters.Set(ctx, minter, allowance)
}

//

func (k *Keeper) DeletePauser(ctx context.Context, pauser string) error {
	return k.Pausers.Remove(ctx, pauser)
}

func (k *Keeper) GetPausers(ctx context.Context) (pausers []string) {
	_ = k.Pausers.Walk(ctx, nil, func(pauser string, _ []byte) (stop bool, err error) {
		pausers = append(pausers, pauser)
		return false, nil
	})

	return
}

func (k *Keeper) HasPauser(ctx context.Context, pauser string) bool {
	has, _ := k.Pausers.Has(ctx, pauser)
	return has
}

func (k *Keeper) SetPauser(ctx context.Context, pauser string) error {
	return k.Pausers.Set(ctx, pauser, []byte{})
}

//

func (k *Keeper) DeleteBlockedChannel(ctx context.Context, channel string) error {
	return k.BlockedChannels.Remove(ctx, channel)
}

func (k *Keeper) GetBlockedChannels(ctx context.Context) (channels []string) {
	_ = k.BlockedChannels.Walk(ctx, nil, func(channel string, _ []byte) (stop bool, err error) {
		channels = append(channels, channel)
		return false, nil
	})

	return
}

func (k *Keeper) HasBlockedChannel(ctx context.Context, channel string) bool {
	has, _ := k.BlockedChannels.Has(ctx, channel)
	return has
}

func (k *Keeper) SetBlockedChannel(ctx context.Context, channel string) error {
	return k.BlockedChannels.Set(ctx, channel, []byte{})
}
