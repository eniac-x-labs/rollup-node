package eip4844

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"

	metrics2 "github.com/eniac-x-labs/rollup-node/metrics"
	"github.com/eniac-x-labs/rollup-node/x/eip4844/config"
	eth "github.com/eniac-x-labs/rollup-node/x/eip4844/eth-serivce"
	"github.com/eniac-x-labs/rollup-node/x/eip4844/log"
	"github.com/eniac-x-labs/rollup-node/x/eip4844/metrics"
	"github.com/eniac-x-labs/rollup-node/x/eip4844/txmgr"
)

type Eip4844Rollup struct {
	Txmgr   txmgr.TxManager
	Config  config.CLIConfig
	Metrics metrics.Metricer
	Log     log.Logger
}

func (e *Eip4844Rollup) NewEip4844Rollup(cliCtx *cli.Context) error {

	cfg, err := config.NewConfig(cliCtx)
	if err != nil {
		return err
	}

	logger := log.NewLogger(log.AppOut(cliCtx), cfg.LogConfig).New("eip-4844")
	log.SetGlobalLogHandler(logger.GetHandler())
	e.Log = logger

	e.initMetrics(cfg)

	if err := e.initTxManager(cfg, e.Log); err != nil {
		return err
	}
	if err := e.initMetricsServer(cfg); err != nil {
		return err
	}

	return nil
}

func (e *Eip4844Rollup) initTxManager(cfg *config.CLIConfig, logger log.Logger) error {
	txManager, err := txmgr.NewSimpleTxManager("eip-4844-rollup", logger, e.Metrics, cfg.TxMgrConfig)
	if err != nil {
		return err
	}
	e.Txmgr = txManager
	return nil
}

func (e *Eip4844Rollup) initMetrics(cfg *config.CLIConfig) {
	if cfg.MetricsConfig.Enabled {
		procName := "default"
		e.Metrics = metrics.NewMetrics(procName)
	}
}

func (e *Eip4844Rollup) initMetricsServer(cfg *config.CLIConfig) error {
	if !cfg.MetricsConfig.Enabled {
		e.Log.Info("metrics disabled")
		return nil
	}
	m, ok := e.Metrics.(metrics.RegistryMetricer)
	if !ok {
		return fmt.Errorf("metrics were enabled, but metricer %T does not expose registry for metrics-server", e.Metrics)
	}
	e.Log.Debug("starting metrics server", "addr", cfg.MetricsConfig.ListenAddr, "port", cfg.MetricsConfig.ListenPort)
	addr, err := metrics2.StartServer(m.Registry(), cfg.MetricsConfig.ListenAddr, cfg.MetricsConfig.ListenPort)
	if err != nil {
		return fmt.Errorf("failed to start metrics server: %w", err)
	}
	e.Log.Info("started metrics server", "addr", addr)
	return nil
}

// sendTransaction creates & submits a transaction to the batch inbox address with the given `txData`.
// It currently uses the underlying `txmgr` to handle transaction sending & price management.
// This is a blocking method. It should not be called concurrently.
func (e *Eip4844Rollup) sendTransaction(tx types.Transaction, queue *txmgr.Queue[types.Transaction], receiptsCh chan txmgr.TxReceipt[types.Transaction]) error {
	// Do the gas estimation offline. A value of 0 will cause the [txmgr] to estimate the gas limit.
	data := tx.Data()

	var candidate *txmgr.TxCandidate
	if e.Config.UseBlobs {
		var err error
		if candidate, err = e.blobTxCandidate(data); err != nil {
			// We could potentially fall through and try a calldata tx instead, but this would
			// likely result in the chain spending more in gas fees than it is tuned for, so best
			// to just fail. We do not expect this error to trigger unless there is a serious bug
			// or configuration issue.
			return fmt.Errorf("could not create blob tx candidate: %w", err)
		}
	} else {
		candidate = e.calldataTxCandidate(data)
	}

	intrinsicGas, err := core.IntrinsicGas(candidate.TxData, nil, false, true, true, false)
	if err != nil {
		// we log instead of return an error here because txmgr can do its own gas estimation
		e.Log.Error("Failed to calculate intrinsic gas", "err", err)
	} else {
		candidate.GasLimit = intrinsicGas
	}

	queue.Send(tx, *candidate, receiptsCh)
	return nil
}

func (e *Eip4844Rollup) blobTxCandidate(data []byte) (*txmgr.TxCandidate, error) {
	var b eth.Blob
	if err := b.FromData(data); err != nil {
		return nil, fmt.Errorf("data could not be converted to blob: %w", err)
	}
	return &txmgr.TxCandidate{
		To:    &e.Config.BatchInboxAddress,
		Blobs: []*eth.Blob{&b},
	}, nil
}

func (e *Eip4844Rollup) calldataTxCandidate(data []byte) *txmgr.TxCandidate {
	return &txmgr.TxCandidate{
		To:     &e.Config.BatchInboxAddress,
		TxData: data,
	}
}

func (e *Eip4844Rollup) getBlobData(blob eth.Blob) (eth.Data, error) {
	data, err := blob.ToData()
	if err != nil {
		e.Log.Error("ignoring blob due to parse failure", "err", err)
		return nil, err
	}

	return data, nil
}
