package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/ondo/x/usdy/types"
	"github.com/noble-assets/ondo/x/usdy/types/blocklist"
)

type Keeper struct {
	cdc          codec.Codec
	logger       log.Logger
	storeService store.KVStoreService
	eventService event.Service

	Denom  string
	Schema collections.Schema

	Paused collections.Item[bool]
	Pauser collections.Item[string]

	Owner            collections.Item[string]
	PendingOwner     collections.Item[string]
	BlockedAddresses collections.Map[[]byte, bool]
}

func NewKeeper(
	cdc codec.Codec,
	logger log.Logger,
	storeService store.KVStoreService,
	eventService event.Service,
	denom string,
) *Keeper {
	builder := collections.NewSchemaBuilder(storeService)

	keeper := &Keeper{
		cdc:          cdc,
		logger:       logger,
		storeService: storeService,
		eventService: eventService,

		Denom: denom,

		Paused: collections.NewItem(builder, types.PausedKey, "paused", collections.BoolValue),
		Pauser: collections.NewItem(builder, types.PauserKey, "pauser", collections.StringValue),

		Owner:            collections.NewItem(builder, blocklist.OwnerKey, "owner", collections.StringValue),
		PendingOwner:     collections.NewItem(builder, blocklist.PendingOwnerKey, "pending_owner", collections.StringValue),
		BlockedAddresses: collections.NewMap(builder, blocklist.BlockedAddressPrefix, "blocked_addresses", collections.BytesKey, collections.BoolValue),
	}

	schema, err := builder.Build()
	if err != nil {
		panic(err)
	}

	keeper.Schema = schema
	return keeper
}

// SendRestrictionFn checks every USDY transfer on the Noble chain to and checks
// if transfers are currently paused.
func (k *Keeper) SendRestrictionFn(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error) {
	if amount := amt.AmountOf(k.Denom); !amount.IsZero() {
		paused, _ := k.Paused.Get(ctx)
		if paused {
			return toAddr, fmt.Errorf("%s transfers are paused", k.Denom)
		}

		// TODO(@john): Implement blocklist checks.
	}

	return toAddr, nil
}
