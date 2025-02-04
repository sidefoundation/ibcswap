package types

import (
	"testing"

	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v4/testing/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgDeposit_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDepositRequest
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDepositRequest{
				Sender: "invalid_address",
			},
			err: errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", ""),
		}, {
			name: "valid address",
			msg: MsgDepositRequest{
				Sender: sample.AccAddress(),
				Tokens: []*types.Coin{
					{
						Denom:  "atom",
						Amount: types.NewInt(0),
					},
					{
						Denom:  "marscoin",
						Amount: types.NewInt(0),
					}},
			},
		},

		{
			name: "invalid denom length",
			msg: MsgDepositRequest{
				Sender: sample.AccAddress(),
			},
			err: errorsmod.Wrapf(ErrInvalidTokenLength, "invalid token length (%d)", 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
