package types

import (
	"time"
)

type GenesisState struct {
	Canvases []Canvas `json:"canvases"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Canvases: []Canvas{
			{
				Id:               "genesis",
				Width:            500,
				Height:           500,
				RefundDuration:   time.Minute * 5,
				AllowDenomPrefix: "",
				PriceForPoint:    1000000,
			},
		},
	}
}
