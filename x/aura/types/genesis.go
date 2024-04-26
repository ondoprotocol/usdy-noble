package types

import "github.com/noble-assets/aura/x/aura/types/blocklist"

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		BlocklistState: blocklist.DefaultGenesisState(),
		Paused:         false,
	}
}

func (gs *GenesisState) Validate() error {
	return nil
}
