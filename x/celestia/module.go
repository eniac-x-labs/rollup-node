package celestia

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"sync/atomic"
	"time"

	cli_config "github.com/eniac-x-labs/rollup-node/config/cli-config"
	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/eniac-x-labs/rollup-node/client"
	eth "github.com/eniac-x-labs/rollup-node/eth-serivce"
	"github.com/eniac-x-labs/rollup-node/signer"
)

var ErrAlreadyStopped = errors.New("already stopped")

type CelestiaRollup struct {
	CelestiaConfig CLIConfig
	Config         *cli_config.CLIConfig
	Log            log.Logger
	DAClient       *DAClient
	ethClients     client.EthClient
	Signer         signer.SignerFn
	From           common.Address
	stopped        atomic.Bool
}

func (c *CelestiaRollup) Stop(ctx context.Context) error {
	if c.stopped.Load() {
		return ErrAlreadyStopped
	}

	c.Log.Info("Stopping Celestia rollup service")

	c.stopped.Store(true)
	c.Log.Info("Celestia rollup service stopped")

	return nil
}

func (c *CelestiaRollup) Stopped() bool {
	return c.stopped.Load()

}

func NewCelestiaRollup(cliCtx *cli.Context, logger log.Logger) (*CelestiaRollup, error) {
	cfg, err := cli_config.NewConfig(cliCtx)
	if err != nil {
		return nil, err
	}
	return CelestiaServiceFromCLIConfig(cliCtx.Context, cfg, ReadCLIConfig(cliCtx, cfg.L1ChainID), logger)
}

func CelestiaServiceFromCLIConfig(ctx context.Context, cfg *cli_config.CLIConfig, celestiaConfig CLIConfig, logger log.Logger) (*CelestiaRollup, error) {
	var c CelestiaRollup
	if err := c.initFromCLIConfig(ctx, cfg, celestiaConfig, logger); err != nil {
		return nil, errors.Join(err, c.Stop(ctx)) // try to clean up our failed initialization attempt
	}
	return &c, nil
}

func NewCelestiaRollupWithConfig(ctx context.Context, cfg *cli_config.CLIConfig, config *CelestiaConfig) (*CelestiaRollup, error) {
	if cfg == nil || config == nil {
		log.Error("celestia config is nil pointer")
		return nil, nil
	}

	var c CelestiaRollup
	if err := c.initFromCLIConfig(ctx, cfg, config.celestiaConfig, config.logger); err != nil {
		return nil, errors.Join(err, c.Stop(ctx)) // try to clean up our failed initialization attempt
	}
	return &c, nil
}

func (c *CelestiaRollup) initFromCLIConfig(ctx context.Context, cfg *cli_config.CLIConfig, celestiaConfig CLIConfig, logger log.Logger) error {

	c.Config = cfg
	c.CelestiaConfig = celestiaConfig
	c.Log = logger

	l1Client, err := client.DialEthClient(ctx, cfg.L1Rpc)
	if err != nil {
		log.Error("failed to dial eth client", "err", err)
		return err
	}
	c.ethClients = l1Client

	signerFactory, from, err := signer.SignerFactoryFromPrivateKey(cfg.PrivateKey)
	if err != nil {
		log.Error(fmt.Errorf("could not init signer: %w", err).Error())
		return err
	}
	c.Signer = signerFactory(cfg.L1ChainID)
	c.From = from

	if err := c.initDA(celestiaConfig); err != nil {
		return err
	}

	return nil
}

func (c *CelestiaRollup) initDA(celestiaConfig CLIConfig) error {
	client, err := NewDAClient(celestiaConfig.DaRpc, celestiaConfig.AuthToken, celestiaConfig.Namespace, celestiaConfig.EthFallbackDisabled)
	if err != nil {
		return err
	}
	c.DAClient = client
	return nil
}

// SendTransaction creates & submits a transaction to the batch inbox address with the given `txData`.
// It currently uses the underlying `txmgr` to handle transaction sending & price management.
// This is a blocking method. It should not be called concurrently.
func (c *CelestiaRollup) SendTransaction(ctx context.Context, data []byte) ([]byte, error) {

	candidate, err := c.calldataTxCandidate(data)
	if err != nil {
		c.Log.Error("building Calldata transaction candidate", "err", err)
		return nil, err
	}

	intrinsicGas, err := core.IntrinsicGas(data, nil, false, true, true, false)
	if err != nil {
		c.Log.Error("Failed to calculate intrinsic gas", "error", err)
		return nil, err
	}
	candidate.GasLimit = intrinsicGas

	tx, err := c.craftTx(ctx, *candidate)
	if err != nil {
		c.Log.Error("Failed to create a transaction", "err", err)
		return nil, err
	}

	signTx, err := c.Signer(ctx, c.From, tx)
	if err != nil {
		c.Log.Error("Failed to sign a transaction", "err", err)
		return nil, err
	}

	err = c.ethClients.SendTransaction(ctx, signTx)
	if err != nil {
		c.Log.Error("Failed to send transaction", "err", err)
		return nil, err
	}
	return signTx.Hash().Bytes(), nil
}

