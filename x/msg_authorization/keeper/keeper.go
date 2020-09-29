package keeper

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/msg_authorization/types"
	proto "github.com/gogo/protobuf/proto"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.BinaryMarshaler
	router   sdk.Router
}

// NewKeeper constructs a message authorisation Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryMarshaler, router sdk.Router) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
		router:   router,
	}
}

func (k Keeper) getActorCapabilityKey(grantee sdk.AccAddress, granter sdk.AccAddress, msg sdk.Msg) []byte {
	return []byte(fmt.Sprintf("c/%x/%x/%s/%s", grantee, granter, msg.Route(), msg.Type()))
}

func (k Keeper) getCapabilityGrant(ctx sdk.Context, actor []byte) (grant types.CapabilityGrant, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(actor)
	if bz == nil {
		return grant, false
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &grant)
	return grant, true
}

func (k Keeper) update(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, updated types.CapabilityI) {
	actor := k.getActorCapabilityKey(grantee, granter, updated.MsgType())
	grant, found := k.getCapabilityGrant(ctx, actor)
	if !found {
		return
	}

	msg, ok := updated.(proto.Message)
	if !ok {
		panic(fmt.Errorf("cannot proto marshal %T", updated))
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		panic(err)
	}

	grant.Capability = any
	store := ctx.KVStore(k.storeKey)
	store.Set(actor, k.cdc.MustMarshalBinaryBare(&grant))
}

// DispatchActions attempts to execute the provided messages via authorization
// grants from the message signer to the grantee.
func (k Keeper) DispatchActions(ctx sdk.Context, grantee sdk.AccAddress, msgs []sdk.Msg) (*sdk.Result, error) {
	var msgResult *sdk.Result
	var err error
	for _, msg := range msgs {
		signers := msg.GetSigners()
		if len(signers) != 1 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "authorization can be given to msg with only one signer")
		}
		granter := signers[0]
		if !bytes.Equal(granter, grantee) {
			capability, _ := k.GetCapability(ctx, grantee, granter, msg)
			if capability == nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "authorization not found")
			}
			allow, updated, del := capability.Accept(msg, ctx.BlockHeader())
			if !allow {
				return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "authorization not found")
			}
			if del {
				k.Revoke(ctx, grantee, granter, msg)
			} else if updated != nil {
				k.update(ctx, grantee, granter, updated)
			}
		}
		handler := k.router.Route(ctx, msg.Route())
		if handler == nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message route: %s", msg.Route())
		}

		msgResult, err = handler(ctx, msg)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "failed to execute message; message %s", msg.Type())
		}
	}

	return msgResult, nil
}

// Grant method grants the provided authorization to the grantee on the granter's account with the provided expiration
// time. If there is an existing authorization grant for the same `sdk.Msg` type, this grant
// overwrites that.
func (k Keeper) Grant(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, capability types.CapabilityI, expiration time.Time) {
	store := ctx.KVStore(k.storeKey)
	grant, err := types.NewGrantCapability(capability, expiration.Unix())
	if err != nil {
		sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "authorization can not be given to msg")
	}
	bz := k.cdc.MustMarshalBinaryBare(grant)
	actor := k.getActorCapabilityKey(grantee, granter, capability.MsgType())
	store.Set(actor, bz)

}

// Revoke method revokes any capability for the provided message type granted to the grantee by the granter.
func (k Keeper) Revoke(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg) error {
	store := ctx.KVStore(k.storeKey)
	actor := k.getActorCapabilityKey(grantee, granter, msgType)
	_, found := k.getCapabilityGrant(ctx, actor)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "capability not found")
	}
	store.Delete(actor)

	return nil
}

// GetCapability Returns any `Capability` (or `nil`), with the expiration time,
// granted to the grantee by the granter for the provided msg type.
func (k Keeper) GetCapability(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg) (cap types.CapabilityI, expiration int64) {
	grant, found := k.getCapabilityGrant(ctx, k.getActorCapabilityKey(grantee, granter, msgType))
	if !found {
		return nil, 0
	}

	if grant.Expiration != 0 && grant.Expiration < (ctx.BlockHeader().Time.Unix()) {
		k.Revoke(ctx, grantee, granter, msgType)
		return nil, 0
	}

	var capability types.CapabilityI
	err := k.cdc.UnpackAny(grant.Capability, &capability)
	if err != nil {
		return nil, 0
	}

	return capability, grant.Expiration
}