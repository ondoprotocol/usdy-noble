package blocklist

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgTransferOwnership{}, "ondo/usdy/blocklist/TransferOwnership", nil)
	cdc.RegisterConcrete(&MsgAcceptOwnership{}, "ondo/usdy/blocklist/AcceptOwnership", nil)

	cdc.RegisterConcrete(&MsgAddToBlocklist{}, "ondo/usdy/blocklist/AddToBlocklist", nil)
	cdc.RegisterConcrete(&MsgRemoveFromBlocklist{}, "ondo/usdy/blocklist/RemoveFromBlocklist", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgTransferOwnership{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgAcceptOwnership{})

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgAddToBlocklist{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgRemoveFromBlocklist{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
