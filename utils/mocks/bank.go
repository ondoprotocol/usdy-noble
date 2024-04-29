package mocks

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/noble-assets/aura/x/aura/types"
)

var cdc = codec.NewBech32Codec("noble")

var _ types.BankKeeper = BankKeeper{}

type BankKeeper struct {
	Balances map[string]sdk.Coins
}

func (k BankKeeper) BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	balance := k.Balances[moduleName]
	newBalance, negative := balance.SafeSub(amt...)
	if negative {
		return sdkerrors.Wrapf(errors.ErrInsufficientFunds, "spendable balance %s is smaller than %s", balance, amt)
	}

	k.Balances[moduleName] = newBalance

	return nil
}

func (k BankKeeper) MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	k.Balances[moduleName] = k.Balances[moduleName].Add(amt...)

	return nil
}

func (k BankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	sender, _ := cdc.BytesToString(senderAddr)

	balance := k.Balances[sender]
	newBalance, negative := balance.SafeSub(amt...)
	if negative {
		return sdkerrors.Wrapf(errors.ErrInsufficientFunds, "spendable balance %s is smaller than %s", balance, amt)
	}

	k.Balances[sender] = newBalance
	k.Balances[recipientModule] = k.Balances[recipientModule].Add(amt...)

	return nil
}

func (k BankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	recipient, _ := cdc.BytesToString(recipientAddr)

	balance := k.Balances[senderModule]
	newBalance, negative := balance.SafeSub(amt...)
	if negative {
		return sdkerrors.Wrapf(errors.ErrInsufficientFunds, "spendable balance %s is smaller than %s", balance, amt)
	}

	k.Balances[senderModule] = newBalance
	k.Balances[recipient] = k.Balances[recipient].Add(amt...)

	return nil
}
