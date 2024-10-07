package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ondoprotocol/usdy-noble/v2/x/aura/types"
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

	return &types.QueryPausedResponse{Paused: k.GetPaused(ctx)}, nil
}

func (k queryServer) Owner(ctx context.Context, req *types.QueryOwner) (*types.QueryOwnerResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	return &types.QueryOwnerResponse{
		Owner:        k.GetOwner(ctx),
		PendingOwner: k.GetPendingOwner(ctx),
	}, nil
}

func (k queryServer) Burners(ctx context.Context, req *types.QueryBurners) (*types.QueryBurnersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	return &types.QueryBurnersResponse{Burners: k.GetBurners(ctx)}, nil
}

func (k queryServer) Minters(ctx context.Context, req *types.QueryMinters) (*types.QueryMintersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	return &types.QueryMintersResponse{Minters: k.GetMinters(ctx)}, nil
}

func (k queryServer) Pausers(ctx context.Context, req *types.QueryPausers) (*types.QueryPausersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	return &types.QueryPausersResponse{Pausers: k.GetPausers(ctx)}, nil
}

func (k queryServer) BlockedChannels(ctx context.Context, req *types.QueryBlockedChannels) (*types.QueryBlockedChannelsResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	return &types.QueryBlockedChannelsResponse{BlockedChannels: k.GetBlockedChannels(ctx)}, nil
}
