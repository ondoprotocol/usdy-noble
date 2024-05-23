package keeper

import (
	"context"

	"cosmossdk.io/errors"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/noble-assets/aura/x/aura/types/blocklist"
)

var _ blocklist.QueryServer = &blocklistQueryServer{}

type blocklistQueryServer struct {
	*Keeper
}

func NewBlocklistQueryServer(keeper *Keeper) blocklist.QueryServer {
	return &blocklistQueryServer{Keeper: keeper}
}

func (k blocklistQueryServer) Owner(ctx context.Context, req *blocklist.QueryOwner) (*blocklist.QueryOwnerResponse, error) {
	if req == nil {
		return nil, errorstypes.ErrInvalidRequest
	}

	owner, err := k.BlocklistOwner.Get(ctx)
	pendingOwner, _ := k.BlocklistPendingOwner.Get(ctx)

	return &blocklist.QueryOwnerResponse{
		Owner:        owner,
		PendingOwner: pendingOwner,
	}, err
}

func (k blocklistQueryServer) Addresses(ctx context.Context, req *blocklist.QueryAddresses) (*blocklist.QueryAddressesResponse, error) {
	if req == nil {
		return nil, errorstypes.ErrInvalidRequest
	}

	addresses, pagination, err := query.CollectionPaginate(
		ctx, k.BlockedAddresses, req.Pagination,
		func(account []byte, blocked bool) (string, error) {
			return k.accountKeeper.AddressCodec().BytesToString(account)
		},
	)

	return &blocklist.QueryAddressesResponse{
		Addresses:  addresses,
		Pagination: pagination,
	}, err
}

func (k blocklistQueryServer) Address(ctx context.Context, req *blocklist.QueryAddress) (*blocklist.QueryAddressResponse, error) {
	if req == nil {
		return nil, errorstypes.ErrInvalidRequest
	}

	address, err := k.accountKeeper.AddressCodec().StringToBytes(req.Address)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to decode address %s", req.Address)
	}

	blocked, err := k.BlockedAddresses.Has(ctx, address)
	return &blocklist.QueryAddressResponse{Blocked: blocked}, err
}
