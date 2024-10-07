package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/core/address"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ondoprotocol/usdy-noble/v2/x/aura/types"
)

type Keeper struct {
	Denom string

	storeService store.KVStoreService
	eventService event.Service

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
	return &Keeper{
		Denom: denom,

		storeService: storeService,
		eventService: eventService,

		addressCodec: addressCodec,
		bankKeeper:   bankKeeper,
	}
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
