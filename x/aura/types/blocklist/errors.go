package blocklist

import "cosmossdk.io/errors"

var (
	Codespace = "aura/blocklist"

	ErrNoOwner             = errors.Register(Codespace, 1, "there is no blocklist owner")
	ErrInvalidOwner        = errors.Register(Codespace, 2, "signer is not blocklist owner")
	ErrNoPendingOwner      = errors.Register(Codespace, 3, "there is no pending blocklist owner")
	ErrInvalidPendingOwner = errors.Register(Codespace, 4, "signer is not blocklist pending owner")
)
