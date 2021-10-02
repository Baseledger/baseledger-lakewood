package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unibrightio/baseledger/common"
	"github.com/unibrightio/baseledger/logger"
	"github.com/unibrightio/baseledger/x/baseledger/types"
)

func (k msgServer) CreateBaseledgerTransaction(goCtx context.Context, msg *types.MsgCreateBaseledgerTransaction) (*types.MsgCreateBaseledgerTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	txCreatorAddress, err := sdk.AccAddressFromBech32(msg.Creator)

	// faucet address needs to be hard coded like this, otherwise some node could change configuration and send to arbitrary acc
	faucetAccAddress, err := sdk.AccAddressFromBech32(common.UbtFaucetAddress)

	if err != nil {
		panic(err)
	}

	coinFee, err := sdk.ParseCoinsNormalized(common.WorkTokenFee)
	if err != nil {
		panic(err)
	}

	err = k.bankKeeper.SendCoins(ctx, txCreatorAddress, faucetAccAddress, coinFee)
	if err != nil {
		logger.Errorf("send 1 token error %v\n", err.Error())
		return nil, err
	}

	var BaseledgerTransaction = types.BaseledgerTransaction{
		Id:                      msg.Id,
		Creator:                 msg.Creator,
		BaseledgerTransactionId: msg.BaseledgerTransactionId,
		Payload:                 msg.Payload,
	}

	id := k.AppendBaseledgerTransaction(
		ctx,
		BaseledgerTransaction,
	)

	return &types.MsgCreateBaseledgerTransactionResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateBaseledgerTransaction(goCtx context.Context, msg *types.MsgUpdateBaseledgerTransaction) (*types.MsgUpdateBaseledgerTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var BaseledgerTransaction = types.BaseledgerTransaction{
		Creator:                 msg.Creator,
		Id:                      msg.Id,
		BaseledgerTransactionId: msg.BaseledgerTransactionId,
		Payload:                 msg.Payload,
	}

	// Checks that the element exists
	if !k.HasBaseledgerTransaction(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", msg.Id))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Creator != k.GetBaseledgerTransactionOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetBaseledgerTransaction(ctx, BaseledgerTransaction)

	return &types.MsgUpdateBaseledgerTransactionResponse{}, nil
}

func (k msgServer) DeleteBaseledgerTransaction(goCtx context.Context, msg *types.MsgDeleteBaseledgerTransaction) (*types.MsgDeleteBaseledgerTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasBaseledgerTransaction(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", msg.Id))
	}
	if msg.Creator != k.GetBaseledgerTransactionOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveBaseledgerTransaction(ctx, msg.Id)

	return &types.MsgDeleteBaseledgerTransactionResponse{}, nil
}
