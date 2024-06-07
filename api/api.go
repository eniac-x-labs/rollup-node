package api

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/rollup-node/api/common/httputil"
	"github.com/eniac-x-labs/rollup-node/api/routes"
	api "github.com/eniac-x-labs/rollup-node/api/service"
)

const (
	HealthPath             = "/healthz"
	RollupWithTypePath     = "/api/v1/rollup-with-type"
	RetrieveFromDAWithType = "/api/v1/retrieve-with-type"
)

type API struct {
	log       log.Logger
	router    *chi.Mux
	apiServer *httputil.HTTPServer
	stopped   atomic.Bool
}

func NewApi(ctx context.Context, log log.Logger, apiAddress string, rollup api.RollupInter) error {
	out := &API{log: log}
	if err := out.initFromConfig(ctx, apiAddress, rollup); err != nil {
		return errors.Join(err, out.Stop(ctx))
	}
	return nil
}

func (a *API) initFromConfig(ctx context.Context, apiAddress string, rollup api.RollupInter) error {

	a.initRouter(rollup)
	if err := a.startServer(apiAddress); err != nil {
		return fmt.Errorf("failed to start API server: %w", err)
	}
	return nil
}

func (a *API) initRouter(rollup api.RollupInter) {

	svc := api.New(rollup)
	apiRouter := chi.NewRouter()
	h := routes.NewRoutes(a.log, apiRouter, svc)

	apiRouter.Use(middleware.Timeout(time.Second * 12))
	apiRouter.Use(middleware.Recoverer)
	apiRouter.Use(middleware.Heartbeat(HealthPath))

	apiRouter.Post(fmt.Sprintf(RollupWithTypePath), h.RollupWithTypePathHandler)
	apiRouter.Post(fmt.Sprintf(RetrieveFromDAWithType), h.RetrieveWithTypePathHandler)

	a.router = apiRouter
}

func (a *API) Start(ctx context.Context) error {
	return nil
}

func (a *API) Stop(ctx context.Context) error {
	var result error
	if a.apiServer != nil {
		if err := a.apiServer.Stop(ctx); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to stop API server: %w", err))
		}
	}

	a.stopped.Store(true)
	a.log.Info("API service shutdown complete")
	return result
}

func (a *API) startServer(addr string) error {
	a.log.Debug("API server listening...", "address", addr)
	srv, err := httputil.StartHTTPServer(addr, a.router)
	if err != nil {
		return fmt.Errorf("failed to start API server: %w", err)
	}
	a.log.Info("API server started", "addr", srv.Addr().String())
	a.apiServer = srv
	return nil
}

func (a *API) Stopped() bool {
	return a.stopped.Load()
}
