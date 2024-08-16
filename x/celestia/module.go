package celestia

import (
	"context"
	"errors"
	"sync/atomic"

	client "github.com/celestiaorg/celestia-openrpc"
	"github.com/celestiaorg/celestia-openrpc/types/blob"
	"github.com/celestiaorg/celestia-openrpc/types/share"
	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/log"

	cli_config "github.com/eniac-x-labs/rollup-node/config/cli-config"
)

var ErrAlreadyStopped = errors.New("already stopped")

type CelestiaRollup struct {
	CelestiaConfig CLIConfig
	Config         *cli_config.CLIConfig
	Log            log.Logger
	DAClient       *client.Client
	Namespace      share.Namespace
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
	if err := c.initFromCLIConfig(ctx, celestiaConfig, logger); err != nil {
		return nil, errors.Join(err, c.Stop(ctx)) // try to clean up our failed initialization attempt
	}
	return &c, nil
}

func NewCelestiaRollupWithConfig(ctx context.Context, config *CelestiaConfig) (*CelestiaRollup, error) {
	if config == nil {
		log.Error("celestia config is nil pointer")
		return nil, nil
	}

	var c CelestiaRollup
	if err := c.initFromCLIConfig(ctx, config.celestiaConfig, config.logger); err != nil {
		return nil, errors.Join(err, c.Stop(ctx)) // try to clean up our failed initialization attempt
	}
	return &c, nil
}

func (c *CelestiaRollup) initFromCLIConfig(ctx context.Context, celestiaConfig CLIConfig, logger log.Logger) error {

	c.CelestiaConfig = celestiaConfig
	c.Log = logger

	if err := c.initDA(ctx, celestiaConfig); err != nil {
		return err
	}

	return nil
}

func (c *CelestiaRollup) initDA(ctx context.Context, celestiaConfig CLIConfig) error {
	namespace, err := share.NewBlobNamespaceV0([]byte{0xDE, 0xAD, 0xBE, 0xEF})
	if err != nil {
		return err
	}

	client, err := client.NewClient(ctx, celestiaConfig.DaRpc, celestiaConfig.AuthToken)
	if err != nil {
		return err
	}

	c.Namespace = namespace
	c.DAClient = client

	return nil
}

func (c *CelestiaRollup) SubmitBlob(ctx context.Context, data []byte) (uint64, error) {

	blobData, err := blob.NewBlobV0(c.Namespace, data)
	if err != nil {
		return 0, err
	}

	// submit the blob to the network
	height, err := c.DAClient.Blob.Submit(ctx, []*blob.Blob{blobData}, blob.DefaultGasPrice())
	if err != nil {
		return 0, err
	}

	return height, nil
}

func (c *CelestiaRollup) RetrievedBlobs(ctx context.Context, height uint64) ([]byte, error) {
	// fetch the blob back from the network
	retrievedBlobs, err := c.DAClient.Blob.GetAll(ctx, height, []share.Namespace{c.Namespace})
	if err != nil {
		return nil, err
	}

	return retrievedBlobs[0].Data, nil
}
