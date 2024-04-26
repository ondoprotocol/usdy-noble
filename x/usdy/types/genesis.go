package types

import "github.com/noble-assets/ondo/x/usdy/types/blocklist"

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		BlocklistState: blocklist.DefaultGenesisState(),
		Paused:         false,
	}
}

func (gs *GenesisState) Validate() error {
	return nil
}
