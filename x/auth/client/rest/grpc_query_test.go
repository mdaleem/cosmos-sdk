package rest_test

import (
	"encoding/base64"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gogo/protobuf/proto"
)

func (s *IntegrationTestSuite) TestAuthAccountsGRPCHandler() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	// TODO: need to pass bech32 string instead of base64 encoding string.
	// ref: https://github.com/cosmos/cosmos-sdk/issues/7195
	addressBase64 := base64.URLEncoding.EncodeToString(val.Address)
	fmt.Println(addressBase64)
	testCases := []struct {
		name      string
		url       string
		expectErr bool
		respType  proto.Message
		expected  proto.Message
	}{
		{
			"test GRPC account invalid address",
			fmt.Sprintf("%s/cosmos/auth/v1beta1/accounts/%s", baseURL, "invalid"),
			true,
			&types.QueryAccountResponse{},
			&types.QueryAccountResponse{},
		},
		{
			"test GRPC account empty address",
			fmt.Sprintf("%s/cosmos/auth/v1beta1/accounts/%s", baseURL, ""),
			true,
			&types.QueryAccountResponse{},
			&types.QueryAccountResponse{},
		},
		{
			"test GRPC params valid address",
			fmt.Sprintf("%s/cosmos/auth/v1beta1/accounts/%s", baseURL, addressBase64),
			false,
			&types.QueryAccountResponse{},
			&types.QueryAccountResponse{
				Account: nil,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)
			err = val.ClientCtx.JSONMarshaler.UnmarshalJSON(resp, tc.respType)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestTotalSupplyGRPCHandler() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name     string
		url      string
		respType proto.Message
		expected proto.Message
	}{
		{
			"test GRPC params",
			fmt.Sprintf("%s/cosmos/auth/v1beta1/params", baseURL),
			&types.QueryParamsResponse{},
			&types.QueryParamsResponse{
				Params: types.DefaultParams(),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)
			s.Require().NoError(val.ClientCtx.LegacyAmino.UnmarshalJSON(resp, tc.respType))
			s.Require().Equal(tc.expected.String(), tc.respType.String())
		})
	}
}
