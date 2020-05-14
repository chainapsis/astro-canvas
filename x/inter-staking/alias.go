package inter_staking

import (
	"github.com/chainapsis/astro-canvas/x/inter-staking/keeper"
	"github.com/chainapsis/astro-canvas/x/inter-staking/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	RouterKey    = types.RouterKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper          = keeper.NewInterStakingKeeper
	RegisterCodec      = types.RegisterCodec
	RegisterInterfaces = types.RegisterInterfaces
	NewMsgRegister     = types.NewMsgRegister
	NewMsgDelegate     = types.NewMsgDelegate
)

type (
	Keeper      = keeper.Keeper
	MsgRegister = types.MsgRegister
	MsgDelegate = types.MsgDelegate
)
