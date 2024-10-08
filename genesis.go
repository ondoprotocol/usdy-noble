package aura

import (
	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/usdy-noble/v2/keeper"
	"github.com/ondoprotocol/usdy-noble/v2/types"
	"github.com/ondoprotocol/usdy-noble/v2/types/blocklist"
)

func InitGenesis(ctx sdk.Context, k *keeper.Keeper, addressCodec address.Codec, genesis types.GenesisState) {
	if err := k.SetBlocklistOwner(ctx, genesis.BlocklistState.Owner); err != nil {
		panic(err)
	}
	if err := k.SetBlocklistPendingOwner(ctx, genesis.BlocklistState.PendingOwner); err != nil {
		panic(err)
	}
	for _, account := range genesis.BlocklistState.BlockedAddresses {
		address, _ := addressCodec.StringToBytes(account)
		if err := k.SetBlockedAddress(ctx, address); err != nil {
			panic(err)
		}
	}

	if err := k.SetPaused(ctx, genesis.Paused); err != nil {
		panic(err)
	}
	if err := k.SetOwner(ctx, genesis.Owner); err != nil {
		panic(err)
	}
	if err := k.SetPendingOwner(ctx, genesis.PendingOwner); err != nil {
		panic(err)
	}
	for _, burner := range genesis.Burners {
		if err := k.SetBurner(ctx, burner.Address, burner.Allowance); err != nil {
			panic(err)
		}
	}
	for _, minter := range genesis.Minters {
		if err := k.SetMinter(ctx, minter.Address, minter.Allowance); err != nil {
			panic(err)
		}
	}
	for _, pauser := range genesis.Pausers {
		if err := k.SetPauser(ctx, pauser); err != nil {
			panic(err)
		}
	}
}

func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		BlocklistState: blocklist.GenesisState{
			Owner:            k.GetBlocklistOwner(ctx),
			PendingOwner:     k.GetBlocklistPendingOwner(ctx),
			BlockedAddresses: k.GetBlockedAddresses(ctx),
		},
		Paused:       k.GetPaused(ctx),
		Owner:        k.GetOwner(ctx),
		PendingOwner: k.GetPendingOwner(ctx),
		Burners:      k.GetBurners(ctx),
		Minters:      k.GetMinters(ctx),
		Pausers:      k.GetPausers(ctx),
	}
}
