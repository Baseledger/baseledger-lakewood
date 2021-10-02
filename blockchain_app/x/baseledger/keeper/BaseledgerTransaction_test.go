package keeper

// import (
// 	"testing"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/unibrightio/baseledger/x/baseledger/types"
// )

// func createNBaseledgerTransaction(keeper *Keeper, ctx sdk.Context, n int) []types.BaseledgerTransaction {
// 	items := make([]types.BaseledgerTransaction, n)
// 	for i := range items {
// 		items[i].Creator = "any"
// 		items[i].Id = keeper.AppendBaseledgerTransaction(ctx, items[i])
// 	}
// 	return items
// }

// func TestBaseledgerTransactionGet(t *testing.T) {
// 	keeper, ctx := setupKeeper(t)
// 	items := createNBaseledgerTransaction(keeper, ctx, 10)
// 	for _, item := range items {
// 		assert.Equal(t, item, keeper.GetBaseledgerTransaction(ctx, item.Id))
// 	}
// }

// func TestBaseledgerTransactionExist(t *testing.T) {
// 	keeper, ctx := setupKeeper(t)
// 	items := createNBaseledgerTransaction(keeper, ctx, 10)
// 	for _, item := range items {
// 		assert.True(t, keeper.HasBaseledgerTransaction(ctx, item.Id))
// 	}
// }

// func TestBaseledgerTransactionRemove(t *testing.T) {
// 	keeper, ctx := setupKeeper(t)
// 	items := createNBaseledgerTransaction(keeper, ctx, 10)
// 	for _, item := range items {
// 		keeper.RemoveBaseledgerTransaction(ctx, item.Id)
// 		assert.False(t, keeper.HasBaseledgerTransaction(ctx, item.Id))
// 	}
// }

// func TestBaseledgerTransactionGetAll(t *testing.T) {
// 	keeper, ctx := setupKeeper(t)
// 	items := createNBaseledgerTransaction(keeper, ctx, 10)
// 	assert.Equal(t, items, keeper.GetAllBaseledgerTransaction(ctx))
// }

// func TestBaseledgerTransactionCount(t *testing.T) {
// 	keeper, ctx := setupKeeper(t)
// 	items := createNBaseledgerTransaction(keeper, ctx, 10)
// 	count := uint64(len(items))
// 	assert.Equal(t, count, keeper.GetBaseledgerTransactionCount(ctx))
// }
