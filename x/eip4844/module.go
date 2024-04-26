package eip4844

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/eniac-x-labs/rollup-node/client"
	"github.com/eniac-x-labs/rollup-node/common/cliapp"
	"github.com/eniac-x-labs/rollup-node/config"
	eth "github.com/eniac-x-labs/rollup-node/eth-serivce"
	"github.com/eniac-x-labs/rollup-node/log"
	metrics2 "github.com/eniac-x-labs/rollup-node/metrics"
	"github.com/eniac-x-labs/rollup-node/txmgr"
	"github.com/eniac-x-labs/rollup-node/txmgr/metrics"
)

var ErrAlreadyStopped = errors.New("already stopped")

type Eip4844Rollup struct {
	Eip4844Config  CLIConfig
	Txmgr          txmgr.TxManager
	Config         *config.CLIConfig
	Metrics        metrics.Metricer
	l1BeaconClient *eth.L1BeaconClient
	Log            log.Logger
	stopped        atomic.Bool
	driverCtx      context.Context
}

func (e *Eip4844Rollup) Start(ctx context.Context) error {
	client, _ := ethclient.DialContext(context.Background(), "https://eth-mainnet.g.alchemy.com/v2/-EK96JwUb8C_l_EfZqaLumlZG6PV8SDq")

	block, _ := client.BlockByNumber(context.Background(), new(big.Int).SetUint64(19560971))

	data, err := e.DataFromEVMTransactions(block)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(len(data))

	return nil
}

func (e *Eip4844Rollup) Stop(ctx context.Context) error {
	if e.stopped.Load() {
		return ErrAlreadyStopped
	}

	e.Log.Info("Stopping eip4844 rollup service")

	if e.Txmgr != nil {
		e.Txmgr.Close()
	}

	e.stopped.Store(true)
	e.Log.Info("eip4844 rollup service stopped")

	return nil
}

func (e *Eip4844Rollup) Stopped() bool {
	return e.stopped.Load()
}

func NewEip4844Rollup(cliCtx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {

	cfg, err := config.NewConfig(cliCtx)
	if err != nil {
		return nil, err
	}
	eip4844Config := ReadCLIConfig(cliCtx)

	logger := log.NewLogger(log.AppOut(cliCtx), cfg.LogConfig).New("eip-4844")
	log.SetGlobalLogHandler(logger.GetHandler())

	return Eip4844ServiceFromCLIConfig(cliCtx.Context, cfg, eip4844Config, logger)
}

func Eip4844ServiceFromCLIConfig(ctx context.Context, cfg *config.CLIConfig, eip4844Config CLIConfig, logger log.Logger) (*Eip4844Rollup, error) {
	var e Eip4844Rollup
	if err := e.initFromCLIConfig(ctx, cfg, eip4844Config, logger); err != nil {
		return nil, errors.Join(err, e.Stop(ctx)) // try to clean up our failed initialization attempt
	}
	return &e, nil
}

func (e *Eip4844Rollup) initFromCLIConfig(ctx context.Context, cfg *config.CLIConfig, eip4844Config CLIConfig, logger log.Logger) error {

	e.Config = cfg
	e.Eip4844Config = eip4844Config
	e.Log = logger

	bCl := client.NewBasicHTTPClient(eip4844Config.L1BeaconAddr)
	beaconCfg := eth.L1BeaconClientConfig{
		FetchAllSidecars: eip4844Config.ShouldFetchAllSidecars,
	}
	e.l1BeaconClient = eth.NewL1BeaconClient(bCl, beaconCfg)

	e.initMetrics(cfg)

	if err := e.initTxManager(cfg, e.Log); err != nil {
		return err
	}
	if err := e.initMetricsServer(cfg); err != nil {
		return err
	}

	e.driverCtx = ctx

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
	if e.Eip4844Config.UseBlobs {
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
		To:    &e.Eip4844Config.DSConfig.batchInboxAddress,
		Blobs: []*eth.Blob{&b},
	}, nil
}

func (e *Eip4844Rollup) calldataTxCandidate(data []byte) *txmgr.TxCandidate {
	return &txmgr.TxCandidate{
		To:     &e.Eip4844Config.DSConfig.batchInboxAddress,
		TxData: data,
	}
}

func (e *Eip4844Rollup) DataFromEVMTransactions(block *types.Block) (datas []eth.Data, err error) {

	_, hashes := dataAndHashesFromTxs(block.Transactions(), e.Eip4844Config.DSConfig, e.Log)
	if len(hashes) == 0 {
		// there are no blobs to fetch so we can return immediately
		return nil, nil
	}

	ref := eth.L1BlockRef{
		Hash:       block.Hash(),
		Number:     block.NumberU64(),
		ParentHash: block.ParentHash(),
		Time:       block.Time(),
	}
	blobs, err := e.l1BeaconClient.GetBlobs(e.driverCtx, ref, hashes)
	if errors.Is(err, ethereum.NotFound) {
		// If the L1 block was available, then the blobs should be available too. The only
		// exception is if the blob retention window has expired, which we will ultimately handle
		// by failing over to a blob archival service.
		return nil, fmt.Errorf("failed to fetch blobs: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("failed to fetch blobs: %w", err)
	}

	for _, blob := range blobs {
		data, err := blob.ToData()
		if err != nil {
			return nil, fmt.Errorf("decodes the blob into raw byte data failed: %w", err)
		}

		datas = append(datas, data)
	}

	if err != nil {
		e.Log.Error("ignoring blob due to parse failure", "err", err)
		return nil, err
	}

	return datas, nil
}
