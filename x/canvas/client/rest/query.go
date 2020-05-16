package rest

import (
	"fmt"
	"github.com/chainapsis/astro-canvas/x/canvas/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func QueryCanvasRequestHandlerFn(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		id := vars["id"]

		ctx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, ctx, r)
		if !ok {
			return
		}

		var marshaler codec.JSONMarshaler

		if ctx.Marshaler != nil {
			marshaler = ctx.Marshaler
		} else {
			marshaler = ctx.Codec
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCanvas)
		params := types.QueryCanvasParams{Id: id}

		bz, err := marshaler.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, height, err := ctx.QueryWithData(route, bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		ctx = ctx.WithHeight(height)
		rest.PostProcessResponse(w, ctx, res)
	}
}

func QueryPointsRequestHandlerFn(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		id := vars["id"]

		ctx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, ctx, r)
		if !ok {
			return
		}

		var marshaler codec.JSONMarshaler

		if ctx.Marshaler != nil {
			marshaler = ctx.Marshaler
		} else {
			marshaler = ctx.Codec
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPoints)
		params := types.QueryPointsParams{Id: id}

		bz, err := marshaler.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, height, err := ctx.QueryWithData(route, bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		ctx = ctx.WithHeight(height)
		rest.PostProcessResponse(w, ctx, res)
	}
}
