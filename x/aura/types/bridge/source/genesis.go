package source

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func (gs *GenesisState) Validate() error {
	if _, err := sdk.AccAddressFromBech32(gs.Owner); err != nil {
		return fmt.Errorf("invalid owner address (%s): %s", gs.Owner, err)
	}

	return nil
}
