package types

import (
	fmt "fmt"

	types "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/proto/tendermint/types"
)

type CapabilityI interface {
	MsgType() sdk.Msg
	Accept(msg sdk.Msg, block abci.Header) (allow bool, updated CapabilityI, delete bool)
}

// NewGrantCapability returns new CapabilityGrant
func NewGrantCapability(capability CapabilityI, expiration int64) (*CapabilityGrant, error) {
	auth := CapabilityGrant{
		Expiration: expiration,
	}
	msg, ok := capability.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("cannot proto marshal %T", capability)
	}

	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	auth.Capability = any

	return &auth, nil
}
