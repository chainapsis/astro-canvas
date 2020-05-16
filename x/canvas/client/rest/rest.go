package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/canvas/{id}", QueryCanvasRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/canvas/{id}/points", QueryPointsRequestHandlerFn(cliCtx)).Methods("GET")
}
