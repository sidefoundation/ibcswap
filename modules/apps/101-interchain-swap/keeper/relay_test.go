package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	clienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
	ibctesting "github.com/sideprotocol/ibcswap/v4/testing"
	"github.com/sideprotocol/ibcswap/v4/testing/testutil/sample"
)

// test sending from chainA to chainB using both coin that orignate on
// chainA and coin that orignate on chainB
func (suite *KeeperTestSuite) TestSendSwap() {
	var (
		//amount sdk.Coin
		msgbyte []byte
		path    *ibctesting.Path
		err     error
	)

	testCases := []struct {
		msg            string
		malleate       func()
		sendFromSource bool
		expPass        bool
	}{
		{
			"successful transfer swap request",
			func() {
				suite.coordinator.CreateInterchainSwapChannels(path)
				msg := &types.MsgSwapRequest{
					SwapType: types.SwapMsgType_LEFT,
					Sender:   sample.AccAddress(),
					TokenIn: &sdk.Coin{
						Denom:  sdk.DefaultBondDenom,
						Amount: sdk.NewInt(100),
					},
					TokenOut: &sdk.Coin{
						Denom:  sdk.DefaultBondDenom,
						Amount: sdk.NewInt(100),
					},
				}

				msgbyte, err = types.ModuleCdc.Marshal(msg)
				suite.Require().NoError(err)

			}, true, true,
		},
		{
			"successful transfer creat pool request",
			func() {
				suite.coordinator.CreateInterchainSwapChannels(path)
				msg := types.NewMsgCreatePool(
					path.EndpointA.ChannelConfig.PortID,
					path.EndpointA.ChannelID,
					suite.chainA.SenderAccount.GetAddress().String(),
					"1:2",
					[]string{sdk.DefaultBondDenom, "venuscoin"},
					[]uint32{10, 100},
				)

				msgbyte, err = types.ModuleCdc.Marshal(msg)
				suite.Require().NoError(err)
			}, true, true,
		},
		{
			"successful transfer deposit request",
			func() {
				suite.coordinator.CreateInterchainSwapChannels(path)
				msg := types.NewMsgDeposit(
					"test pool id",
					suite.chainA.SenderAccount.GetAddress().String(),
					[]*sdk.Coin{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(1000)}},
				)

				msgbyte, err = types.ModuleCdc.Marshal(msg)
				suite.Require().NoError(err)
			}, true, true,
		},
		{
			"successful transfer withdraw request",
			func() {
				suite.coordinator.CreateInterchainSwapChannels(path)
				msg := types.NewMsgWithdraw(
					"test_id",
					suite.chainA.SenderAccount.GetAddress().String(),
					sdk.DefaultBondDenom,
				)

				msgbyte, err = types.ModuleCdc.Marshal(msg)
				suite.Require().NoError(err)
			}, true, true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			path = NewInterchainSwapPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)

			tc.malleate()
			packet := types.IBCSwapDataPacket{
				Type: types.MessageType_LEFTSWAP,
				Data: msgbyte,
			}

			err = suite.chainA.GetSimApp().IBCInterchainSwapKeeper.SendIBCSwapPacket(
				suite.chainA.GetContext(),
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
				clienttypes.NewHeight(0, 110), 0,
				packet,
			)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
