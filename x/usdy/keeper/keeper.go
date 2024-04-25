package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Keeper struct {
	cdc          codec.Codec
	logger       log.Logger
	storeService store.KVStoreService

	Schema collections.Schema
}

func NewKeeper(
	cdc codec.Codec,
	logger log.Logger,
	storeService store.KVStoreService,
) *Keeper {
	builder := collections.NewSchemaBuilder(storeService)

	keeper := &Keeper{
		cdc:          cdc,
		logger:       logger,
		storeService: storeService,
	}

	schema, err := builder.Build()
	if err != nil {
		panic(err)
	}

	keeper.Schema = schema
	return keeper
}
