package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreateCanvas{}
	_ sdk.Msg = &MsgPaint{}
)

func NewMsgCreateCanvas(id string, width uint64, height uint64, refundDuration time.Duration, allowDenomPrefix string, priceForPoint uint64, sender sdk.AccAddress) MsgCreateCanvas {
	return MsgCreateCanvas{
		Id:               id,
		Width:            width,
		Height:           height,
		RefundDuration:   refundDuration,
		AllowDenomPrefix: allowDenomPrefix,
		PriceForPoint:    priceForPoint,
		Sender:           sender,
	}
}

func (msg MsgCreateCanvas) Route() string {
	return RouterKey
}

func (msg MsgCreateCanvas) Type() string {
	return "create_canvas"
}

func (msg MsgCreateCanvas) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgCreateCanvas) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return bz
}

func (msg MsgCreateCanvas) ValidateBasic() error {
	return nil
}

func NewMsgPaint(id string, x uint64, y uint64, amount sdk.Coin, sender sdk.AccAddress) MsgPaint {
	return MsgPaint{
		Id:     id,
		X:      x,
		Y:      y,
		Amount: amount,
		Sender: sender,
	}
}

func (msg MsgPaint) Route() string {
	return RouterKey
}

func (msg MsgPaint) Type() string {
	return "paint"
}

func (msg MsgPaint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgPaint) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return bz
}

func (msg MsgPaint) ValidateBasic() error {
	return nil
}
