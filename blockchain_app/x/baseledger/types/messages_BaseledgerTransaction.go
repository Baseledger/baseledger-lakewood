package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateBaseledgerTransaction{}

func NewMsgCreateBaseledgerTransaction(id string, creator string, baseledgerTransactionId string, payload string, opCode uint32) *MsgCreateBaseledgerTransaction {
	return &MsgCreateBaseledgerTransaction{
		Id:                      id,
		Creator:                 creator,
		BaseledgerTransactionId: baseledgerTransactionId,
		Payload:                 payload,
		OpCode:                  opCode,
	}
}

func (msg *MsgCreateBaseledgerTransaction) Route() string {
	return RouterKey
}

func (msg *MsgCreateBaseledgerTransaction) Type() string {
	return "CreateBaseledgerTransaction"
}

func (msg *MsgCreateBaseledgerTransaction) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateBaseledgerTransaction) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateBaseledgerTransaction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateBaseledgerTransaction{}

func NewMsgUpdateBaseledgerTransaction(creator string, id string, baseledgerTransactionId string, payload string, opCode uint32) *MsgUpdateBaseledgerTransaction {
	return &MsgUpdateBaseledgerTransaction{
		Id:                      id,
		Creator:                 creator,
		BaseledgerTransactionId: baseledgerTransactionId,
		Payload:                 payload,
		OpCode:                  opCode,
	}
}

func (msg *MsgUpdateBaseledgerTransaction) Route() string {
	return RouterKey
}

func (msg *MsgUpdateBaseledgerTransaction) Type() string {
	return "UpdateBaseledgerTransaction"
}

func (msg *MsgUpdateBaseledgerTransaction) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateBaseledgerTransaction) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateBaseledgerTransaction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteBaseledgerTransaction{}

func NewMsgDeleteBaseledgerTransaction(creator string, id string) *MsgDeleteBaseledgerTransaction {
	return &MsgDeleteBaseledgerTransaction{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeleteBaseledgerTransaction) Route() string {
	return RouterKey
}

func (msg *MsgDeleteBaseledgerTransaction) Type() string {
	return "DeleteBaseledgerTransaction"
}

func (msg *MsgDeleteBaseledgerTransaction) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteBaseledgerTransaction) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteBaseledgerTransaction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
