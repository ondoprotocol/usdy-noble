package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
)

//

var _ legacytx.LegacyMsg = &MsgBurn{}

func (msg *MsgBurn) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.From); err != nil {
		return fmt.Errorf("invalid from address (%s): %w", msg.From, err)
	}

	return nil
}

func (msg *MsgBurn) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgBurn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgBurn) Route() string { return ModuleName }

func (*MsgBurn) Type() string { return "aura/Burn" }

//

var _ legacytx.LegacyMsg = &MsgMint{}

func (msg *MsgMint) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.To); err != nil {
		return fmt.Errorf("invalid to address (%s): %w", msg.To, err)
	}

	return nil
}

func (msg *MsgMint) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgMint) Route() string { return ModuleName }

func (*MsgMint) Type() string { return "aura/Mint" }

//

var _ legacytx.LegacyMsg = &MsgPause{}

func (msg *MsgPause) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	return nil
}

func (msg *MsgPause) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgPause) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgPause) Route() string { return ModuleName }

func (*MsgPause) Type() string { return "aura/Pause" }

//

var _ legacytx.LegacyMsg = &MsgUnpause{}

func (msg *MsgUnpause) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	return nil
}

func (msg *MsgUnpause) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgUnpause) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgUnpause) Route() string { return ModuleName }

func (*MsgUnpause) Type() string { return "aura/Unpause" }

//

var _ legacytx.LegacyMsg = &MsgTransferOwnership{}

func (msg *MsgTransferOwnership) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.NewOwner); err != nil {
		return fmt.Errorf("invalid new owner address (%s): %w", msg.NewOwner, err)
	}

	return nil
}

func (msg *MsgTransferOwnership) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgTransferOwnership) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgTransferOwnership) Route() string { return ModuleName }

func (*MsgTransferOwnership) Type() string { return "aura/TransferOwnership" }

//

var _ legacytx.LegacyMsg = &MsgAcceptOwnership{}

func (msg *MsgAcceptOwnership) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	return nil
}

func (msg *MsgAcceptOwnership) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgAcceptOwnership) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgAcceptOwnership) Route() string { return ModuleName }

func (*MsgAcceptOwnership) Type() string { return "aura/AcceptOwnership" }

//

var _ legacytx.LegacyMsg = &MsgAddBurner{}

func (msg *MsgAddBurner) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Burner); err != nil {
		return fmt.Errorf("invalid burner address (%s): %w", msg.Burner, err)
	}

	return nil
}

func (msg *MsgAddBurner) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgAddBurner) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgAddBurner) Route() string { return ModuleName }

func (*MsgAddBurner) Type() string { return "aura/AddBurner" }

//

var _ legacytx.LegacyMsg = &MsgRemoveBurner{}

func (msg *MsgRemoveBurner) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Burner); err != nil {
		return fmt.Errorf("invalid burner address (%s): %w", msg.Burner, err)
	}

	return nil
}

func (msg *MsgRemoveBurner) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgRemoveBurner) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgRemoveBurner) Route() string { return ModuleName }

func (*MsgRemoveBurner) Type() string { return "aura/RemoveBurner" }

//

var _ legacytx.LegacyMsg = &MsgSetBurnerAllowance{}

func (msg *MsgSetBurnerAllowance) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Burner); err != nil {
		return fmt.Errorf("invalid burner address (%s): %w", msg.Burner, err)
	}

	return nil
}

func (msg *MsgSetBurnerAllowance) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgSetBurnerAllowance) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgSetBurnerAllowance) Route() string { return ModuleName }

func (*MsgSetBurnerAllowance) Type() string { return "aura/SetBurnerAllowance" }

//

var _ legacytx.LegacyMsg = &MsgAddMinter{}

func (msg *MsgAddMinter) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Minter); err != nil {
		return fmt.Errorf("invalid minter address (%s): %w", msg.Minter, err)
	}

	return nil
}

func (msg *MsgAddMinter) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgAddMinter) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgAddMinter) Route() string { return ModuleName }

func (*MsgAddMinter) Type() string { return "aura/AddMinter" }

//

var _ legacytx.LegacyMsg = &MsgRemoveMinter{}

func (msg *MsgRemoveMinter) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Minter); err != nil {
		return fmt.Errorf("invalid minter address (%s): %w", msg.Minter, err)
	}

	return nil
}

func (msg *MsgRemoveMinter) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgRemoveMinter) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgRemoveMinter) Route() string { return ModuleName }

func (*MsgRemoveMinter) Type() string { return "aura/RemoveMinter" }

//

var _ legacytx.LegacyMsg = &MsgSetMinterAllowance{}

func (msg *MsgSetMinterAllowance) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Minter); err != nil {
		return fmt.Errorf("invalid minter address (%s): %w", msg.Minter, err)
	}

	return nil
}

func (msg *MsgSetMinterAllowance) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgSetMinterAllowance) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgSetMinterAllowance) Route() string { return ModuleName }

func (*MsgSetMinterAllowance) Type() string { return "aura/SetMinterAllowance" }

//

var _ legacytx.LegacyMsg = &MsgAddPauser{}

func (msg *MsgAddPauser) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Pauser); err != nil {
		return fmt.Errorf("invalid pauser address (%s): %w", msg.Pauser, err)
	}

	return nil
}

func (msg *MsgAddPauser) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgAddPauser) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgAddPauser) Route() string { return ModuleName }

func (*MsgAddPauser) Type() string { return "aura/AddPauser" }

//

var _ legacytx.LegacyMsg = &MsgRemovePauser{}

func (msg *MsgRemovePauser) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Pauser); err != nil {
		return fmt.Errorf("invalid pauser address (%s): %w", msg.Pauser, err)
	}

	return nil
}

func (msg *MsgRemovePauser) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgRemovePauser) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgRemovePauser) Route() string { return ModuleName }

func (*MsgRemovePauser) Type() string { return "aura/RemovePauser" }

//

var _ legacytx.LegacyMsg = &MsgAllowChannel{}

func (msg *MsgAllowChannel) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if !channeltypes.IsValidChannelID(msg.Channel) {
		return fmt.Errorf("invalid channel (%s)", msg.Channel)
	}

	return nil
}

func (msg *MsgAllowChannel) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgAllowChannel) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgAllowChannel) Route() string { return ModuleName }

func (*MsgAllowChannel) Type() string { return "aura/AllowChannel" }
