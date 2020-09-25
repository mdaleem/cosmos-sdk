package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterLegacyAminoCodec registers all the necessary types and interfaces for the
// governance module.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgGrantAuthorization{}, "cosmos-sdk/GrantAuthorization", nil)
	cdc.RegisterConcrete(MsgRevokeAuthorization{}, "cosmos-sdk/RevokeAuthorization", nil)
	cdc.RegisterConcrete(MsgExecDelegated{}, "cosmos-sdk/ExecDelegated", nil)
	cdc.RegisterConcrete(SendAuthorization{}, "cosmos-sdk/SendAuthorization", nil)
	cdc.RegisterConcrete(GenericAuthorization{}, "cosmos-sdk/GenericAuthorization", nil)

	cdc.RegisterInterface((*AuthorizationI)(nil), nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgGrantAuthorization{},
		&MsgRevokeAuthorization{},
		&MsgExecDelegated{},
	)
	registry.RegisterInterface(
		"cosmos.msg_authorization.v1beta1.AuthorizationI",
		(*AuthorizationI)(nil),
		&SendAuthorization{},
	)
}

// RegisterProposalTypeCodec registers an external proposal content type defined
// in another module for the internal ModuleCdc. This allows the MsgSubmitProposal
// to be correctly Amino encoded and decoded.
//
// NOTE: This should only be used for applications that are still using a concrete
// Amino codec for serialization.
func RegisterProposalTypeCodec(o interface{}, name string) {
	amino.RegisterConcrete(o, name, nil)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/gov module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/gov and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
}
