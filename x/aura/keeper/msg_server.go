package keeper

import (
	"context"
	"errors"
	"fmt"

	sdkerrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/aura/x/aura/types"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasBurner(ctx, msg.Signer) {
		return nil, types.ErrInvalidBurner
	}
	allowance := k.GetBurner(ctx, msg.Signer)
	if allowance.LT(msg.Amount) {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientAllowance, "burner %s has an allowance of %s", msg.Signer, allowance.String())
	}

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to decode account address %s", msg.From)
	}

	if !msg.Amount.IsPositive() {
		return nil, errors.New("amount must be positive")
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

	k.SetBurner(ctx, msg.Signer, allowance.Sub(msg.Amount))

	// NOTE: The bank module emits an event for us.
	return &types.MsgBurnResponse{}, nil
}

func (k msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasMinter(ctx, msg.Signer) {
		return nil, types.ErrInvalidMinter
	}
	allowance := k.GetMinter(ctx, msg.Signer)
	if allowance.LT(msg.Amount) {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientAllowance, "minter %s has an allowance of %s", msg.Signer, allowance.String())
	}

	to, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to decode account address %s", msg.To)
	}

	if !msg.Amount.IsPositive() {
		return nil, errors.New("amount must be positive")
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

	k.SetMinter(ctx, msg.Signer, allowance.Sub(msg.Amount))

	// NOTE: The bank module emits an event for us.
	return &types.MsgMintResponse{}, nil
}

func (k msgServer) Pause(goCtx context.Context, msg *types.MsgPause) (*types.MsgPauseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasPauser(ctx, msg.Signer) {
		return nil, types.ErrInvalidPauser
	}
	if k.GetPaused(ctx) {
		return nil, errors.New("module is already paused")
	}

	k.SetPaused(ctx, true)

	return &types.MsgPauseResponse{}, ctx.EventManager().EmitTypedEvent(&types.Paused{
		Account: msg.Signer,
	})
}

func (k msgServer) Unpause(goCtx context.Context, msg *types.MsgUnpause) (*types.MsgUnpauseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	if !k.GetPaused(ctx) {
		return nil, errors.New("module is already unpaused")
	}

	k.SetPaused(ctx, false)

	return &types.MsgUnpauseResponse{}, ctx.EventManager().EmitTypedEvent(&types.Unpaused{
		Account: msg.Signer,
	})
}

func (k msgServer) TransferOwnership(goCtx context.Context, msg *types.MsgTransferOwnership) (*types.MsgTransferOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	if msg.NewOwner == owner {
		return nil, types.ErrSameOwner
	}

	k.SetPendingOwner(ctx, msg.NewOwner)

	return &types.MsgTransferOwnershipResponse{}, ctx.EventManager().EmitTypedEvent(&types.OwnershipTransferStarted{
		PreviousOwner: owner,
		NewOwner:      msg.NewOwner,
	})
}

func (k msgServer) AcceptOwnership(goCtx context.Context, msg *types.MsgAcceptOwnership) (*types.MsgAcceptOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pendingOwner := k.GetPendingOwner(ctx)
	if pendingOwner == "" {
		return nil, types.ErrNoPendingOwner
	}
	if msg.Signer != pendingOwner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPendingOwner, "expected %s, got %s", pendingOwner, msg.Signer)
	}

	owner := k.GetOwner(ctx)

	k.SetOwner(ctx, msg.Signer)
	k.DeletePendingOwner(ctx)

	return &types.MsgAcceptOwnershipResponse{}, ctx.EventManager().EmitTypedEvent(&types.OwnershipTransferred{
		PreviousOwner: owner,
		NewOwner:      msg.Signer,
	})
}

func (k msgServer) AddBurner(goCtx context.Context, msg *types.MsgAddBurner) (*types.MsgAddBurnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	if k.HasBurner(ctx, msg.Burner) {
		return nil, fmt.Errorf("%s is already a burner", msg.Burner)
	}

	if msg.Allowance.IsNegative() {
		return nil, errors.New("allowance cannot be negative")
	}

	k.SetBurner(ctx, msg.Burner, msg.Allowance)

	return &types.MsgAddBurnerResponse{}, ctx.EventManager().EmitTypedEvent(&types.BurnerAdded{
		Address:   msg.Burner,
		Allowance: msg.Allowance,
	})
}

