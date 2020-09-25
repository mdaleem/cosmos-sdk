package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(cliCtx client.CLIContext, r *mux.Router) {

	registerTxRoutes(cliCtx, r)
}
