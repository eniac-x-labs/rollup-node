package routes

import (
	"github.com/eniac-x-labs/rollup-node/api/service"
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"
)

type Routes struct {
	logger log.Logger
	router *chi.Mux
	svc    service.HandlerSvc
}

// NewRoutes ... Construct a new route handler instance
func NewRoutes(l log.Logger, r *chi.Mux, svc service.HandlerSvc) Routes {
	return Routes{
		logger: l,
		router: r,
		svc:    svc,
	}
}
