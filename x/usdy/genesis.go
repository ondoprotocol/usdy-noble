package usdy

import (
	"context"

	"github.com/noble-assets/ondo/x/usdy/keeper"
	"github.com/noble-assets/ondo/x/usdy/types"
)

func InitGenesis(ctx context.Context, k *keeper.Keeper, genesis types.GenesisState) {
	_ = k.Paused.Set(ctx, genesis.Paused)
	_ = k.Pauser.Set(ctx, genesis.Pauser)
}

func ExportGenesis(ctx context.Context, k *keeper.Keeper) *types.GenesisState {
	paused, _ := k.Paused.Get(ctx)
	pauser, _ := k.Pauser.Get(ctx)

	return &types.GenesisState{
		Paused: paused,
		Pauser: pauser,
	}
}
