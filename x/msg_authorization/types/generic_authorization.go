package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

var (
	_ CapabilityI = &GenericCapability{}
)

func (cap GenericCapability) MsgType() sdk.Msg {
	var msg sdk.Msg
	ModuleCdc.UnpackAny(cap.Msg, msg)
	return msg
}

func (cap GenericCapability) Accept(msg sdk.Msg, block tmproto.Header) (allow bool, updated CapabilityI, delete bool) {
	return true, &cap, false
}
