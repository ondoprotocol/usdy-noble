package blocklist

import (
	"fmt"

	"cosmossdk.io/core/address"
)

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func (gs *GenesisState) Validate(cdc address.Codec) error {
	if gs.Owner != "" {
		if _, err := cdc.StringToBytes(gs.Owner); err != nil {
			return fmt.Errorf("invalid blocklist owner address (%s): %s", gs.Owner, err)
		}
	}

	if gs.PendingOwner != "" {
		if _, err := cdc.StringToBytes(gs.PendingOwner); err != nil {
			return fmt.Errorf("invalid pending blocklist owner address (%s): %s", gs.PendingOwner, err)
		}
	}

	for _, address := range gs.BlockedAddresses {
		if _, err := cdc.StringToBytes(address); err != nil {
			return fmt.Errorf("invalid blocked address (%s): %s", address, err)
		}
	}

	return nil
}
