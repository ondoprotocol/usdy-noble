package usdy

import (
	"context"

	"github.com/noble-assets/ondo/x/usdy/keeper"
	"github.com/noble-assets/ondo/x/usdy/types"
	"github.com/noble-assets/ondo/x/usdy/types/blocklist"
)

func InitGenesis(ctx context.Context, k *keeper.Keeper, accountKeeper types.AccountKeeper, genesis types.GenesisState) {
	_ = k.Paused.Set(ctx, genesis.Paused)
	_ = k.Burner.Set(ctx, genesis.Burner)
	_ = k.Minter.Set(ctx, genesis.Minter)
	_ = k.Pauser.Set(ctx, genesis.Pauser)

	_ = k.Owner.Set(ctx, genesis.BlocklistState.Owner)
	_ = k.PendingOwner.Set(ctx, genesis.BlocklistState.PendingOwner)

	for address, blocked := range genesis.BlocklistState.BlockedAddresses {
		if blocked {
			account, _ := accountKeeper.AddressCodec().StringToBytes(address)
			_ = k.BlockedAddresses.Set(ctx, account, blocked)
		}
	}
}

func ExportGenesis(ctx context.Context, k *keeper.Keeper, accountKeeper types.AccountKeeper) *types.GenesisState {
	paused, _ := k.Paused.Get(ctx)
	burner, _ := k.Burner.Get(ctx)
	minter, _ := k.Minter.Get(ctx)
	pauser, _ := k.Pauser.Get(ctx)

	owner, _ := k.Owner.Get(ctx)
	pendingOwner, _ := k.PendingOwner.Get(ctx)

	blockedAddresses := make(map[string]bool)
	_ = k.BlockedAddresses.Walk(ctx, nil, func(account []byte, blocked bool) (stop bool, err error) {
		address, _ := accountKeeper.AddressCodec().BytesToString(account)
		blockedAddresses[address] = blocked

		return false, nil
	})

	return &types.GenesisState{
		BlocklistState: blocklist.GenesisState{
			Owner:            owner,
			PendingOwner:     pendingOwner,
			BlockedAddresses: blockedAddresses,
		},
		Paused: paused,
		Burner: burner,
		Minter: minter,
		Pauser: pauser,
	}
}
