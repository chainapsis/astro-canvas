package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateCanvas{}, "canvas/MsgCreateCanvas", nil)
	cdc.RegisterConcrete(MsgPaint{}, "canvas/MsgPaint", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgCreateCanvas{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgPaint{})
}

var (
	amino = codec.New()

	ModuleCdc = codec.NewHybridCodec(amino, cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(amino)
	amino.Seal()
}
