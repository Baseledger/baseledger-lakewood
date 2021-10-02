package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	uuid "github.com/kthomas/go.uuid"
	"github.com/unibrightio/baseledger/x/baseledger/types"
)

// GetBaseledgerTransactionCount get the total number of BaseledgerTransaction
func (k Keeper) GetBaseledgerTransactionCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BaseledgerTransactionCountKey))
	byteKey := types.KeyPrefix(types.BaseledgerTransactionCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to iint64
		panic("cannot decode count")
	}

	return count
}

// SetBaseledgerTransactionCount set the total number of BaseledgerTransaction
func (k Keeper) SetBaseledgerTransactionCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BaseledgerTransactionCountKey))
	byteKey := types.KeyPrefix(types.BaseledgerTransactionCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendBaseledgerTransaction appends a BaseledgerTransaction in the store with a new id and update the count
func (k Keeper) AppendBaseledgerTransaction(
	ctx sdk.Context,
	BaseledgerTransaction types.BaseledgerTransaction,
) string {
	// Create the BaseledgerTransaction
	count := k.GetBaseledgerTransactionCount(ctx)

	// Set the ID of the appended value
	// BaseledgerTransaction.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BaseledgerTransactionKey))
	appendedValue := k.cdc.MustMarshalBinaryBare(&BaseledgerTransaction)
	store.Set(GetBaseledgerTransactionUUIDBytes(BaseledgerTransaction.Id), appendedValue)

	// Update BaseledgerTransaction count
	k.SetBaseledgerTransactionCount(ctx, count+1)

	return BaseledgerTransaction.Id
}

// SetBaseledgerTransaction set a specific BaseledgerTransaction in the store
func (k Keeper) SetBaseledgerTransaction(ctx sdk.Context, BaseledgerTransaction types.BaseledgerTransaction) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BaseledgerTransactionKey))
	b := k.cdc.MustMarshalBinaryBare(&BaseledgerTransaction)
	store.Set(GetBaseledgerTransactionUUIDBytes(BaseledgerTransaction.Id), b)
}

// GetBaseledgerTransaction returns a BaseledgerTransaction from its id
func (k Keeper) GetBaseledgerTransaction(ctx sdk.Context, id string) types.BaseledgerTransaction {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BaseledgerTransactionKey))
	var BaseledgerTransaction types.BaseledgerTransaction
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetBaseledgerTransactionUUIDBytes(id)), &BaseledgerTransaction)
	return BaseledgerTransaction
}

// HasBaseledgerTransaction checks if the BaseledgerTransaction exists in the store
func (k Keeper) HasBaseledgerTransaction(ctx sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BaseledgerTransactionKey))
	return store.Has(GetBaseledgerTransactionUUIDBytes(id))
}

// GetBaseledgerTransactionOwner returns the creator of the BaseledgerTransaction
func (k Keeper) GetBaseledgerTransactionOwner(ctx sdk.Context, id string) string {
	return k.GetBaseledgerTransaction(ctx, id).Creator
}

// RemoveBaseledgerTransaction removes a BaseledgerTransaction from the store
func (k Keeper) RemoveBaseledgerTransaction(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BaseledgerTransactionKey))
	store.Delete(GetBaseledgerTransactionUUIDBytes(id))
}

// GetAllBaseledgerTransaction returns all BaseledgerTransaction
func (k Keeper) GetAllBaseledgerTransaction(ctx sdk.Context) (list []types.BaseledgerTransaction) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BaseledgerTransactionKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.BaseledgerTransaction
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func GetBaseledgerTransactionUUIDBytes(id string) []byte {
	uuidFromString, _ := uuid.FromString(id)
	return uuidFromString.Bytes()
}
