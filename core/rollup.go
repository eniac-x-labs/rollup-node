package core

import (
	"context"
	"errors"
	"github.com/eniac-x-labs/rollup-node/common/cliapp"
	_errors "github.com/eniac-x-labs/rollup-node/common/errors"
	"github.com/eniac-x-labs/rollup-node/config"
	"github.com/eniac-x-labs/rollup-node/log"
	"github.com/eniac-x-labs/rollup-node/x/anytrust"
	"github.com/eniac-x-labs/rollup-node/x/celestia"
	"github.com/eniac-x-labs/rollup-node/x/eigenda"
	"github.com/eniac-x-labs/rollup-node/x/eip4844"
	"github.com/eniac-x-labs/rollup-node/x/nearda"
	"github.com/urfave/cli/v2"
	"sync/atomic"
)

var ErrAlreadyStopped = errors.New("already stopped")

type RollupModule struct {
	anytrustDA *anytrust.AnytrustDA
	celestiaDA *celestia.CelestiaRollup
	eigenDA    *eigenda.EigenDAClient
	eip4844    *eip4844.Eip4844Rollup
	nearDA     *nearda.NearDAClient
	stopped    atomic.Bool
	Log        log.Logger
}

func (r *RollupModule) Start(ctx context.Context) error {
	return nil
}

func (r *RollupModule) Stop(ctx context.Context) error {
	if r.stopped.Load() {
		return ErrAlreadyStopped
	}
	r.Log.Info("Stopping rollup node service")

	r.stopped.Store(true)
	r.Log.Info("rollup node service stopped")
	return nil
}

func (r *RollupModule) Stopped() bool {
	return r.stopped.Load()
}

func NewRollupModule(cliCtx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {

	cfg, err := config.NewConfig(cliCtx)
	if err != nil {
		return nil, err
	}

	logger := log.NewLogger(log.AppOut(cliCtx), cfg.LogConfig).New("rollup-node")
	log.SetGlobalLogHandler(logger.GetHandler())

	anytrustDA, err := anytrust.NewAnytrustDA()
	if err != nil {

	}
	celestiaDa, err := celestia.NewCelestiaRollup(cliCtx, logger)
	if err != nil {

	}
	eip4844, err := eip4844.NewEip4844Rollup(cliCtx, logger)
	if err != nil {

	}

	return &RollupModule{
		anytrustDA: anytrustDA,
		celestiaDA: celestiaDa,
		eigenDA:    nil,
		eip4844:    eip4844,
		nearDA:     nil,
		Log:        logger,
	}, nil
}

func (r *RollupModule) RollupWithType(data []byte, daType int) ([]interface{}, error) {
	switch daType {
	case AnytrustType:
	case CelestiaType:
		r.celestiaDA.SendTransaction(data)
	case EigenDAType:
		r.eip4844.SendTransaction(data)
	case Eip4844Type:
	case NearDAType:
	default:
		log.Error("rollup with unknown da type", "daType", daType, "expected", "[0,4]")
		return nil, _errors.UnknownDATypeErr
	}
}
