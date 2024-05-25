package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type Account struct {
	Address string
	Bytes   []byte
}

func TestAccount() Account {
	bytes := secp256k1.GenPrivKey().PubKey().Address().Bytes()
	address, _ := sdk.Bech32ifyAddressBytes("noble", bytes)

	return Account{
		Address: address,
		Bytes:   bytes,
	}
}
