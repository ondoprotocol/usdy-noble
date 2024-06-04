package aura

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/aura/x/aura/keeper"
	"github.com/noble-assets/aura/x/aura/types"
	"github.com/noble-assets/aura/x/aura/types/blocklist"
	"github.com/noble-assets/aura/x/aura/types/bridge"
	"github.com/noble-assets/aura/x/aura/types/bridge/source"
)

func InitGenesis(ctx sdk.Context, k *keeper.Keeper, genesis types.GenesisState) {
	k.SetBlocklistOwner(ctx, genesis.BlocklistState.Owner)
	k.SetBlocklistPendingOwner(ctx, genesis.BlocklistState.PendingOwner)
	for _, account := range genesis.BlocklistState.BlockedAddresses {
		address, _ := sdk.AccAddressFromBech32(account)
		k.SetBlockedAddress(ctx, address)
	}

	k.SetBridgeSourcePaused(ctx, genesis.BridgeState.SourceState.Paused)
	k.SetBridgeSourceOwner(ctx, genesis.BridgeState.SourceState.Owner)
	k.SetBridgeSourceNonce(ctx, genesis.BridgeState.SourceState.Nonce)
	k.SetBridgeChannel(ctx, genesis.BridgeState.SourceState.Channel)
	for chain, destination := range genesis.BridgeState.SourceState.Destinations {
		k.SetBridgeDestination(ctx, chain, destination)
	}

	k.SetPaused(ctx, genesis.Paused)
	k.SetOwner(ctx, genesis.Owner)
	k.SetPendingOwner(ctx, genesis.PendingOwner)
	for _, burner := range genesis.Burners {
		k.SetBurner(ctx, burner.Address, burner.Allowance)
	}
	for _, minter := range genesis.Minters {
		k.SetMinter(ctx, minter.Address, minter.Allowance)
	}
	for _, pauser := range genesis.Pausers {
		k.SetPauser(ctx, pauser)
	}
}

func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		BlocklistState: blocklist.GenesisState{
			Owner:            k.GetBlocklistOwner(ctx),
			PendingOwner:     k.GetBlocklistPendingOwner(ctx),
			BlockedAddresses: k.GetBlockedAddresses(ctx),
		},
		BridgeState: bridge.GenesisState{
			SourceState: source.GenesisState{
				Paused:       k.GetBridgeSourcePaused(ctx),
				Owner:        k.GetBridgeSourceOwner(ctx),
				Nonce:        k.GetBridgeSourceNonce(ctx),
				Channel:      k.GetBridgeChannel(ctx),
				Destinations: k.GetBridgeDestinations(ctx),
			},
		},
		Paused:       k.GetPaused(ctx),
		Owner:        k.GetOwner(ctx),
		PendingOwner: k.GetPendingOwner(ctx),
		Burners:      k.GetBurners(ctx),
		Minters:      k.GetMinters(ctx),
		Pausers:      k.GetPausers(ctx),
	}
}
