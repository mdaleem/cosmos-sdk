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
	TypeExecDelegate     = "exec_delegated"
	TypeGrantCapability  = "grant_capability"
	TypeRevokeCapability = "revoke_capability"
)

var (
	_ sdk.Msg = &MsgGrantCapability{}
	_ sdk.Msg = &MsgRevokeCapability{}
	_ sdk.Msg = &MsgExecDelegated{}
)

func NewMsgGrantCapability(granter sdk.AccAddress, grantee sdk.AccAddress, capability CapabilityI, expiration time.Time) (*MsgGrantCapability, error) {
	m := &MsgGrantCapability{
		Granter:    granter.String(),
		Grantee:    grantee.String(),
		Expiration: expiration,
	}
	err := m.SetCapability(capability)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m MsgGrantCapability) GetCapability() CapabilityI {
	capability, ok := m.Capability.GetCachedValue().(CapabilityI)
	if !ok {
		return nil
	}
	return capability
}

func (m MsgGrantCapability) SetCapability(capability CapabilityI) error {
	msg, ok := capability.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}
	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return err
	}
	m.Capability = any
	return nil
}

func (MsgGrantCapability) Route() string { return RouterKey }
func (MsgGrantCapability) Type() string  { return TypeGrantCapability }

func (m MsgGrantCapability) GetSigners() []sdk.AccAddress {
	granter, _ := sdk.AccAddressFromBech32(m.Granter)
	return []sdk.AccAddress{granter}
}

func (m MsgGrantCapability) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgGrantCapability) ValidateBasic() error {
	if m.Granter == "" {
		return sdkerrors.Wrap(ErrInvalidGranter, " Granter address is empty")
	}
	if m.Grantee == "" {
		return sdkerrors.Wrap(ErrInvalidGrantee, " Grantee address is empty")
	}
	if m.Expiration.Unix() < time.Now().Unix() {
		return sdkerrors.Wrap(ErrInvalidExpirationTime, " Time can't be in the past")
	}

	return nil
}

func NewMsgRevokeCapability(granter sdk.AccAddress, grantee sdk.AccAddress, msgType string) *MsgRevokeCapability {
	m := &MsgRevokeCapability{
		Granter:           granter.String(),
		Grantee:           grantee.String(),
		CapabilityMsgType: msgType,
	}
	return m
}

func (msg MsgRevokeCapability) Route() string { return RouterKey }
func (msg MsgRevokeCapability) Type() string  { return TypeRevokeCapability }

func (msg MsgRevokeCapability) GetSigners() []sdk.AccAddress {
	granter, _ := sdk.AccAddressFromBech32(msg.Granter)
	return []sdk.AccAddress{granter}
}

func (msg MsgRevokeCapability) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgRevokeCapability) ValidateBasic() error {
	if msg.Granter == "" {
		return ErrEmptyGranter
	}
	if msg.Grantee == "" {
		return ErrEmptyGrantee
	}
	return nil
}

func NewMsgExecDelegated(grantee sdk.AccAddress, msgs []sdk.Msg) (*MsgExecDelegated, error) {
	m := &MsgExecDelegated{
		Grantee: grantee.String(),
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
	grantee, _ := sdk.AccAddressFromBech32(m.Grantee)
	return []sdk.AccAddress{grantee}
}

func (m MsgExecDelegated) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgExecDelegated) ValidateBasic() error {
	if m.Grantee == "" {
		return sdkerrors.Wrap(ErrInvalidGrantee, " Grantee address is empty")
	}
	return nil
}
