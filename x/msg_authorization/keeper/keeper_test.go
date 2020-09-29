package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/msg_authorization/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type TestSuite struct {
	suite.Suite

	app   *simapp.SimApp
	ctx   sdk.Context
	addrs []sdk.AccAddress
}

func (s *TestSuite) SetupTest() {
	s.app = simapp.Setup(false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
}

func (s *TestSuite) TestKeeper() {
	app, ctx := s.app, s.ctx

	granterAddr := sdk.AccAddress("")
	granteeAddr := sdk.AccAddress("")
	recipientAddr := sdk.AccAddress("")
	err := app.BankKeeper.SetBalances(ctx, granterAddr, sdk.NewCoins(sdk.NewInt64Coin("steak", 10000)))
	s.Require().Nil(err)
	s.Require().True(app.BankKeeper.GetBalance(ctx, granterAddr, "steak").IsEqual(sdk.NewCoin("steak", sdk.NewInt(10000))))

	s.T().Log("verify that no authorization returns nil")
	authorization, expiration := app.MsgAuthKeeper.GetCapability(ctx, granteeAddr, granterAddr, &banktypes.MsgSend{})
	s.Require().Nil(authorization)
	s.Require().Zero(expiration)
	now := s.ctx.BlockHeader().Time
	s.Require().NotNil(now)

	newCoins := sdk.NewCoins(sdk.NewInt64Coin("steak", 100))
	s.T().Log("verify if expired authorization is rejected")
	x := types.SendCapability{SpendLimit: newCoins}
	s.app.MsgAuthKeeper.Grant(ctx, granterAddr, granteeAddr, x, now.Add(-1*time.Hour))
	authorization, _ = s.app.MsgAuthKeeper.GetCapability(ctx, granteeAddr, granterAddr, &banktypes.MsgSend{})
	s.Require().Nil(authorization)

	s.T().Log("verify if authorization is accepted")
	x = types.SendCapability{SpendLimit: newCoins}
	s.app.MsgAuthKeeper.Grant(ctx, granteeAddr, granterAddr, x, now.Add(time.Hour))
	authorization, _ = s.app.MsgAuthKeeper.GetCapability(ctx, granteeAddr, granterAddr, &banktypes.MsgSend{})
	s.Require().NotNil(authorization)
	s.Require().Equal(authorization.MsgType(), banktypes.MsgSend{}.Type())

	s.T().Log("verify fetching authorization with wrong msg type fails")
	authorization, _ = s.app.MsgAuthKeeper.GetCapability(ctx, granteeAddr, granterAddr, &banktypes.MsgMultiSend{})
	s.Require().Nil(authorization)

	s.T().Log("verify fetching authorization with wrong grantee fails")
	authorization, _ = s.app.MsgAuthKeeper.GetCapability(ctx, recipientAddr, granterAddr, &banktypes.MsgMultiSend{})
	s.Require().Nil(authorization)

	s.T().Log("")

	s.T().Log("verify revoke fails with wrong information")
	s.app.MsgAuthKeeper.Revoke(ctx, recipientAddr, granterAddr, &banktypes.MsgSend{})
	authorization, _ = s.app.MsgAuthKeeper.GetCapability(ctx, recipientAddr, granterAddr, &banktypes.MsgSend{})
	s.Require().Nil(authorization)

	s.T().Log("verify revoke executes with correct information")
	s.app.MsgAuthKeeper.Revoke(ctx, recipientAddr, granterAddr, &banktypes.MsgSend{})
	authorization, _ = s.app.MsgAuthKeeper.GetCapability(ctx, granteeAddr, granterAddr, &banktypes.MsgSend{})
	s.Require().NotNil(authorization)

}

func (s *TestSuite) TestKeeperFees() {
	app := s.app

	granterAddr := sdk.AccAddress("")
	granteeAddr := sdk.AccAddress("")
	recipientAddr := sdk.AccAddress("")
	err := app.BankKeeper.SetBalances(s.ctx, granterAddr, sdk.NewCoins(sdk.NewInt64Coin("steak", 10000)))
	s.Require().Nil(err)
	s.Require().True(app.BankKeeper.GetBalance(s.ctx, granterAddr, "steak").IsEqual(sdk.NewCoin("steak", sdk.NewInt(10000))))

	now := s.ctx.BlockHeader().Time
	s.Require().NotNil(now)

	smallCoin := sdk.NewCoins(sdk.NewInt64Coin("steak", 20))
	someCoin := sdk.NewCoins(sdk.NewInt64Coin("steak", 123))
	//lotCoin := sdk.NewCoins(sdk.NewInt64Coin("steak", 4567))
	m := []sdk.Msg{
		&banktypes.MsgSend{
			Amount:      sdk.NewCoins(sdk.NewInt64Coin("steak", 2)),
			FromAddress: granterAddr.String(),
			ToAddress:   recipientAddr.String(),
		},
	}

	msgs := types.MsgExecDelegated{
		Grantee: granteeAddr,
		Msgs:    &m,
	}

	s.T().Log("verify dispatch fails with invalid authorization")
	result, error := s.app.MsgAuthKeeper.DispatchActions(s.ctx, granteeAddr, msgs.Msgs)
	s.Require().Nil(result)
	s.Require().NotNil(error)

	s.T().Log("verify dispatch executes with correct information")
	// grant authorization
	s.app.MsgAuthKeeper.Grant(s.ctx, granteeAddr, granterAddr, types.SendCapability{SpendLimit: smallCoin}, now)
	authorization, expiration := s.app.MsgAuthKeeper.GetCapability(s.ctx, granteeAddr, granterAddr, &banktypes.MsgSend{})
	s.Require().NotNil(authorization)
	s.Require().Zero(expiration)
	s.Require().Equal(authorization.MsgType(), banktypes.MsgSend{}.Type())
	result, error = s.app.MsgAuthKeeper.DispatchActions(s.ctx, granteeAddr, msgs.Msgs)
	s.Require().NotNil(result)
	s.Require().Nil(error)

	authorization, _ = s.app.MsgAuthKeeper.GetCapability(s.ctx, granteeAddr, granterAddr, &banktypes.MsgSend{})
	s.Require().NotNil(authorization)

	s.T().Log("verify dispatch fails with overlimit")
	// grant authorization

	msgs = types.MsgExecDelegated{
		Grantee: granteeAddr,
		Msgs: []sdk.Msg{
			&banktypes.MsgSend{
				Amount:      someCoin,
				FromAddress: granterAddr.String(),
				ToAddress:   recipientAddr.String(),
			},
		},
	}

	result, error = s.app.MsgAuthKeeper.DispatchActions(s.ctx, granteeAddr, msgs.Msgs)
	s.Require().Nil(result)
	s.Require().NotNil(error)

	authorization, _ = s.app.MsgAuthKeeper.GetCapability(s.ctx, granteeAddr, granterAddr, &banktypes.MsgSend{})
	s.Require().NotNil(authorization)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
