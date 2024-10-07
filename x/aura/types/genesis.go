package types

import (
	"fmt"

	"cosmossdk.io/core/address"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/ondoprotocol/usdy-noble/v2/x/aura/types/blocklist"
)

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		BlocklistState: blocklist.DefaultGenesisState(),
		Paused:         false,
	}
}

func (gs *GenesisState) Validate(cdc address.Codec) error {
	if err := gs.BlocklistState.Validate(cdc); err != nil {
		return err
	}

	if gs.Owner != "" {
		if _, err := cdc.StringToBytes(gs.Owner); err != nil {
			return fmt.Errorf("invalid owner address (%s): %s", gs.Owner, err)
		}
	}

	if gs.PendingOwner != "" {
		if _, err := cdc.StringToBytes(gs.PendingOwner); err != nil {
			return fmt.Errorf("invalid pending owner address (%s): %s", gs.PendingOwner, err)
		}
	}

	for _, burner := range gs.Burners {
		if _, err := cdc.StringToBytes(burner.Address); err != nil {
			return fmt.Errorf("invalid burner address (%s): %s", burner.Address, err)
		}

		if burner.Allowance.IsNegative() {
			return fmt.Errorf("invalid burner allowance (%s)", burner.Address)
		}
	}

	for _, minter := range gs.Minters {
		if _, err := cdc.StringToBytes(minter.Address); err != nil {
			return fmt.Errorf("invalid minter address (%s): %s", minter.Address, err)
		}

		if minter.Allowance.IsNegative() {
			return fmt.Errorf("invalid minter allowance (%s)", minter.Address)
		}
	}

	for _, pauser := range gs.Pausers {
		if _, err := cdc.StringToBytes(pauser); err != nil {
			return fmt.Errorf("invalid pauser address (%s): %s", pauser, err)
		}
	}

	for _, channel := range gs.BlockedChannels {
		if !channeltypes.IsValidChannelID(channel) {
			return fmt.Errorf("invalid blocked channel (%s)", channel)
		}
	}

	return nil
}
