package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// GenericAuthorization grants the permission to execute any transaction of the provided
// sdk.Msg type without restrictions
type GenericAuthorization struct {
	// MsgType is the type of Msg this capability grant allows
	Msg sdk.Msg
}

func (cap GenericAuthorization) MsgType() sdk.Msg {
	return cap.MsgType()
}

func (cap GenericAuthorization) Accept(msg sdk.Msg, block tmproto.Header) (allow bool, updated *codectypes.Any, delete bool) {
	genAuth, err := ConvertToAny(cap)
	if err != nil {
		return false, nil, false
	}
	return true, genAuth, false
}
