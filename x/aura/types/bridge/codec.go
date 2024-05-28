package bridge

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/noble-assets/aura/x/aura/types/bridge/source"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	source.RegisterLegacyAminoCodec(cdc)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	source.RegisterInterfaces(registry)
}
