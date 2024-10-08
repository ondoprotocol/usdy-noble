package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ondoprotocol/usdy-noble/v2/types"
	"github.com/ondoprotocol/usdy-noble/v2/types/blocklist"
)

type Keeper struct {
	Denom string

	schema       collections.Schema
	storeService store.KVStoreService
	eventService event.Service

	Paused          collections.Item[bool]
	Owner           collections.Item[string]
	PendingOwner    collections.Item[string]
	Burners         collections.Map[string, math.Int]
	Minters         collections.Map[string, math.Int]
	Pausers         collections.Map[string, []byte]
	BlockedChannels collections.Map[string, []byte]

	BlocklistOwner        collections.Item[string]
	BlocklistPendingOwner collections.Item[string]
	BlockedAddresses      collections.Map[[]byte, []byte]

	addressCodec address.Codec
	bankKeeper   types.BankKeeper
}

func NewKeeper(
	denom string,
	storeService store.KVStoreService,
	eventService event.Service,
	addressCodec address.Codec,
	bankKeeper types.BankKeeper,
) *Keeper {
	builder := collections.NewSchemaBuilder(storeService)

	keeper := &Keeper{
		Denom: denom,

		storeService: storeService,
		eventService: eventService,

		Paused:          collections.NewItem(builder, types.PausedKey, "paused", collections.BoolValue),
		Owner:           collections.NewItem(builder, types.OwnerKey, "owner", collections.StringValue),
		PendingOwner:    collections.NewItem(builder, types.PendingOwnerKey, "pending_owner", collections.StringValue),
		Burners:         collections.NewMap(builder, types.BurnerPrefix, "burners", collections.StringKey, sdk.IntValue),
		Minters:         collections.NewMap(builder, types.MinterPrefix, "minters", collections.StringKey, sdk.IntValue),
		Pausers:         collections.NewMap(builder, types.PauserPrefix, "pausers", collections.StringKey, collections.BytesValue),
		BlockedChannels: collections.NewMap(builder, types.BlockedChannelPrefix, "blocked_channels", collections.StringKey, collections.BytesValue),

		BlocklistOwner:        collections.NewItem(builder, blocklist.OwnerKey, "blocklist_owner", collections.StringValue),
		BlocklistPendingOwner: collections.NewItem(builder, blocklist.PendingOwnerKey, "blocklist_pending_owner", collections.StringValue),
		BlockedAddresses:      collections.NewMap(builder, blocklist.BlockedAddressPrefix, "blocked_address", collections.BytesKey, collections.BytesValue),

		addressCodec: addressCodec,
		bankKeeper:   bankKeeper,
	}

	schema, err := builder.Build()
	if err != nil {
		panic(err)
	}

	keeper.schema = schema
	return keeper
}

// SetBankKeeper overwrites the bank keeper used in this module.
func (k *Keeper) SetBankKeeper(bankKeeper types.BankKeeper) {
	k.bankKeeper = bankKeeper
}

// SendRestrictionFn executes necessary checks against all USDY transfers.
func (k *Keeper) SendRestrictionFn(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error) {
	if amount := amt.AmountOf(k.Denom); !amount.IsZero() {
		burning := !fromAddr.Equals(types.ModuleAddress) && toAddr.Equals(types.ModuleAddress)
		if burning {
			return toAddr, nil
		}

		if k.GetPaused(ctx) {
			return toAddr, fmt.Errorf("%s transfers are paused", k.Denom)
		}

		minting := fromAddr.Equals(types.ModuleAddress) && !toAddr.Equals(types.ModuleAddress)

		if !minting {
			if k.HasBlockedAddress(ctx, fromAddr) {
				return toAddr, fmt.Errorf("%s is blocked from sending %s", fromAddr, k.Denom)
			}
		}

		if k.HasBlockedAddress(ctx, toAddr) {
			return toAddr, fmt.Errorf("%s is blocked from receiving %s", toAddr, k.Denom)
		}

		for _, channel := range k.GetBlockedChannels(ctx) {
			escrow := transfertypes.GetEscrowAddress(transfertypes.PortID, channel)

			if toAddr.Equals(escrow) {
				return toAddr, fmt.Errorf("%s transfers are blocked on %s", k.Denom, channel)
			}
		}
	}

	return toAddr, nil
}
