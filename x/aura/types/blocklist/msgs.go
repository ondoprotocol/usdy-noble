package blocklist

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)

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

func (*MsgTransferOwnership) Route() string { return SubmoduleName }

func (*MsgTransferOwnership) Type() string { return "aura/blocklist/TransferOwnership" }

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

func (*MsgAcceptOwnership) Route() string { return SubmoduleName }

func (*MsgAcceptOwnership) Type() string { return "aura/blocklist/AcceptOwnership" }

//

var _ legacytx.LegacyMsg = &MsgAddToBlocklist{}

func (msg *MsgAddToBlocklist) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	for _, account := range msg.Accounts {
		if _, err := sdk.AccAddressFromBech32(account); err != nil {
			return fmt.Errorf("invalid account address (%s): %w", account, err)
		}
	}

	return nil
}

func (msg *MsgAddToBlocklist) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgAddToBlocklist) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgAddToBlocklist) Route() string { return SubmoduleName }

func (*MsgAddToBlocklist) Type() string { return "aura/blocklist/AddToBlocklist" }

//

var _ legacytx.LegacyMsg = &MsgRemoveFromBlocklist{}

func (msg *MsgRemoveFromBlocklist) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	for _, account := range msg.Accounts {
		if _, err := sdk.AccAddressFromBech32(account); err != nil {
			return fmt.Errorf("invalid account address (%s): %w", account, err)
		}
	}

	return nil
}

func (msg *MsgRemoveFromBlocklist) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgRemoveFromBlocklist) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgRemoveFromBlocklist) Route() string { return SubmoduleName }

func (*MsgRemoveFromBlocklist) Type() string { return "aura/blocklist/RemoveFromBlocklist" }
