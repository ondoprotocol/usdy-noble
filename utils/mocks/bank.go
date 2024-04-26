package mocks

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/ondo/x/usdy/types"
)

var _ types.BankKeeper = BankKeeper{}

type BankKeeper struct{}

func (k BankKeeper) BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	return nil
}

func (k BankKeeper) MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	return nil
}

func (k BankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return nil
}

func (k BankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	return nil
}
