package usdy

import (
	"context"

	"github.com/noble-assets/ondo/x/usdy/keeper"
	"github.com/noble-assets/ondo/x/usdy/types"
)

func InitGenesis(ctx context.Context, k *keeper.Keeper, genesis types.GenesisState) {}

func ExportGenesis(ctx context.Context, k *keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{}
}
