package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/aura/x/aura/types"
	"github.com/noble-assets/aura/x/aura/types/blocklist"
)

type Keeper struct {
	cdc          codec.Codec
	logger       log.Logger
	storeService store.KVStoreService
	eventService event.Service

	Denom  string
	Schema collections.Schema

	Paused       collections.Item[bool]
	Owner        collections.Item[string]
	PendingOwner collections.Item[string]
	Burners      collections.Map[string, math.Int]
	Minters      collections.Map[string, math.Int]
	Pausers      collections.KeySet[string]

	BlocklistOwner        collections.Item[string]
	BlocklistPendingOwner collections.Item[string]
	BlockedAddresses      collections.Map[[]byte, bool]

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

func NewKeeper(
	cdc codec.Codec,
	logger log.Logger,
	storeService store.KVStoreService,
	eventService event.Service,
	denom string,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) *Keeper {
	builder := collections.NewSchemaBuilder(storeService)

	keeper := &Keeper{
		cdc:          cdc,
		logger:       logger,
		storeService: storeService,
		eventService: eventService,

		Denom:         denom,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,

		Paused: collections.NewItem(builder, types.PausedKey, "paused", collections.BoolValue),

		Owner:        collections.NewItem(builder, types.OwnerKey, "owner", collections.StringValue),
		PendingOwner: collections.NewItem(builder, types.PendingOwnerKey, "pending_owner", collections.StringValue),

		Burners: collections.NewMap(builder, types.BurnerPrefix, "burners", collections.StringKey, sdk.IntValue),
		Minters: collections.NewMap(builder, types.MinterPrefix, "minters", collections.StringKey, sdk.IntValue),
		Pausers: collections.NewKeySet(builder, types.PauserPrefix, "pausers", collections.StringKey),

		BlocklistOwner:        collections.NewItem(builder, blocklist.OwnerKey, "blocklist_owner", collections.StringValue),
		BlocklistPendingOwner: collections.NewItem(builder, blocklist.PendingOwnerKey, "blocklist_pending_owner", collections.StringValue),
		BlockedAddresses:      collections.NewMap(builder, blocklist.BlockedAddressPrefix, "blocked_addresses", collections.BytesKey, collections.BoolValue),
	}

	schema, err := builder.Build()
	if err != nil {
		panic(err)
	}

	keeper.Schema = schema
	return keeper
}

// SendRestrictionFn executes necessary checks against all USDY transfers.
func (k *Keeper) SendRestrictionFn(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error) {
	if amount := amt.AmountOf(k.Denom); !amount.IsZero() {
		burning := !fromAddr.Equals(types.ModuleAddress) && toAddr.Equals(types.ModuleAddress)
		if burning {
			return toAddr, nil
		}

		paused, _ := k.Paused.Get(ctx)
		if paused {
			return toAddr, fmt.Errorf("%s transfers are paused", k.Denom)
		}

		minting := fromAddr.Equals(types.ModuleAddress) && !toAddr.Equals(types.ModuleAddress)

		if !minting {
			blocked, _ := k.BlockedAddresses.Get(ctx, fromAddr)
			if blocked {
				address, _ := k.accountKeeper.AddressCodec().BytesToString(fromAddr)
				return toAddr, fmt.Errorf("%s is blocked from sending %s", address, k.Denom)
			}
		}

		blocked, _ := k.BlockedAddresses.Get(ctx, toAddr)
		if blocked {
			address, _ := k.accountKeeper.AddressCodec().BytesToString(toAddr)
			return toAddr, fmt.Errorf("%s is blocked from receiving %s", address, k.Denom)
		}
	}

	return toAddr, nil
}
