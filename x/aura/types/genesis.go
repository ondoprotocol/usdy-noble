package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/aura/x/aura/types/blocklist"
)

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		BlocklistState: blocklist.DefaultGenesisState(),
		Paused:         false,
	}
}

func (gs *GenesisState) Validate() error {
	if _, err := sdk.AccAddressFromBech32(gs.Owner); err != nil {
		return fmt.Errorf("invalid owner address (%s): %s", gs.Owner, err)
	}

	if gs.PendingOwner != "" {
		if _, err := sdk.AccAddressFromBech32(gs.PendingOwner); err != nil {
			return fmt.Errorf("invalid pending owner address (%s): %s", gs.PendingOwner, err)
		}
	}

	for _, burner := range gs.Burners {
		if _, err := sdk.AccAddressFromBech32(burner.Address); err != nil {
			return fmt.Errorf("invalid burner address (%s): %s", burner.Address, err)
		}

		if burner.Allowance.IsNegative() {
			return fmt.Errorf("invalid burner allowance (%s)", burner.Address)
		}
	}

	for _, minter := range gs.Minters {
		if _, err := sdk.AccAddressFromBech32(minter.Address); err != nil {
			return fmt.Errorf("invalid minter address (%s): %s", minter.Address, err)
		}

		if minter.Allowance.IsNegative() {
			return fmt.Errorf("invalid minter allowance (%s)", minter.Address)
		}
	}

	for _, pauser := range gs.Pausers {
		if _, err := sdk.AccAddressFromBech32(pauser); err != nil {
			return fmt.Errorf("invalid pauser address (%s): %s", pauser, err)
		}
	}

	return gs.BlocklistState.Validate()
}
