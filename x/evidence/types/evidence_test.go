package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/evidence/types"
)

func TestEquivocation_Valid(t *testing.T) {
	n, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	e := types.Equivocation{
		Height:           100,
		Time:             n,
		Power:            1000000,
		ConsensusAddress: sdk.ConsAddress("__________foo_______").String(),
	}

	require.Equal(t, e.GetTotalPower(), int64(0))
	require.Equal(t, e.GetValidatorPower(), e.Power)
	require.Equal(t, e.GetTime(), e.Time)
	require.Equal(t, e.GetConsensusAddress().String(), e.ConsensusAddress)
	require.Equal(t, e.GetHeight(), e.Height)
	require.Equal(t, e.Type(), types.TypeEquivocation)
	require.Equal(t, e.Route(), types.RouteEquivocation)
	require.Equal(t, e.Hash().String(), "1ECFBCF3081A428FC28F3175595BC780DFD41499D880BCA69C2B33EBB1EDC186")
	require.Equal(t, e.String(), "height: 100\ntime: 2006-01-02T15:04:05Z\npower: 1000000\nconsensus_address: cosmosvalcons1ta047h6lta047h6lvehk7h6lta047h6lyp568c\n")
	require.NoError(t, e.ValidateBasic())
}

func TestEquivocationValidateBasic(t *testing.T) {
	var zeroTime time.Time

	n, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	testCases := []struct {
		name      string
		e         types.Equivocation
		expectErr bool
	}{
		{"valid", types.Equivocation{100, n, 1000000, sdk.ConsAddress("foo").String()}, false},
		{"invalid time", types.Equivocation{100, zeroTime, 1000000, sdk.ConsAddress("foo").String()}, true},
		{"invalid height", types.Equivocation{0, n, 1000000, sdk.ConsAddress("foo").String()}, true},
		{"invalid power", types.Equivocation{100, n, 0, sdk.ConsAddress("foo").String()}, true},
		{"invalid address", types.Equivocation{100, n, 1000000, ""}, true},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectErr, tc.e.ValidateBasic() != nil)
		})
	}
}
