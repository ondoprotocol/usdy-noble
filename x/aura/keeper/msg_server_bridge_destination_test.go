package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/noble-assets/aura/x/aura/types/bridge"
	"github.com/stretchr/testify/require"
)

func TestBridgePayloadDecoding(t *testing.T) {
	// https://etherscan.io/tx/0xf2372565408ed06136c6094007bafaa40023dc4668fe3bb8fef52e9aaea0bdd2

	amount, _ := sdk.NewIntFromString("12300000000000000")
	expected := bridge.Payload{
		Version: "1.0",
		Chain:   5000,
		Sender:  common.HexToAddress("0x26621f75cECaD3501202961E81f74B648F9DCe80"),
		Amount:  amount,
		Nonce:   0,
	}

	var payload bridge.Payload
	err := payload.Parse("0x312e300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000138800000000000000000000000026621f75cecad3501202961e81f74b648f9dce80000000000000000000000000000000000000000000000000002bb2c8eabcc0000000000000000000000000000000000000000000000000000000000000000000")

	require.NoError(t, err)
	require.EqualValues(t, expected, payload)
}
