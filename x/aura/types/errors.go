package types

import "cosmossdk.io/errors"

var (
	ErrInvalidBurner = errors.Register(ModuleName, 1, "signer is not a burner")
	ErrInvalidMinter = errors.Register(ModuleName, 2, "signer is not a minter")
	ErrInvalidPauser = errors.Register(ModuleName, 3, "signer is not pauser")
)
