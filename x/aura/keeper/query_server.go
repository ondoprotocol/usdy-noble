package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ondoprotocol/aura/x/aura/types"
)

var _ types.QueryServer = &queryServer{}

type queryServer struct {
	*Keeper
}

func NewQueryServer(keeper *Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

func (k queryServer) Denom(_ context.Context, req *types.QueryDenom) (*types.QueryDenomResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	return &types.QueryDenomResponse{Denom: k.Keeper.Denom}, nil
}

func (k queryServer) Paused(goCtx context.Context, req *types.QueryPaused) (*types.QueryPausedResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryPausedResponse{Paused: k.GetPaused(ctx)}, nil
}

func (k queryServer) Owner(goCtx context.Context, req *types.QueryOwner) (*types.QueryOwnerResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryOwnerResponse{
		Owner:        k.GetOwner(ctx),
		PendingOwner: k.GetPendingOwner(ctx),
	}, nil
}

func (k queryServer) Burners(goCtx context.Context, req *types.QueryBurners) (*types.QueryBurnersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryBurnersResponse{Burners: k.GetBurners(ctx)}, nil
}

func (k queryServer) Minters(goCtx context.Context, req *types.QueryMinters) (*types.QueryMintersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryMintersResponse{Minters: k.GetMinters(ctx)}, nil
}

func (k queryServer) Pausers(goCtx context.Context, req *types.QueryPausers) (*types.QueryPausersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryPausersResponse{Pausers: k.GetPausers(ctx)}, nil
}
