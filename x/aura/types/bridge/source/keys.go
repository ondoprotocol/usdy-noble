package source

import authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

const SubmoduleName = "aura-bridge-source"

var SubmoduleAddress = authtypes.NewModuleAddress(SubmoduleName)

var (
	PausedKey         = []byte("bridge/source/paused")
	OwnerKey          = []byte("bridge/source/owner")
	NonceKey          = []byte("bridge/source/nonce")
	ChannelKey        = []byte("bridge/source/channel")
	DestinationPrefix = []byte("bridge/source/destination/")
)

func DestinationKey(chain string) []byte {
	return append(DestinationPrefix, []byte(chain)...)
}
