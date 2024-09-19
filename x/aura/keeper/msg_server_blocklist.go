package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/usdy-noble/v2/x/aura/types/blocklist"
)

var _ blocklist.MsgServer = &blocklistMsgServer{}

type blocklistMsgServer struct {
	*Keeper
}

func NewBlocklistMsgServer(keeper *Keeper) blocklist.MsgServer {
	return &blocklistMsgServer{Keeper: keeper}
}

func (k blocklistMsgServer) TransferOwnership(goCtx context.Context, msg *blocklist.MsgTransferOwnership) (*blocklist.MsgTransferOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetBlocklistOwner(ctx)
	if owner == "" {
		return nil, blocklist.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, errors.Wrapf(blocklist.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	if msg.NewOwner == owner {
		return nil, blocklist.ErrSameOwner
	}

	k.SetBlocklistPendingOwner(ctx, msg.NewOwner)

	return &blocklist.MsgTransferOwnershipResponse{}, ctx.EventManager().EmitTypedEvent(&blocklist.OwnershipTransferStarted{
		PreviousOwner: owner,
		NewOwner:      msg.NewOwner,
	})
}

func (k blocklistMsgServer) AcceptOwnership(goCtx context.Context, msg *blocklist.MsgAcceptOwnership) (*blocklist.MsgAcceptOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pendingOwner := k.GetBlocklistPendingOwner(ctx)
	if pendingOwner == "" {
		return nil, blocklist.ErrNoPendingOwner
	}
	if msg.Signer != pendingOwner {
		return nil, errors.Wrapf(blocklist.ErrInvalidPendingOwner, "expected %s, got %s", pendingOwner, msg.Signer)
	}

	owner := k.GetBlocklistOwner(ctx)
	k.SetBlocklistOwner(ctx, msg.Signer)
	k.DeleteBlocklistPendingOwner(ctx)

	return &blocklist.MsgAcceptOwnershipResponse{}, ctx.EventManager().EmitTypedEvent(&blocklist.OwnershipTransferred{
		PreviousOwner: owner,
		NewOwner:      msg.Signer,
	})
}

func (k blocklistMsgServer) AddToBlocklist(goCtx context.Context, msg *blocklist.MsgAddToBlocklist) (*blocklist.MsgAddToBlocklistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetBlocklistOwner(ctx)
	if owner == "" {
		return nil, blocklist.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, errors.Wrapf(blocklist.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	for _, account := range msg.Accounts {
		address, err := sdk.AccAddressFromBech32(account)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to decode account address %s", account)
		}

		k.SetBlockedAddress(ctx, address)
	}

	return &blocklist.MsgAddToBlocklistResponse{}, ctx.EventManager().EmitTypedEvent(&blocklist.BlockedAddressesAdded{
		Accounts: msg.Accounts,
	})
}

func (k blocklistMsgServer) RemoveFromBlocklist(goCtx context.Context, msg *blocklist.MsgRemoveFromBlocklist) (*blocklist.MsgRemoveFromBlocklistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetBlocklistOwner(ctx)
	if owner == "" {
		return nil, blocklist.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, errors.Wrapf(blocklist.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	for _, account := range msg.Accounts {
		address, err := sdk.AccAddressFromBech32(account)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to decode account address %s", account)
		}

		k.DeleteBlockedAddress(ctx, address)
	}

	return &blocklist.MsgRemoveFromBlocklistResponse{}, ctx.EventManager().EmitTypedEvent(&blocklist.BlockedAddressesRemoved{
		Accounts: msg.Accounts,
	})
}
