package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = MsgRegister{}

func NewMsgRegister(sourcePort string, sourceChannel string, sender sdk.AccAddress) MsgRegister {
	return MsgRegister{SourcePort: sourcePort, SourceChannel: sourceChannel, Sender: sender}
}

func (MsgRegister) Route() string {
	return RouterKey
}

func (MsgRegister) Type() string {
	return RouterKey
}

func (MsgRegister) ValidateBasic() error {
	return nil
}

func (msg MsgRegister) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRegister) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

var _ sdk.Msg = MsgDelegate{}

func NewMsgDelegate(transferSourcePort, transferSourceChannel, iaSourcePort, iaSourceChannel string, counterpartyBech32Addr string, delegator sdk.AccAddress, validator sdk.ValAddress, amount sdk.Coin) MsgDelegate {
	return MsgDelegate{
		TransferSourcePort:             transferSourcePort,
		TransferSourceChannel:          transferSourceChannel,
		InterchainAccountSourcePort:    iaSourcePort,
		InterchainAccountSourceChannel: iaSourceChannel,
		CounterpartyBech32Addr:         counterpartyBech32Addr,
		DelegatorAddress:               delegator,
		ValidatorAddress:               validator,
		Amount:                         amount,
	}
}

func (MsgDelegate) Route() string {
	return RouterKey
}

func (MsgDelegate) Type() string {
	return RouterKey
}

func (MsgDelegate) ValidateBasic() error {
	return nil
}

func (msg MsgDelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelegatorAddress}
}
