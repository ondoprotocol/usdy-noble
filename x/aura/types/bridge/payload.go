package bridge

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var ABI abi.Arguments

type Payload struct {
	Version string
	Chain   uint64
	// NOTE: Current EVM implementation will need to be adjusted to accept a
	//  recipient address, which will be included here as bytes32.
	Sender common.Address
	Amount sdk.Int
	Nonce  uint64
}

func (payload *Payload) Parse(input string) error {
	bz := common.FromHex(input)
	raw, err := ABI.Unpack(bz)
	if err != nil {
		return err
	}

	version, ok := raw[0].([32]uint8)
	if !ok {
		return errors.New("invalid version")
	}
	payload.Version = string(common.TrimRightZeroes(version[:]))

	chain, ok := raw[1].(*big.Int)
	if !ok {
		return errors.New("invalid chain")
	}
	payload.Chain = chain.Uint64()

	sender, ok := raw[2].(common.Address)
	if !ok {
		return errors.New("invalid sender")
	}
	payload.Sender = sender

	amount, ok := raw[3].(*big.Int)
	if !ok {
		return errors.New("invalid amount")
	}
	payload.Amount = sdk.NewIntFromBigInt(amount)

	nonce, ok := raw[4].(*big.Int)
	if !ok {
		return errors.New("invalid nonce")
	}
	payload.Nonce = nonce.Uint64()

	return nil
}

func (payload *Payload) Bytes() ([]byte, error) {
	return ABI.Pack(
		[32]uint8(common.RightPadBytes([]byte(payload.Version), 32)),
		big.NewInt(int64(payload.Chain)),
		payload.Sender,
		payload.Amount.BigInt(),
		big.NewInt(int64(payload.Nonce)),
	)
}

//

func init() {
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)
	addressType, _ := abi.NewType("address", "", nil)

	ABI = abi.Arguments{
		{Type: bytes32Type},
		{Type: uint256Type},
		{Type: addressType},
		{Type: uint256Type},
		{Type: uint256Type},
	}
}
