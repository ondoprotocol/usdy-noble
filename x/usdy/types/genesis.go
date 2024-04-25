package types

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Paused: false,
	}
}

func (gs *GenesisState) Validate() error {
	return nil
}
