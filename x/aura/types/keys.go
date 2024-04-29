package types

import authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

const ModuleName = "aura"

var ModuleAddress = authtypes.NewModuleAddress(ModuleName)

var (
	PausedKey       = []byte("paused")
	OwnerKey        = []byte("owner")
	PendingOwnerKey = []byte("pending_owner")
	BurnerPrefix    = []byte("burner/")
	MinterPrefix    = []byte("minter/")
	PauserPrefix    = []byte("pauser/")
)
