package source

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgBurnAndCallAxelar{}, "aura/bridge/BurnAndCallAxelar", nil)

	cdc.RegisterConcrete(&MsgPause{}, "aura/bridge/source/Pause", nil)
	cdc.RegisterConcrete(&MsgUnpause{}, "aura/bridge/source/Unpause", nil)

	cdc.RegisterConcrete(&MsgTransferOwnership{}, "aura/bridge/source/TransferOwnership", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgBurnAndCallAxelar{})

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgPause{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUnpause{})

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgTransferOwnership{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
