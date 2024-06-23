package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/noble-assets/aura/x/aura/types/blocklist"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	blocklist.RegisterLegacyAminoCodec(cdc)

	cdc.RegisterConcrete(&MsgBurn{}, "aura/Burn", nil)
	cdc.RegisterConcrete(&MsgMint{}, "aura/Mint", nil)

	cdc.RegisterConcrete(&MsgPause{}, "aura/Pause", nil)
	cdc.RegisterConcrete(&MsgUnpause{}, "aura/Unpause", nil)

	cdc.RegisterConcrete(&MsgTransferOwnership{}, "aura/TransferOwnership", nil)
	cdc.RegisterConcrete(&MsgAcceptOwnership{}, "aura/AcceptOwnership", nil)

	cdc.RegisterConcrete(&MsgAddBurner{}, "aura/AddBurner", nil)
	cdc.RegisterConcrete(&MsgRemoveBurner{}, "aura/RemoveBurner", nil)
	cdc.RegisterConcrete(&MsgSetBurnerAllowance{}, "aura/SetBurnerAllowance", nil)

	cdc.RegisterConcrete(&MsgAddMinter{}, "aura/AddMinter", nil)
	cdc.RegisterConcrete(&MsgRemoveMinter{}, "aura/RemoveMinter", nil)
	cdc.RegisterConcrete(&MsgSetMinterAllowance{}, "aura/SetMinterAllowance", nil)

	cdc.RegisterConcrete(&MsgAddPauser{}, "aura/AddPauser", nil)
	cdc.RegisterConcrete(&MsgRemovePauser{}, "aura/RemovePauser", nil)

	cdc.RegisterConcrete(&MsgAllowChannel{}, "aura/AllowChannel", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	blocklist.RegisterInterfaces(registry)

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgBurn{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgMint{})

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgPause{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUnpause{})

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgTransferOwnership{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgAcceptOwnership{})

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgAddBurner{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgRemoveBurner{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSetBurnerAllowance{})

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgAddMinter{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgRemoveMinter{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSetMinterAllowance{})

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgAddPauser{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgRemovePauser{})

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgAllowChannel{})

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
