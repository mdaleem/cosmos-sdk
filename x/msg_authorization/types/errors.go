package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidGranter        = sdkerrors.Register(ModuleName, 101, "Invalid granter address")
	ErrEmptyGranter          = sdkerrors.Register(ModuleName, 102, "Empty granter address")
	ErrInvalidGrantee        = sdkerrors.Register(ModuleName, 103, "Invalid grantee address")
	ErrEmptyGrantee          = sdkerrors.Register(ModuleName, 104, "Empty grantee address")
	ErrInvalidExpirationTime = sdkerrors.Register(ModuleName, 105, "expiration time of authorization should be more than current time")
)
