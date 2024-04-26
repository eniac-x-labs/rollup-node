package celestia

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/eniac-x-labs/rollup-node/common/cliapp"
	"github.com/eniac-x-labs/rollup-node/config"
	eth "github.com/eniac-x-labs/rollup-node/eth-serivce"
	"github.com/eniac-x-labs/rollup-node/log"
	metrics2 "github.com/eniac-x-labs/rollup-node/metrics"
	"github.com/eniac-x-labs/rollup-node/txmgr"
	"github.com/eniac-x-labs/rollup-node/txmgr/metrics"
)

var ErrAlreadyStopped = errors.New("already stopped")

type CelestiaRollup struct {
	CelestiaConfig CLIConfig
	Txmgr          txmgr.TxManager
	Config         *config.CLIConfig
	Metrics        metrics.Metricer
	Log            log.Logger
	DAClient       *DAClient
	stopped        atomic.Bool
	driverCtx      context.Context
}

func (c *CelestiaRollup) Start(ctx context.Context) error {
	return nil
}

func (c *CelestiaRollup) Stop(ctx context.Context) error {
	if c.stopped.Load() {
		return ErrAlreadyStopped
	}

	c.Log.Info("Stopping eip4844 rollup service")

	if c.Txmgr != nil {
		c.Txmgr.Close()
	}

	c.stopped.Store(true)
	c.Log.Info("eip4844 rollup service stopped")

	return nil
}

func (c *CelestiaRollup) Stopped() bool {
	return c.stopped.Load()

}

func NewCelestiaRollup(cliCtx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {

	cfg, err := config.NewConfig(cliCtx)
	if err != nil {
		return nil, err
	}
	CelestiaConfig := ReadCLIConfig(cliCtx)

	logger := log.NewLogger(log.AppOut(cliCtx), cfg.LogConfig).New("celestia")
	log.SetGlobalLogHandler(logger.GetHandler())

	return CelestiaServiceFromCLIConfig(cliCtx.Context, cfg, CelestiaConfig, logger)
}

func CelestiaServiceFromCLIConfig(ctx context.Context, cfg *config.CLIConfig, celestiaConfig CLIConfig, logger log.Logger) (*CelestiaRollup, error) {
	var c CelestiaRollup
	if err := c.initFromCLIConfig(ctx, cfg, celestiaConfig, logger); err != nil {
		return nil, errors.Join(err, c.Stop(ctx)) // try to clean up our failed initialization attempt
	}
	return &c, nil
}

func (c *CelestiaRollup) initFromCLIConfig(ctx context.Context, cfg *config.CLIConfig, celestiaConfig CLIConfig, logger log.Logger) error {

	c.Config = cfg
	c.CelestiaConfig = celestiaConfig
	c.Log = logger

	if err := c.initDA(celestiaConfig); err != nil {
		return err
	}

	c.initMetrics(cfg)

	if err := c.initTxManager(cfg, c.Log); err != nil {
		return err
	}
	if err := c.initMetricsServer(cfg); err != nil {
		return err
	}

	c.driverCtx = ctx

	return nil
}

func (c *CelestiaRollup) initDA(celestiaConfig CLIConfig) error {
	client, err := NewDAClient(celestiaConfig.DaRpc, celestiaConfig.AuthToken, celestiaConfig.Namespace)
	if err != nil {
		return err
	}
	c.DAClient = client
	return nil
}

func (c *CelestiaRollup) initTxManager(cfg *config.CLIConfig, logger log.Logger) error {
	txManager, err := txmgr.NewSimpleTxManager("celestia-rollup", logger, c.Metrics, cfg.TxMgrConfig)
	if err != nil {
		return err
	}
	c.Txmgr = txManager
	return nil
}

func (c *CelestiaRollup) initMetrics(cfg *config.CLIConfig) {
	if cfg.MetricsConfig.Enabled {
		procName := "default"
		c.Metrics = metrics.NewMetrics(procName)
	}
}

func (c *CelestiaRollup) initMetricsServer(cfg *config.CLIConfig) error {
	if !cfg.MetricsConfig.Enabled {
		c.Log.Info("metrics disabled")
		return nil
	}
	m, ok := c.Metrics.(metrics.RegistryMetricer)
	if !ok {
		return fmt.Errorf("metrics were enabled, but metricer %T does not expose registry for metrics-server", c.Metrics)
	}
	c.Log.Debug("starting metrics server", "addr", cfg.MetricsConfig.ListenAddr, "port", cfg.MetricsConfig.ListenPort)
	addr, err := metrics2.StartServer(m.Registry(), cfg.MetricsConfig.ListenAddr, cfg.MetricsConfig.ListenPort)
	if err != nil {
		return fmt.Errorf("failed to start metrics server: %w", err)
	}
	c.Log.Info("started metrics server", "addr", addr)
	return nil
}

// sendTransaction creates & submits a transaction to the batch inbox address with the given `txData`.
// It currently uses the underlying `txmgr` to handle transaction sending & price management.
// This is a blocking method. It should not be called concurrently.
func (c *CelestiaRollup) sendTransaction(tx types.Transaction, queue *txmgr.Queue[types.Transaction], receiptsCh chan txmgr.TxReceipt[types.Transaction]) error {
	// Do the gas estimation offline. A value of 0 will cause the [txmgr] to estimate the gas limit.
	data := tx.Data()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Duration(c.Config.BlockTime)*time.Second)
	ids, err := c.DAClient.Client.Submit(ctx, [][]byte{data}, -1, c.DAClient.Namespace)
	cancel()
	if err == nil && len(ids) == 1 {
		c.Log.Info("celestia: blob successfully submitted", "id", hex.EncodeToString(ids[0]))
		data = append([]byte{DerivationVersionCelestia}, ids[0]...)
	} else {
		c.Log.Info("celestia: blob submission failed; falling back to eth", "err", err)
	}

	intrinsicGas, err := core.IntrinsicGas(data, nil, false, true, true, false)
	if err != nil {
		c.Log.Error("Failed to calculate intrinsic gas", "error", err)
		return err
	}

	candidate := txmgr.TxCandidate{
		To:       &c.CelestiaConfig.DSConfig.batchInboxAddress,
		TxData:   data,
		GasLimit: intrinsicGas,
	}

	queue.Send(tx, candidate, receiptsCh)
	return nil
}

