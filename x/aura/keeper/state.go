package keeper

import (
	"context"

	"cosmossdk.io/math"
	"github.com/noble-assets/aura/x/aura/types"
)

func (k *Keeper) GetBurners(ctx context.Context) ([]types.Burner, error) {
	var burners []types.Burner

	err := k.Burners.Walk(ctx, nil, func(address string, allowance math.Int) (stop bool, err error) {
		burners = append(burners, types.Burner{
			Address:   address,
			Allowance: allowance,
		})

		return false, nil
	})

	return burners, err
}

func (k *Keeper) GetMinters(ctx context.Context) ([]types.Minter, error) {
	var minters []types.Minter

	err := k.Minters.Walk(ctx, nil, func(address string, allowance math.Int) (stop bool, err error) {
		minters = append(minters, types.Minter{
			Address:   address,
			Allowance: allowance,
		})

		return false, nil
	})

	return minters, err
}

func (k *Keeper) GetPausers(ctx context.Context) ([]string, error) {
	var pausers []string

	err := k.Pausers.Walk(ctx, nil, func(pauser string) (stop bool, err error) {
		pausers = append(pausers, pauser)
		return false, nil
	})

	return pausers, err
}
