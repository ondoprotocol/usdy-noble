package keeper

import (
	"context"
	"errors"

	sdkerrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/aura/x/aura/types"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) Burn(ctx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	burner, err := k.Burner.Get(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to retrieve burner from state")
	}
	if msg.Signer != burner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidBurner, "expected %s, got %s", burner, msg.Signer)
	}

	from, err := k.accountKeeper.AddressCodec().StringToBytes(msg.From)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to decode account address %s", msg.From)
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.Denom, msg.Amount))
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, coins)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to transfer from user to module")
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to burn from module")
	}

	// TODO(@john): Do we emit an event here?
	return &types.MsgBurnResponse{}, nil
}

func (k msgServer) Mint(ctx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	minter, err := k.Minter.Get(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to retrieve minter from state")
	}
	if msg.Signer != minter {
		return nil, sdkerrors.Wrapf(types.ErrInvalidMinter, "expected %s, got %s", minter, msg.Signer)
	}

	to, err := k.accountKeeper.AddressCodec().StringToBytes(msg.To)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to decode account address %s", msg.To)
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.Denom, msg.Amount))
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to mint to module")
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, coins)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to transfer from module to user")
	}

	// TODO(@john): Do we emit an event here?
	return &types.MsgMintResponse{}, nil
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
