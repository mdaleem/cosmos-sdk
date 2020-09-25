package types

import (
	fmt "fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	proto "github.com/gogo/protobuf/proto"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func (authorization SendAuthorization) MsgType() sdk.Msg {
	return &banktypes.MsgSend{}
}

func (authorization SendAuthorization) Accept(msg sdk.Msg, block tmproto.Header) (allow bool, updated *codectypes.Any, delete bool) {
	switch msg := msg.(type) {
	case *banktypes.MsgSend:
		limitLeft, isNegative := authorization.SpendLimit.SafeSub(msg.Amount)
		if isNegative {
			return false, nil, false
		}
		if limitLeft.IsZero() {
			return true, nil, true
		}
		authorization, err := ConvertToAny(&SendAuthorization{SpendLimit: limitLeft})
		if err != nil {
			return false, nil, false
		}
		return true, authorization, false
	}
	return false, nil, false
}

// ConvertToAny converts interface(types.Authorization) to any
func ConvertToAny(authorization AuthorizationI) (*codectypes.Any, error) {
	msg, ok := authorization.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("can't proto marshal %T", msg)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return any, nil
}
