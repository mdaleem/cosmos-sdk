package types

import (
	types "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/proto/tendermint/types"
)

type AuthorizationI interface {
	MsgType() sdk.Msg
	Accept(msg sdk.Msg, block abci.Header) (allow bool, updated *types.Any, delete bool)
}
