package canvas

import (
	"github.com/chainapsis/astro-canvas/x/canvas/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	for _, canvas := range data.Canvases {
		err := keeper.CreateCanvas(ctx, canvas.Id, canvas.Width, canvas.Height, canvas.RefundDuration, canvas.AllowDenomPrefix, canvas.PriceForPoint)
		if err != nil {
			panic(err)
		}
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	// TODO...
	return types.GenesisState{Canvases: []types.Canvas{}}
}
