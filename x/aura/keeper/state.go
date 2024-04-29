package keeper

import "context"

func (k *Keeper) GetBurners(ctx context.Context) ([]string, error) {
	var burners []string

	err := k.Burners.Walk(ctx, nil, func(burner string) (stop bool, err error) {
		burners = append(burners, burner)
		return false, nil
	})

	return burners, err
}

func (k *Keeper) GetMinters(ctx context.Context) ([]string, error) {
	var minters []string

	err := k.Minters.Walk(ctx, nil, func(minter string) (stop bool, err error) {
		minters = append(minters, minter)
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
