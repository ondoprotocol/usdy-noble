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
	ChannelPrefix   = []byte("channel/")
)

func BurnerKey(address string) []byte {
	return append(BurnerPrefix, []byte(address)...)
}

func MinterKey(address string) []byte {
	return append(MinterPrefix, []byte(address)...)
}

func PauserKey(address string) []byte {
	return append(PauserPrefix, []byte(address)...)
}

func ChannelKey(channel string) []byte {
	return append(ChannelPrefix, []byte(channel)...)
}
