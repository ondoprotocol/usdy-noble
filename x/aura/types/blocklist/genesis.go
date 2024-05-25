package blocklist

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func (gs *GenesisState) Validate() error {
	if _, err := sdk.AccAddressFromBech32(gs.Owner); err != nil {
		return fmt.Errorf("invalid blocklist owner address (%s): %s", gs.Owner, err)
	}

	if gs.PendingOwner != "" {
		if _, err := sdk.AccAddressFromBech32(gs.PendingOwner); err != nil {
			return fmt.Errorf("invalid pending blocklist owner address (%s): %s", gs.PendingOwner, err)
		}
	}

	for _, address := range gs.BlockedAddresses {
		if _, err := sdk.AccAddressFromBech32(address); err != nil {
			return fmt.Errorf("invalid blocked address (%s): %s", address, err)
		}
	}

	return nil
}
