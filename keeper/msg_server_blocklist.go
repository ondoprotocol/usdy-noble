package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/ondoprotocol/usdy-noble/v2/types/blocklist"
)

var _ blocklist.MsgServer = &blocklistMsgServer{}

type blocklistMsgServer struct {
	*Keeper
}

func NewBlocklistMsgServer(keeper *Keeper) blocklist.MsgServer {
	return &blocklistMsgServer{Keeper: keeper}
}

func (k blocklistMsgServer) TransferOwnership(ctx context.Context, msg *blocklist.MsgTransferOwnership) (*blocklist.MsgTransferOwnershipResponse, error) {
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

	if err := k.SetBlocklistPendingOwner(ctx, msg.NewOwner); err != nil {
		return nil, err
	}

	return &blocklist.MsgTransferOwnershipResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &blocklist.OwnershipTransferStarted{
		PreviousOwner: owner,
		NewOwner:      msg.NewOwner,
	})
}

func (k blocklistMsgServer) AcceptOwnership(ctx context.Context, msg *blocklist.MsgAcceptOwnership) (*blocklist.MsgAcceptOwnershipResponse, error) {
	pendingOwner := k.GetBlocklistPendingOwner(ctx)
	if pendingOwner == "" {
		return nil, blocklist.ErrNoPendingOwner
	}
	if msg.Signer != pendingOwner {
		return nil, errors.Wrapf(blocklist.ErrInvalidPendingOwner, "expected %s, got %s", pendingOwner, msg.Signer)
	}

	owner := k.GetBlocklistOwner(ctx)
	if err := k.SetBlocklistOwner(ctx, msg.Signer); err != nil {
		return nil, err
	}
	if err := k.DeleteBlocklistPendingOwner(ctx); err != nil {
		return nil, err
	}

	return &blocklist.MsgAcceptOwnershipResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &blocklist.OwnershipTransferred{
		PreviousOwner: owner,
		NewOwner:      msg.Signer,
	})
}

func (k blocklistMsgServer) AddToBlocklist(ctx context.Context, msg *blocklist.MsgAddToBlocklist) (*blocklist.MsgAddToBlocklistResponse, error) {
	owner := k.GetBlocklistOwner(ctx)
	if owner == "" {
		return nil, blocklist.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, errors.Wrapf(blocklist.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	for _, account := range msg.Accounts {
		address, err := k.addressCodec.StringToBytes(account)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to decode account address %s", account)
		}

		if err := k.SetBlockedAddress(ctx, address); err != nil {
			return nil, err
		}
	}

	return &blocklist.MsgAddToBlocklistResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &blocklist.BlockedAddressesAdded{
		Accounts: msg.Accounts,
	})
}

func (k blocklistMsgServer) RemoveFromBlocklist(ctx context.Context, msg *blocklist.MsgRemoveFromBlocklist) (*blocklist.MsgRemoveFromBlocklistResponse, error) {
	owner := k.GetBlocklistOwner(ctx)
	if owner == "" {
		return nil, blocklist.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, errors.Wrapf(blocklist.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	for _, account := range msg.Accounts {
		address, err := k.addressCodec.StringToBytes(account)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to decode account address %s", account)
		}

		if err := k.DeleteBlockedAddress(ctx, address); err != nil {
			return nil, err
		}
	}

	return &blocklist.MsgRemoveFromBlocklistResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &blocklist.BlockedAddressesRemoved{
		Accounts: msg.Accounts,
	})
}
