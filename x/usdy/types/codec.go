package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/noble-assets/ondo/x/usdy/types/blocklist"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	blocklist.RegisterLegacyAminoCodec(cdc)

	cdc.RegisterConcrete(&MsgBurn{}, "ondo/usdy/Burn", nil)
	cdc.RegisterConcrete(&MsgMint{}, "ondo/usdy/Mint", nil)
	cdc.RegisterConcrete(&MsgPause{}, "ondo/usdy/Pause", nil)
	cdc.RegisterConcrete(&MsgUnpause{}, "ondo/usdy/Unpause", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	blocklist.RegisterInterfaces(registry)

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgBurn{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgMint{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgPause{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUnpause{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
