package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ondoprotocol/aura/x/aura/types/blocklist"
)

var _ blocklist.QueryServer = &blocklistQueryServer{}

type blocklistQueryServer struct {
	*Keeper
}

func NewBlocklistQueryServer(keeper *Keeper) blocklist.QueryServer {
	return &blocklistQueryServer{Keeper: keeper}
}

func (k blocklistQueryServer) Owner(goCtx context.Context, req *blocklist.QueryOwner) (*blocklist.QueryOwnerResponse, error) {
	if req == nil {
		return nil, errorstypes.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &blocklist.QueryOwnerResponse{
		Owner:        k.GetBlocklistOwner(ctx),
		PendingOwner: k.GetBlocklistPendingOwner(ctx),
	}, nil
}

func (k blocklistQueryServer) Addresses(goCtx context.Context, req *blocklist.QueryAddresses) (*blocklist.QueryAddressesResponse, error) {
	if req == nil {
		return nil, errorstypes.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), blocklist.BlockedAddressPrefix)

	var addresses []string
	pagination, err := query.Paginate(store, req.Pagination, func(key []byte, _ []byte) error {
		addresses = append(addresses, sdk.AccAddress(key).String())
		return nil
	})

	return &blocklist.QueryAddressesResponse{
		Addresses:  addresses,
		Pagination: pagination,
	}, err
}

func (k blocklistQueryServer) Address(goCtx context.Context, req *blocklist.QueryAddress) (*blocklist.QueryAddressResponse, error) {
	if req == nil {
		return nil, errorstypes.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to decode address %s", req.Address)
	}

	blocked := k.HasBlockedAddress(ctx, address)
	return &blocklist.QueryAddressResponse{Blocked: blocked}, nil
}
