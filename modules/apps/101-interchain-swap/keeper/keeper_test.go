package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"

	"github.com/stretchr/testify/suite"

	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
	ibctesting "github.com/sideprotocol/ibcswap/v4/testing"
)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain

	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))

	queryHelper := baseapp.NewQueryServerTestHelper(suite.chainA.GetContext(), suite.chainA.GetSimApp().InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.chainA.GetSimApp().IBCInterchainSwapKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func NewInterchainSwapPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.InterChainSwapPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.InterChainSwapPort
	path.EndpointA.ChannelConfig.Version = types.Version
	path.EndpointB.ChannelConfig.Version = types.Version
	return path
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
