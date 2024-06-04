package keeper

import (
	"context"
	"encoding/json"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	"github.com/noble-assets/aura/x/aura/types"
	"github.com/noble-assets/aura/x/aura/types/bridge"
	"github.com/noble-assets/aura/x/aura/types/bridge/source"
)

var _ source.MsgServer = &sourceMsgServer{}

type sourceMsgServer struct {
	*Keeper
}

func NewSourceMsgServer(keeper *Keeper) source.MsgServer {
	return &sourceMsgServer{Keeper: keeper}
}

func (k sourceMsgServer) BurnAndCallAxelar(goCtx context.Context, msg *source.MsgBurnAndCallAxelar) (*source.MsgBurnAndCallAxelarResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.GetBridgeSourcePaused(ctx) {
		return nil, errors.New("source bridge is currently paused")
	}
	destination := k.GetBridgeDestination(ctx, msg.Destination)
	if destination == "" {
		return nil, errors.New("unknown bridge destination")
	}

	server := NewMsgServer(k.Keeper)
	_, err := server.Burn(goCtx, &types.MsgBurn{
		Signer: source.SubmoduleAddress.String(),
		From:   msg.Signer,
		Amount: msg.Amount,
	})
	if err != nil {
		return nil, err
	}

	payload := bridge.Payload{
		Version: bridge.PayloadVersion,
		Amount:  msg.Amount,
		Nonce:   k.IncrementBridgeSourceNonce(ctx),
	}
	bz, err := payload.Bytes()
	if err != nil {
		return nil, err
	}

	memo, err := json.Marshal(source.Message{
		DestinationChain:   msg.Destination,
		DestinationAddress: destination,
		Payload:            bz,
		Type:               1,   // TODO: Abstract into a constant.
		Fee:                nil, // TODO: Investigate fees.
	})
	if err != nil {
		return nil, err
	}

	timeout := uint64(ctx.BlockTime().UnixNano()) + transfertypes.DefaultRelativePacketTimeoutTimestamp
	_, err = k.transferKeeper.Transfer(goCtx, &transfertypes.MsgTransfer{
		SourcePort:       transfertypes.PortID,
		SourceChannel:    k.GetBridgeChannel(ctx),
		Token:            sdk.Coin{}, // TODO: Investigate fees.
		Sender:           msg.Signer,
		Receiver:         bridge.AxelarGMP,
		TimeoutHeight:    clienttypes.ZeroHeight(),
		TimeoutTimestamp: timeout,
		Memo:             string(memo),
	})
	if err != nil {
		return nil, err
	}

	return &source.MsgBurnAndCallAxelarResponse{}, ctx.EventManager().EmitTypedEvent(&source.BridgeInitiated{
		User:    msg.Signer,
		Nonce:   payload.Nonce,
		Chain:   payload.Chain,
		Version: payload.Version,
		Amount:  payload.Amount,
	})
}

func (k sourceMsgServer) Pause(goCtx context.Context, msg *source.MsgPause) (*source.MsgPauseResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (k sourceMsgServer) Unpause(goCtx context.Context, msg *source.MsgUnpause) (*source.MsgUnpauseResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (k sourceMsgServer) TransferOwnership(goCtx context.Context, msg *source.MsgTransferOwnership) (*source.MsgTransferOwnershipResponse, error) {
	// TODO implement me
	panic("implement me")
}
