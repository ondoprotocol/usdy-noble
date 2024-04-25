package utils

import (
	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cosmos/cosmos-sdk/x/auth/codec"
)

var cdc = codec.NewBech32Codec("noble")

type Account struct {
	Address string
	Bytes   []byte
}

func TestAccount() Account {
	bytes := secp256k1.GenPrivKey().PubKey().Address().Bytes()
	address, _ := cdc.BytesToString(bytes)

	return Account{
		Address: address,
		Bytes:   bytes,
	}
}
