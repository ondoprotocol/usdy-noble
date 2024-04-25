package keeper

import (
	"context"
	"errors"

	sdkerrors "cosmossdk.io/errors"
	"github.com/noble-assets/ondo/x/usdy/types"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) Pause(ctx context.Context, msg *types.MsgPause) (*types.MsgPauseResponse, error) {
	pauser, err := k.Pauser.Get(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to retrieve pauser from state")
	}
	if msg.Signer != pauser {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPauser, "expected %s, got %s", pauser, msg.Signer)
	}

	if paused, _ := k.Paused.Get(ctx); paused {
		return nil, errors.New("module is already paused")
	}

	err = k.Paused.Set(ctx, true)
	if err != nil {
		return nil, errors.New("unable to set paused state")
	}

	return &types.MsgPauseResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &types.Paused{
		Account: msg.Signer,
	})
}

func (k msgServer) Unpause(ctx context.Context, msg *types.MsgUnpause) (*types.MsgUnpauseResponse, error) {
	pauser, err := k.Pauser.Get(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to retrieve pauser from state")
	}
	if msg.Signer != pauser {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPauser, "expected %s, got %s", pauser, msg.Signer)
	}

	if paused, _ := k.Paused.Get(ctx); !paused {
		return nil, errors.New("module is already unpaused")
	}

	err = k.Paused.Set(ctx, false)
	if err != nil {
		return nil, errors.New("unable to set paused state")
	}

	return &types.MsgUnpauseResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &types.Unpaused{
		Account: msg.Signer,
	})
}
