package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/errors"
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
		return nil, errors.ErrInvalidRequest
	}

	owner, err := k.Keeper.BlocklistOwner.Get(ctx)
	pendingOwner, _ := k.BlocklistPendingOwner.Get(ctx)

	return &blocklist.QueryOwnerResponse{
		Owner:        owner,
		PendingOwner: pendingOwner,
	}, err
}