func (k msgServer) RemoveBurner(goCtx context.Context, msg *types.MsgRemoveBurner) (*types.MsgRemoveBurnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	if !k.HasBurner(ctx, msg.Burner) {
		return nil, fmt.Errorf("%s is not a burner", msg.Burner)
	}

	k.DeleteBurner(ctx, msg.Burner)

	return &types.MsgRemoveBurnerResponse{}, ctx.EventManager().EmitTypedEvent(&types.BurnerRemoved{
		Address: msg.Burner,
	})
}

func (k msgServer) SetBurnerAllowance(goCtx context.Context, msg *types.MsgSetBurnerAllowance) (*types.MsgSetBurnerAllowanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	if !k.HasBurner(ctx, msg.Burner) {
		return nil, fmt.Errorf("%s is not a burner", msg.Burner)
	}

	if msg.Allowance.IsNegative() {
		return nil, errors.New("allowance cannot be negative")
	}

	allowance := k.GetBurner(ctx, msg.Burner)
	k.SetBurner(ctx, msg.Burner, msg.Allowance)

	return &types.MsgSetBurnerAllowanceResponse{}, ctx.EventManager().EmitTypedEvent(&types.BurnerUpdated{
		Address:           msg.Burner,
		PreviousAllowance: allowance,
		NewAllowance:      msg.Allowance,
	})
}

func (k msgServer) AddMinter(goCtx context.Context, msg *types.MsgAddMinter) (*types.MsgAddMinterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	if k.HasMinter(ctx, msg.Minter) {
		return nil, fmt.Errorf("%s is already a minter", msg.Minter)
	}

	if msg.Allowance.IsNegative() {
		return nil, errors.New("allowance cannot be negative")
	}

	k.SetMinter(ctx, msg.Minter, msg.Allowance)

	return &types.MsgAddMinterResponse{}, ctx.EventManager().EmitTypedEvent(&types.MinterAdded{
		Address:   msg.Minter,
		Allowance: msg.Allowance,
	})
}

func (k msgServer) RemoveMinter(goCtx context.Context, msg *types.MsgRemoveMinter) (*types.MsgRemoveMinterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	if !k.HasMinter(ctx, msg.Minter) {
		return nil, fmt.Errorf("%s is not a minter", msg.Minter)
	}

	k.DeleteMinter(ctx, msg.Minter)

	return &types.MsgRemoveMinterResponse{}, ctx.EventManager().EmitTypedEvent(&types.MinterRemoved{
		Address: msg.Minter,
	})
}

func (k msgServer) SetMinterAllowance(goCtx context.Context, msg *types.MsgSetMinterAllowance) (*types.MsgSetMinterAllowanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	if !k.HasMinter(ctx, msg.Minter) {
		return nil, fmt.Errorf("%s is not a minter", msg.Minter)
	}

	if msg.Allowance.IsNegative() {
		return nil, errors.New("allowance cannot be negative")
	}

	allowance := k.GetMinter(ctx, msg.Minter)
	k.SetMinter(ctx, msg.Minter, msg.Allowance)

	return &types.MsgSetMinterAllowanceResponse{}, ctx.EventManager().EmitTypedEvent(&types.MinterUpdated{
		Address:           msg.Minter,
		PreviousAllowance: allowance,
		NewAllowance:      msg.Allowance,
	})
}

func (k msgServer) AddPauser(goCtx context.Context, msg *types.MsgAddPauser) (*types.MsgAddPauserResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	if k.HasPauser(ctx, msg.Pauser) {
		return nil, fmt.Errorf("%s is already a pauser", msg.Pauser)
	}

	k.SetPauser(ctx, msg.Pauser)

	return &types.MsgAddPauserResponse{}, ctx.EventManager().EmitTypedEvent(&types.PauserAdded{
		Address: msg.Pauser,
	})
}

func (k msgServer) RemovePauser(goCtx context.Context, msg *types.MsgRemovePauser) (*types.MsgRemovePauserResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	if !k.HasPauser(ctx, msg.Pauser) {
		return nil, fmt.Errorf("%s is not a pauser", msg.Pauser)
	}

	k.DeletePauser(ctx, msg.Pauser)

	return &types.MsgRemovePauserResponse{}, ctx.EventManager().EmitTypedEvent(&types.PauserRemoved{
		Address: msg.Pauser,
	})
}
