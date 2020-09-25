package types

import (
	fmt "fmt"
	"time"

	types "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	proto "github.com/gogo/protobuf/proto"
)

const (
	TypeExecDelegate        = "exec_delegated"
	TypeGrantAuthorization  = "grant_authorization"
	TypeRevokeAuthorization = "revoke_authorization"
)

var (
	_ sdk.Msg = &MsgGrantAuthorization{}
	_ sdk.Msg = &MsgRevokeAuthorization{}
	_ sdk.Msg = &MsgExecDelegated{}
)

func NewMsgGrantAuthorization(granter sdk.AccAddress, grantee sdk.AccAddress, authorization AuthorizationI, expiration time.Time) (*MsgGrantAuthorization, error) {
	m := &MsgGrantAuthorization{
		Granter:    granter,
		Grantee:    grantee,
		Expiration: expiration,
	}
	err := m.SetAuthorization(authorization)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m MsgGrantAuthorization) GetAuthorization() AuthorizationI {
	authorization, ok := m.Authorization.GetCachedValue().(AuthorizationI)
	if !ok {
		return nil
	}
	return authorization
}

func (m MsgGrantAuthorization) SetAuthorization(authorization AuthorizationI) error {
	msg, ok := authorization.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}
	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return err
	}
	m.Authorization = any
	return nil
}

func (MsgGrantAuthorization) Route() string { return RouterKey }
func (MsgGrantAuthorization) Type() string  { return TypeGrantAuthorization }

func (m MsgGrantAuthorization) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Granter}
}

func (m MsgGrantAuthorization) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgGrantAuthorization) ValidateBasic() error {
	if m.Granter.Empty() {
		return sdkerrors.Wrap(ErrInvalidGranter, " Granter address is missing")
	}
	if m.Grantee.Empty() {
		return sdkerrors.Wrap(ErrInvalidGrantee, " Grantee address is missing")
	}
	if m.Expiration.Unix() < time.Now().Unix() {
		return sdkerrors.Wrap(ErrInvalidExpirationTime, " Time can't be in the past")
	}

	return nil
}

func NewMsgRevokeAuthorization(granter sdk.AccAddress, grantee sdk.AccAddress, authorizationMsgType sdk.Msg) (*MsgRevokeAuthorization, error) {
	m := &MsgRevokeAuthorization{
		Granter: granter,
		Grantee: grantee,
	}
	err := m.SetMsgType(authorizationMsgType)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (msg MsgRevokeAuthorization) SetMsgType(msgType sdk.Msg) error {
	m, ok := msgType.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", m)
	}
	any, err := types.NewAnyWithValue(m)
	if err != nil {
		return err
	}
	msg.AuthorizationMsgType = any
	return nil
}

func (msg MsgRevokeAuthorization) GetMsgType() sdk.Msg {
	msgType, ok := msg.AuthorizationMsgType.GetCachedValue().(sdk.Msg)
	if !ok {
		return nil
	}
	return msgType
}

func (msg MsgRevokeAuthorization) Route() string { return RouterKey }
func (msg MsgRevokeAuthorization) Type() string  { return TypeRevokeAuthorization }

func (msg MsgRevokeAuthorization) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Granter}
}

func (msg MsgRevokeAuthorization) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgRevokeAuthorization) ValidateBasic() error {
	if msg.Granter.Empty() {
		return ErrEmptyGranter
	}
	if msg.Grantee.Empty() {
		return ErrEmptyGrantee
	}
	return nil
}

func NewMsgExecDelegated(grantee sdk.AccAddress, msgs []sdk.Msg) (*MsgExecDelegated, error) {
	m := &MsgExecDelegated{
		Grantee: grantee,
	}
	err := m.SetMsgs(msgs)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m MsgExecDelegated) SetMsgs(msgs []sdk.Msg) error {
	for _, msg := range msgs {
		msg1, ok := msg.(proto.Message)
		if !ok {
			return fmt.Errorf("can't proto marshal %T", msg1)
		}
		any, err := types.NewAnyWithValue(msg1)
		if err != nil {
			return err
		}
		m.Msgs = append(m.Msgs, any)
	}
	return nil
}

func (m MsgExecDelegated) GetExecMsgs() ([]sdk.Msg, error) {
	msgs := []sdk.Msg{}
	for _, item := range m.Msgs {
		var msgInfo sdk.Msg
		err := ModuleCdc.UnpackAny(item, &msgInfo)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msgInfo)
	}
	return msgs, nil
}

func (MsgExecDelegated) Route() string { return RouterKey }
func (MsgExecDelegated) Type() string  { return TypeExecDelegate }

func (m MsgExecDelegated) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Grantee}
}

func (m MsgExecDelegated) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgExecDelegated) ValidateBasic() error {
	if m.Grantee.Empty() {
		return sdkerrors.Wrap(ErrInvalidGrantee, " Grantee address is empty")
	}
	return nil
}
