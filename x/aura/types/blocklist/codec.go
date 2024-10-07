package blocklist

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgTransferOwnership{}, "aura/blocklist/TransferOwnership", nil)
	cdc.RegisterConcrete(&MsgAcceptOwnership{}, "aura/blocklist/AcceptOwnership", nil)

	cdc.RegisterConcrete(&MsgAddToBlocklist{}, "aura/blocklist/AddToBlocklist", nil)
	cdc.RegisterConcrete(&MsgRemoveFromBlocklist{}, "aura/blocklist/RemoveFromBlocklist", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgTransferOwnership{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgAcceptOwnership{})

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgAddToBlocklist{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgRemoveFromBlocklist{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
