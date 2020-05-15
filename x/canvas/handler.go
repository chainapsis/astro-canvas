package canvas

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgCreateCanvas:
			return handleMsgCreateCanvas(ctx, msg, keeper)
		case MsgPaint:
			return handleMsgPaint(ctx, msg, keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgCreateCanvas(ctx sdk.Context, msg MsgCreateCanvas, keeper Keeper) (*sdk.Result, error) {
	err := keeper.CreateCanvas(ctx, msg.Id, msg.Width, msg.Height, msg.RefundDuration, msg.AllowDenomPrefix, msg.PriceForPoint)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgPaint(ctx sdk.Context, msg MsgPaint, keeper Keeper) (*sdk.Result, error) {
	err := keeper.Paint(ctx, msg.Id, msg.X, msg.Y, msg.Amount, msg.Sender)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
