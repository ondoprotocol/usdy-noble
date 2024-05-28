package source

const SubmoduleName = "aura-bridge-source"

var (
	PausedKey         = []byte("bridge/source/paused")
	OwnerKey          = []byte("bridge/source/owner")
	NonceKey          = []byte("bridge/source/nonce")
	DestinationPrefix = []byte("bridge/source/destination/")
)

func DestinationKey(chain string) []byte {
	return append(DestinationPrefix, []byte(chain)...)
}
