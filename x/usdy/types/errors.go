package types

import "cosmossdk.io/errors"

var ErrInvalidPauser = errors.Register(ModuleName, 1, "signer is not pauser")
