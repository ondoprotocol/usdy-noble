package types

import "cosmossdk.io/errors"

var (
	ErrNoOwner               = errors.Register(ModuleName, 1, "there is no owner")
	ErrSameOwner             = errors.Register(ModuleName, 2, "provided owner is the current owner")
	ErrInvalidOwner          = errors.Register(ModuleName, 3, "signer is not owner")
	ErrNoPendingOwner        = errors.Register(ModuleName, 4, "there is no pending owner")
	ErrInvalidPendingOwner   = errors.Register(ModuleName, 5, "signer is not pending owner")
	ErrInvalidBurner         = errors.Register(ModuleName, 6, "signer is not a burner")
	ErrInvalidMinter         = errors.Register(ModuleName, 7, "signer is not a minter")
	ErrInvalidPauser         = errors.Register(ModuleName, 8, "signer is not a pauser")
	ErrInsufficientAllowance = errors.Register(ModuleName, 9, "insufficient allowance")
)
