package keeper

import (
	"context"
	"errors"

	sdkerrors "cosmossdk.io/errors"
	"github.com/noble-assets/aura/x/aura/types/blocklist"
)

var _ blocklist.MsgServer = &blocklistMsgServer{}

type blocklistMsgServer struct {
	*Keeper
}

func NewBlocklistMsgServer(keeper *Keeper) blocklist.MsgServer {
	return &blocklistMsgServer{Keeper: keeper}
}

func (k blocklistMsgServer) TransferOwnership(ctx context.Context, msg *blocklist.MsgTransferOwnership) (*blocklist.MsgTransferOwnershipResponse, error) {
	owner, err := k.Owner.Get(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to retrieve blocklist owner from state")
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(blocklist.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	err = k.PendingOwner.Set(ctx, msg.NewOwner)
	if err != nil {
		return nil, errors.New("unable to set blocklist pending owner state")
	}

	return &blocklist.MsgTransferOwnershipResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &blocklist.OwnershipTransferStarted{
		OldOwner: owner,
		NewOwner: msg.NewOwner,
	})
}

func (k blocklistMsgServer) AcceptOwnership(ctx context.Context, msg *blocklist.MsgAcceptOwnership) (*blocklist.MsgAcceptOwnershipResponse, error) {
	pendingOwner, err := k.PendingOwner.Get(ctx)
	if err != nil {
		return nil, errors.New("there is no blocklist pending owner")
	}
	if msg.Signer != pendingOwner {
		return nil, sdkerrors.Wrapf(blocklist.ErrInvalidPendingOwner, "expected %s, got %s", pendingOwner, msg.Signer)
	}

	err = k.Owner.Set(ctx, pendingOwner)
	if err != nil {
		return nil, errors.New("unable to set blocklist owner state")
	}
	err = k.PendingOwner.Remove(ctx)
	if err != nil {
		return nil, errors.New("unable to remove blocklist pending owner state")
	}

	return &blocklist.MsgAcceptOwnershipResponse{}, nil
}

func (k blocklistMsgServer) AddToBlocklist(ctx context.Context, msg *blocklist.MsgAddToBlocklist) (*blocklist.MsgAddToBlocklistResponse, error) {
	owner, err := k.Owner.Get(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to retrieve blocklist owner from state")
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(blocklist.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	for _, account := range msg.Accounts {
		address, err := k.accountKeeper.AddressCodec().StringToBytes(account)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "unable to decode account address %s", account)
		}

		err = k.BlockedAddresses.Set(ctx, address, true)
		if err != nil {
			return nil, errors.New("unable to set blocked address state")
		}
	}

	return &blocklist.MsgAddToBlocklistResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &blocklist.BlockedAddressesAdded{
		Accounts: msg.Accounts,
	})
}

func (k blocklistMsgServer) RemoveFromBlocklist(ctx context.Context, msg *blocklist.MsgRemoveFromBlocklist) (*blocklist.MsgRemoveFromBlocklistResponse, error) {
	owner, err := k.Owner.Get(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to retrieve blocklist owner from state")
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(blocklist.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	for _, account := range msg.Accounts {
		address, err := k.accountKeeper.AddressCodec().StringToBytes(account)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "unable to decode account address %s", account)
		}

		err = k.BlockedAddresses.Remove(ctx, address)
		if err != nil {
			return nil, errors.New("unable to remove blocked address state")
		}
	}

	return &blocklist.MsgRemoveFromBlocklistResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &blocklist.BlockedAddressesRemoved{
		Accounts: msg.Accounts,
	})
}
