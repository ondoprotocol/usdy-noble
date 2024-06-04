package aura

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v4/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v4/modules/core/exported"
	"github.com/noble-assets/aura/x/aura/keeper"
	"github.com/noble-assets/aura/x/aura/types/bridge"
	"github.com/noble-assets/aura/x/aura/types/bridge/destination"
)

var _ porttypes.IBCModule = Middleware{}

type Middleware struct {
	app    porttypes.IBCModule
	keeper keeper.Keeper
}

func (m Middleware) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version string) (string, error) {
	return m.app.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)
}

func (m Middleware) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, counterpartyVersion string) (version string, err error) {
	return m.app.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, counterpartyVersion)
}

func (m Middleware) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	return m.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

func (m Middleware) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return m.app.OnChanOpenConfirm(ctx, portID, channelID)
}

func (m Middleware) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return m.app.OnChanCloseInit(ctx, portID, channelID)
}

func (m Middleware) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return m.app.OnChanCloseConfirm(ctx, portID, channelID)
}

func (m Middleware) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) exported.Acknowledgement {
	ack := m.app.OnRecvPacket(ctx, packet, relayer)
	if !ack.Success() {
		return ack
	}

	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return channeltypes.NewErrorAcknowledgement(fmt.Errorf("cannot unmarshal ICS-20 transfer packet data"))
	}

	if packet.DestinationChannel != m.keeper.GetBridgeChannel(ctx) || data.Sender != bridge.AxelarGMP {
		return ack
	}

	var msg destination.Message
	if err := json.Unmarshal([]byte(data.GetMemo()), &msg); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	// TODO: Hand off to destination bridge.
	return ack
}

func (m Middleware) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	return m.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

func (m Middleware) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	return m.app.OnTimeoutPacket(ctx, packet, relayer)
}
