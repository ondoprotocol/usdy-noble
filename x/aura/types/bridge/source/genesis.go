package source

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Destinations: map[string]string{
			// https://docs.ondo.finance/addresses#ondo-bridge
			// https://axelar-api.polkachu.com/axelar/nexus/v1beta1/chains?status=1
			"Ethereum": "0xBd8Fb563a325dc853741907ae06e5F3c02c9235c",
			"mantle":   "0xd5235958c1F8a40641847A0E3BD51d04EFe9eC28",
		},
	}
}

func (gs *GenesisState) Validate() error {
	if _, err := sdk.AccAddressFromBech32(gs.Owner); err != nil {
		return fmt.Errorf("invalid owner address (%s): %s", gs.Owner, err)
	}

	if gs.Channel == "" {
		return errors.New("bridge channel must be defined")
	}

	return nil
}
