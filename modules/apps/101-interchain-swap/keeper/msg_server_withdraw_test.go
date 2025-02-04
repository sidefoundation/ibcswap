package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/keeper"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
	"github.com/sideprotocol/ibcswap/v4/testing/testutil/sample"
)

func (suite *KeeperTestSuite) TestMsgWithdraw() {
	var msg *types.MsgWithdrawRequest
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
		msg = types.NewMsgWithdraw(
			*pooId,
			suite.chainA.SenderAccount.GetAddress().String(),
			sdk.DefaultBondDenom,
		)

		tc.malleate()
		msgSrv := keeper.NewMsgServerImpl(suite.chainA.GetSimApp().IBCInterchainSwapKeeper)

		res, err := msgSrv.Withdraw(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

		if tc.expPass {
			suite.Require().NoError(err)
			suite.Require().NotNil(res)
		} else {
			suite.Require().Error(err)
			suite.Require().Nil(res)
		}
	}
}
