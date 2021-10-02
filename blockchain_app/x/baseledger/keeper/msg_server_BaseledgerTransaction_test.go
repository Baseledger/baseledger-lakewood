package keeper

// import (
// 	"testing"

// 	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"

// 	"github.com/unibrightio/baseledger/x/baseledger/types"
// )

// func TestBaseledgerTransactionMsgServerCreate(t *testing.T) {
// 	srv, ctx := setupMsgServer(t)
// 	creator := "A"
// 	for i := 0; i < 5; i++ {
// 		resp, err := srv.CreateBaseledgerTransaction(ctx, &types.MsgCreateBaseledgerTransaction{Creator: creator})
// 		require.NoError(t, err)
// 		assert.Equal(t, i, int(resp.Id))
// 	}
// }

// func TestBaseledgerTransactionMsgServerUpdate(t *testing.T) {
// 	creator := "A"

// 	for _, tc := range []struct {
// 		desc    string
// 		request *types.MsgUpdateBaseledgerTransaction
// 		err     error
// 	}{
// 		{
// 			desc:    "Completed",
// 			request: &types.MsgUpdateBaseledgerTransaction{Creator: creator},
// 		},
// 		{
// 			desc:    "Unauthorized",
// 			request: &types.MsgUpdateBaseledgerTransaction{Creator: "B"},
// 			err:     sdkerrors.ErrUnauthorized,
// 		},
// 		{
// 			desc:    "Unauthorized",
// 			request: &types.MsgUpdateBaseledgerTransaction{Creator: creator, Id: 10},
// 			err:     sdkerrors.ErrKeyNotFound,
// 		},
// 	} {
// 		tc := tc
// 		t.Run(tc.desc, func(t *testing.T) {
// 			srv, ctx := setupMsgServer(t)
// 			_, err := srv.CreateBaseledgerTransaction(ctx, &types.MsgCreateBaseledgerTransaction{Creator: creator})
// 			require.NoError(t, err)

// 			_, err = srv.UpdateBaseledgerTransaction(ctx, tc.request)
// 			if tc.err != nil {
// 				require.ErrorIs(t, err, tc.err)
// 			} else {
// 				require.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func TestBaseledgerTransactionMsgServerDelete(t *testing.T) {
// 	creator := "A"

// 	for _, tc := range []struct {
// 		desc    string
// 		request *types.MsgDeleteBaseledgerTransaction
// 		err     error
// 	}{
// 		{
// 			desc:    "Completed",
// 			request: &types.MsgDeleteBaseledgerTransaction{Creator: creator},
// 		},
// 		{
// 			desc:    "Unauthorized",
// 			request: &types.MsgDeleteBaseledgerTransaction{Creator: "B"},
// 			err:     sdkerrors.ErrUnauthorized,
// 		},
// 		{
// 			desc:    "KeyNotFound",
// 			request: &types.MsgDeleteBaseledgerTransaction{Creator: creator, Id: 10},
// 			err:     sdkerrors.ErrKeyNotFound,
// 		},
// 	} {
// 		tc := tc
// 		t.Run(tc.desc, func(t *testing.T) {
// 			srv, ctx := setupMsgServer(t)

// 			_, err := srv.CreateBaseledgerTransaction(ctx, &types.MsgCreateBaseledgerTransaction{Creator: creator})
// 			require.NoError(t, err)
// 			_, err = srv.DeleteBaseledgerTransaction(ctx, tc.request)
// 			if tc.err != nil {
// 				require.ErrorIs(t, err, tc.err)
// 			} else {
// 				require.NoError(t, err)
// 			}
// 		})
// 	}
// }
