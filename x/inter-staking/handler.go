package inter_staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgRegister:
			return handleMsgRegister(ctx, msg, k)
		case MsgDelegate:
			return handleMsgDelegate(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgRegister(ctx sdk.Context, msg MsgRegister, k Keeper) (*sdk.Result, error) {
	err := k.RegisterInterchainAccount(ctx, msg.Sender, msg.SourcePort, msg.SourceChannel)

	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDelegate(ctx sdk.Context, msg MsgDelegate, k Keeper) (*sdk.Result, error) {
	err := k.Delegate(ctx, msg.CounterpartyBech32Addr, msg.DelegatorAddress, msg.ValidatorAddress, msg.Amount,
		msg.TransferSourcePort, msg.TransferSourceChannel, msg.InterchainAccountSourcePort, msg.InterchainAccountSourceChannel)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
