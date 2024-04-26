package keeper

import "github.com/noble-assets/ondo/x/usdy/types/blocklist"

var _ blocklist.MsgServer = &blocklistMsgServer{}

type blocklistMsgServer struct {
	*Keeper
}

func NewBlocklistMsgServer(keeper *Keeper) blocklist.MsgServer {
	return &blocklistMsgServer{Keeper: keeper}
}
