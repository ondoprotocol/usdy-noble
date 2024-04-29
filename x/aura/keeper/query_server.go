package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/noble-assets/aura/x/aura/types"
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

func (k queryServer) Paused(ctx context.Context, req *types.QueryPaused) (*types.QueryPausedResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	paused, _ := k.Keeper.Paused.Get(ctx)

	return &types.QueryPausedResponse{Paused: paused}, nil
}

func (k queryServer) Burners(ctx context.Context, req *types.QueryBurners) (*types.QueryBurnersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	burners, err := k.GetBurners(ctx)

	return &types.QueryBurnersResponse{Burners: burners}, err
}

func (k queryServer) Minters(ctx context.Context, req *types.QueryMinters) (*types.QueryMintersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	minters, err := k.GetMinters(ctx)

	return &types.QueryMintersResponse{Minters: minters}, err
}

func (k queryServer) Pauser(ctx context.Context, req *types.QueryPauser) (*types.QueryPauserResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	pauser, err := k.Keeper.Pauser.Get(ctx)

	return &types.QueryPauserResponse{Pauser: pauser}, err
}
