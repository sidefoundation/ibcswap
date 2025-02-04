package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/keeper"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
	"github.com/sideprotocol/ibcswap/v4/testing/testutil/sample"
)

func (suite *KeeperTestSuite) SetupPool() (*string, error) {
	suite.SetupTest()
	path := NewInterchainSwapPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)
	msg := types.NewMsgCreatePool(
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		suite.chainA.SenderAccount.GetAddress().String(),
		"1:2",
		[]string{sdk.DefaultBondDenom, "venuscoin"},
		[]uint32{10, 100},
	)

	ctx := suite.chainA.GetContext()
	suite.chainA.GetSimApp().IBCInterchainSwapKeeper.OnCreatePoolAcknowledged(ctx, msg)
	poolId := types.GetPoolId(msg.Denoms)
	return &poolId, nil
}

func (suite *KeeperTestSuite) TestMsgDeposit() {
	var msg *types.MsgDepositRequest
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {},
			true,
		},
		{
			"invalid address",
			func() {
				msg.Sender = "invalid address"
			},
			false,
		},
		{
			"invalid amount",
			func() {
				msg.Sender = sample.AccAddress()
			},
			false,
		},
	}

	for _, tc := range testCases {
		// create pool first of all.
		pooId, err := suite.SetupPool()
		suite.Require().NoError(err)
		fmt.Println(pooId)

		//
		msg = types.NewMsgDeposit(
			*pooId,
			suite.chainA.SenderAccount.GetAddress().String(),
			[]*sdk.Coin{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(1000)}},
		)

		tc.malleate()
		msgSrv := keeper.NewMsgServerImpl(suite.chainA.GetSimApp().IBCInterchainSwapKeeper)

		res, err := msgSrv.Deposit(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

		if tc.expPass {
			suite.Require().NoError(err)
			suite.Require().NotNil(res)
		} else {
			suite.Require().Error(err)
			suite.Require().Nil(res)
		}
	}
}
