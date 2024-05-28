package keeper_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/noble-assets/aura/x/aura/types/bridge"
	"github.com/stretchr/testify/require"
)

func TestBridgePayloadEncoding(t *testing.T) {
	// https://etherscan.io/tx/0x781e74c855ba859fd7c963f75227c57daa2ae1dacd28c5fe1a1e5658a7fd8c32

	amount, _ := sdk.NewIntFromString("2968688625494144582144")
	payload := bridge.Payload{
		Version: "1.0",
		Chain:   1,
		Sender:  common.HexToAddress("0x90e0d37f59B4d3202880d2FB17f3e50b7056f762"),
		Amount:  amount,
		Nonce:   42,
	}

	bz, err := payload.Bytes()

	require.NoError(t, err)
	require.Equal(t, "312E300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000090E0D37F59B4D3202880D2FB17F3E50B7056F7620000000000000000000000000000000000000000000000A0EED4B019B9240200000000000000000000000000000000000000000000000000000000000000002A", strings.ToUpper(common.Bytes2Hex(bz)))
}
