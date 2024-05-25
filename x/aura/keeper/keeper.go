package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/aura/x/aura/types"
)

type Keeper struct {
	cdc      codec.Codec
	storeKey storetypes.StoreKey

	Denom      string
	bankKeeper types.BankKeeper
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	denom string,
	bankKeeper types.BankKeeper,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,

		Denom:      denom,
		bankKeeper: bankKeeper,
	}
}

// SetBankKeeper overwrites the bank keeper used in this module.
func (k *Keeper) SetBankKeeper(bankKeeper types.BankKeeper) {
	k.bankKeeper = bankKeeper
}

// SendRestrictionFn executes necessary checks against all USDY transfers.
func (k *Keeper) SendRestrictionFn(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error) {
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
	}

	return toAddr, nil
}
