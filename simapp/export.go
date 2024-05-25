package simapp

import (
	"encoding/json"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// ExportAppStateAndValidators exports the state of the application for a genesis
// file.
func (app *SimApp) ExportAppStateAndValidators(
	forZeroHeight bool, jailAllowedAddrs []string,
) (servertypes.ExportedApp, error) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	// We export at last height + 1, because that's the height at which
	// Tendermint will start InitChain.
	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		panic("zero height genesis is unsupported")
	}

	genState := app.mm.ExportGenesis(ctx, app.appCodec)
	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	validators, err := staking.WriteValidators(ctx, app.StakingKeeper)
	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, err
}
