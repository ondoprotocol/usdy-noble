package types

import "cosmossdk.io/errors"

var (
	ErrInvalidBurner = errors.Register(ModuleName, 1, "signer is not burner")
	ErrInvalidMinter = errors.Register(ModuleName, 2, "signer is not minter")
	ErrInvalidPauser = errors.Register(ModuleName, 3, "signer is not pauser")
)
