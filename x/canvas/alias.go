package canvas

import (
	"github.com/chainapsis/astro-canvas/x/canvas/keeper"
	"github.com/chainapsis/astro-canvas/x/canvas/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
)

var (
	NewKeeper = keeper.NewCanvasKeeper

	RegisterCodec      = types.RegisterCodec
	RegisterInterfaces = types.RegisterInterfaces

	DefaultGenesisState = types.DefaultGenesisState
)

type (
	Keeper          = keeper.Keeper
	MsgCreateCanvas = types.MsgCreateCanvas
	MsgPaint        = types.MsgPaint

	GenesisState = types.GenesisState
)
