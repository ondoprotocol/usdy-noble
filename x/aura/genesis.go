package aura

import (
	"context"

	"github.com/noble-assets/aura/x/aura/keeper"
	"github.com/noble-assets/aura/x/aura/types"
	"github.com/noble-assets/aura/x/aura/types/blocklist"
)

func InitGenesis(ctx context.Context, k *keeper.Keeper, accountKeeper types.AccountKeeper, genesis types.GenesisState) {
	_ = k.Paused.Set(ctx, genesis.Paused)
	_ = k.Owner.Set(ctx, genesis.Owner)
	_ = k.PendingOwner.Set(ctx, genesis.PendingOwner)
	for _, burner := range genesis.Burners {
		_ = k.Burners.Set(ctx, burner)
	}
	for _, minter := range genesis.Minters {
		_ = k.Minters.Set(ctx, minter)
	}
	for _, pauser := range genesis.Pausers {
		_ = k.Pausers.Set(ctx, pauser)
	}

	_ = k.BlocklistOwner.Set(ctx, genesis.BlocklistState.Owner)
	_ = k.BlocklistPendingOwner.Set(ctx, genesis.BlocklistState.PendingOwner)

	for address, blocked := range genesis.BlocklistState.BlockedAddresses {
		if blocked {
			account, _ := accountKeeper.AddressCodec().StringToBytes(address)
			_ = k.BlockedAddresses.Set(ctx, account, blocked)
		}
	}
}

func ExportGenesis(ctx context.Context, k *keeper.Keeper, accountKeeper types.AccountKeeper) *types.GenesisState {
	paused, _ := k.Paused.Get(ctx)
	owner, _ := k.Owner.Get(ctx)
	pendingOwner, _ := k.PendingOwner.Get(ctx)
	burners, _ := k.GetBurners(ctx)
	minters, _ := k.GetMinters(ctx)
	pausers, _ := k.GetPausers(ctx)

	blocklistOwner, _ := k.BlocklistOwner.Get(ctx)
	blocklistPendingOwner, _ := k.BlocklistPendingOwner.Get(ctx)

	blockedAddresses := make(map[string]bool)
	_ = k.BlockedAddresses.Walk(ctx, nil, func(account []byte, blocked bool) (stop bool, err error) {
		address, _ := accountKeeper.AddressCodec().BytesToString(account)
		blockedAddresses[address] = blocked

		return false, nil
	})

	return &types.GenesisState{
		BlocklistState: blocklist.GenesisState{
			Owner:            blocklistOwner,
			PendingOwner:     blocklistPendingOwner,
			BlockedAddresses: blockedAddresses,
		},
		Paused:       paused,
		Owner:        owner,
		PendingOwner: pendingOwner,
		Burners:      burners,
		Minters:      minters,
		Pausers:      pausers,
	}
}