func (c *CelestiaRollup) calldataTxCandidate(data []byte) (*eth.TxCandidate, error) {
	c.Log.Info("building Calldata transaction candidate", "size", len(data))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Duration(c.CelestiaConfig.BlockTime)*time.Second)
	ids, err := c.DAClient.Client.Submit(ctx, [][]byte{data}, -1, c.DAClient.Namespace)
	fmt.Println(ids)
	fmt.Println(err)
	cancel()
	if err == nil && len(ids) == 1 {
		c.Log.Info("celestia: blob successfully submitted", "id", hex.EncodeToString(ids[0]))
		data = append([]byte{DerivationVersionCelestia}, ids[0]...)
	} else {
		if c.DAClient.EthFallbackDisabled {
			return nil, fmt.Errorf("celestia: blob submission failed; eth fallback disabled: %w", err)
		}

		c.Log.Info("celestia: blob submission failed; falling back to eth", "err", err)
	}
	return &eth.TxCandidate{
		To:     &c.CelestiaConfig.DSConfig.batchInboxAddress,
		TxData: data,
	}, nil
}

func (c *CelestiaRollup) DataFromEVMTransactions(ctx context.Context, txHashStr string) (eth.Data, error) {
	var out eth.Data

	tx, err := c.ethClients.TxByHash(common.HexToHash(txHashStr))
	if err != nil {
		c.Log.Error("failed to get transaction", "tx_hash", txHashStr)
		return nil, err
	}

	if to := tx.To(); to != nil && *to == c.CelestiaConfig.DSConfig.batchInboxAddress {
		if isValidBatchTx(tx, c.CelestiaConfig.DSConfig.l1Signer, c.CelestiaConfig.DSConfig.batchInboxAddress, c.CelestiaConfig.DSConfig.batcherAddr) {
			data := tx.Data()
			switch len(data) {
			case 0:
				out = data
			default:
				switch data[0] {
				case DerivationVersionCelestia:
					log.Info("celestia: blob request", "id", hex.EncodeToString(tx.Data()))
					ctxT, cancel := context.WithTimeout(ctx, 30*time.Duration(c.CelestiaConfig.BlockTime)*time.Second)
					blobs, err := c.DAClient.Client.Get(ctxT, [][]byte{data[1:]}, c.DAClient.Namespace)
					cancel()
					if err != nil {
						return nil, fmt.Errorf("celestia: failed to resolve frame: %w", err)
					}
					if len(blobs) != 1 {
						log.Warn("celestia: unexpected length for blobs", "expected", 1, "got", len(blobs))
						if len(blobs) == 0 {
							log.Warn("celestia: skipping empty blobs")
						}
					}
					out = blobs[0]
				default:
					out = data
					log.Info("celestia: using eth fallback")
				}
			}
		}

	}

	return out, nil
}

func (c *CelestiaRollup) craftTx(ctx context.Context, candidate eth.TxCandidate) (*types.Transaction, error) {
	c.Log.Debug("crafting Transaction", "blobs", len(candidate.Blobs), "calldata_size", len(candidate.TxData))

	tip, err := c.ethClients.SuggestGasTipCap(ctx)
	if err != nil {
		c.Log.Error(fmt.Errorf("failed to fetch the suggested gas tip cap: %w", err).Error())
		return nil, err
	}

	header, err := c.ethClients.HeaderByNumber(ctx, nil)
	if err != nil {
		c.Log.Error(fmt.Errorf("failed to fetch the suggested base fee: %w", err).Error())
		return nil, err
	}
	baseFee := header.BaseFee
	gasFeeCap := calcGasFeeCap(baseFee, tip)

	txMessage := &types.DynamicFeeTx{
		ChainID:   c.Config.L1ChainID,
		To:        candidate.To,
		GasTipCap: tip,
		GasFeeCap: gasFeeCap,
		Value:     candidate.Value,
		Data:      candidate.TxData,
		Gas:       candidate.GasLimit,
	}

	return types.NewTx(txMessage), err
}

func calcGasFeeCap(baseFee, gasTipCap *big.Int) *big.Int {
	return new(big.Int).Add(
		gasTipCap,
		new(big.Int).Mul(baseFee, big.NewInt(2)),
	)
}
