package types

import (
	"cosmossdk.io/core/address"
)

type AccountKeeper interface {
	AddressCodec() address.Codec
}
