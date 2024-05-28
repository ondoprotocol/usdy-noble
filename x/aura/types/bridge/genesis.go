package bridge

import "github.com/noble-assets/aura/x/aura/types/bridge/source"

func DefaultGenesisState() GenesisState {
	return GenesisState{
		SourceState: source.DefaultGenesisState(),
	}
}

func (gs *GenesisState) Validate() error {
	if err := gs.SourceState.Validate(); err != nil {
		return err
	}

	return nil
}
