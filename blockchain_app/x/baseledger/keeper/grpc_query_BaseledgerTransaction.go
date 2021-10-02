package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/unibrightio/baseledger/x/baseledger/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) BaseledgerTransactionAll(c context.Context, req *types.QueryAllBaseledgerTransactionRequest) (*types.QueryAllBaseledgerTransactionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var BaseledgerTransactions []*types.BaseledgerTransaction
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	BaseledgerTransactionStore := prefix.NewStore(store, types.KeyPrefix(types.BaseledgerTransactionKey))

	pageRes, err := query.Paginate(BaseledgerTransactionStore, req.Pagination, func(key []byte, value []byte) error {
		var BaseledgerTransaction types.BaseledgerTransaction
		if err := k.cdc.UnmarshalBinaryBare(value, &BaseledgerTransaction); err != nil {
			return err
		}

		BaseledgerTransactions = append(BaseledgerTransactions, &BaseledgerTransaction)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllBaseledgerTransactionResponse{BaseledgerTransaction: BaseledgerTransactions, Pagination: pageRes}, nil
}

func (k Keeper) BaseledgerTransaction(c context.Context, req *types.QueryGetBaseledgerTransactionRequest) (*types.QueryGetBaseledgerTransactionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var BaseledgerTransaction types.BaseledgerTransaction
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasBaseledgerTransaction(ctx, req.Id) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BaseledgerTransactionKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetBaseledgerTransactionUUIDBytes(req.Id)), &BaseledgerTransaction)

	return &types.QueryGetBaseledgerTransactionResponse{BaseledgerTransaction: &BaseledgerTransaction}, nil
}
