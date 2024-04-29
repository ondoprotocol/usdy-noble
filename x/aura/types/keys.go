package types

import authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

const ModuleName = "aura"

var ModuleAddress = authtypes.NewModuleAddress(ModuleName)

var (
	PausedKey    = []byte("paused")
	BurnerPrefix = []byte("burner/")
	MinterPrefix = []byte("minter/")
	PauserKey    = []byte("pauser")
)
