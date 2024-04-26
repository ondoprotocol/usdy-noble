package blocklist

import "cosmossdk.io/errors"

var (
	Codespace = "aura/blocklist"

	ErrInvalidOwner        = errors.Register(Codespace, 1, "signer is not blocklist owner")
	ErrInvalidPendingOwner = errors.Register(Codespace, 2, "signer is not blocklist pending owner")
)
