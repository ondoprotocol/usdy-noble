package blocklist

const SubmoduleName = "aura-blocklist"

var (
	OwnerKey             = []byte("blocklist/owner")
	PendingOwnerKey      = []byte("blocklist/pending_owner")
	BlockedAddressPrefix = []byte("blocklist/blocked_address/")
)

func BlockedAddressKey(address []byte) []byte {
	return append(BlockedAddressPrefix, address...)
}