func (c *CelestiaRollup) DataFromEVMTransactions(block types.Block) ([]eth.Data, error) {
	var out []eth.Data
	for _, tx := range block.Transactions() {
		if to := tx.To(); to != nil && *to == c.CelestiaConfig.DSConfig.batchInboxAddress {
			if isValidBatchTx(tx, c.CelestiaConfig.DSConfig.l1Signer, c.CelestiaConfig.DSConfig.batchInboxAddress, c.CelestiaConfig.DSConfig.batcherAddr) {
				data := tx.Data()
				switch len(data) {
				case 0:
					out = append(out, data)
				default:
					switch data[0] {
					case DerivationVersionCelestia:
						log.Info("celestia: blob request", "id", hex.EncodeToString(tx.Data()))
						ctx, cancel := context.WithTimeout(context.Background(), 30*time.Duration(c.Config.BlockTime)*time.Second)
						blobs, err := c.DAClient.Client.Get(ctx, [][]byte{data[1:]}, c.DAClient.Namespace)
						cancel()
						if err != nil {
							return nil, fmt.Errorf("celestia: failed to resolve frame: %w", err)
						}
						if len(blobs) != 1 {
							log.Warn("celestia: unexpected length for blobs", "expected", 1, "got", len(blobs))
							if len(blobs) == 0 {
								log.Warn("celestia: skipping empty blobs")
								continue
							}
						}
						out = append(out, blobs[0])
					default:
						out = append(out, data)
						log.Info("celestia: using eth fallback")
					}
				}
			}

		}
	}

	return out, nil
}
