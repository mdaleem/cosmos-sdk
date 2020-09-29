package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func (authorization SendCapability) MsgType() sdk.Msg {
	return &banktypes.MsgSend{}
}

func (authorization SendCapability) Accept(msg sdk.Msg, block tmproto.Header) (allow bool, updated CapabilityI, delete bool) {
	switch msg := msg.(type) {
	case *banktypes.MsgSend:
		limitLeft, isNegative := authorization.SpendLimit.SafeSub(msg.Amount)
		if isNegative {
			return false, nil, false
		}
		if limitLeft.IsZero() {
			return true, nil, true
		}
		return true, &SendCapability{SpendLimit: limitLeft}, false
	}
	return false, nil, false
}
