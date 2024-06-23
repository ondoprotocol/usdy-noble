package aura

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	"github.com/noble-assets/aura/x/aura/keeper"
)

type Decorator struct {
	*keeper.Keeper
}

var _ sdk.AnteDecorator = Decorator{}

func NewAnteDecorator(keeper *keeper.Keeper) Decorator {
	return Decorator{Keeper: keeper}
}

func (d Decorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	msgs := tx.GetMsgs()

	err = d.CheckMessages(ctx, msgs)
	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func (d Decorator) CheckMessages(ctx sdk.Context, msgs []sdk.Msg) error {
	for _, raw := range msgs {
		if msg, ok := raw.(*authz.MsgExec); ok {
			nestedMsgs, err := msg.GetMessages()
			if err != nil {
				return err
			}

			return d.CheckMessages(ctx, nestedMsgs)
		}

		switch msg := raw.(type) {
		case *transfertypes.MsgTransfer:
			if msg.Token.Denom == d.Denom {
				if !d.HasChannel(ctx, msg.SourceChannel) {
					return fmt.Errorf("%s cannot be transferred over %s", d.Denom, msg.SourceChannel)
				}
			}
		}
	}

	return nil
}
