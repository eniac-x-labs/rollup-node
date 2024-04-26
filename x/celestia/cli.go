package celestia

import (
	"errors"
	"fmt"
	"math/big"
	"net"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	eth "github.com/eniac-x-labs/rollup-node/eth-serivce"
)

const (
	L1ChainIdFlagName = "l1.chain-id"
	DaRpcFlagName     = "celestia.da.rpc"
	// AuthTokenFlagName defines the flag for the auth token
	AuthTokenFlagName = "celestia.da.auth_token"
	// NamespaceFlagName defines the flag for the namespace
	NamespaceFlagName         = "celestia.da.namespace"
	BatcherAddressFlagName    = "celestia..batcher-address"
	BatchInboxAddressFlagName = "celestia..batch-inbox-address"
)

var (
	defaultDaRpc = "localhost:26650"

	ErrInvalidPort = errors.New("invalid port")
)

func Check(address string) error {
	_, port, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}

	if port == "" {
		return ErrInvalidPort
	}

	_, err = net.LookupPort("tcp", port)
	if err != nil {
		return err
	}

	return nil
}

func CLIFlags(envPrefix string) []cli.Flag {
	return []cli.Flag{
		&cli.Uint64Flag{
			Name:     L1ChainIdFlagName,
			Usage:    "The chain id of l1.",
			Required: false,
			EnvVars:  eth.PrefixEnvVar(envPrefix, "L1_CHAIN_ID"),
		},
		&cli.StringFlag{
			Name:    DaRpcFlagName,
			Usage:   "dial address of data availability grpc client",
			Value:   defaultDaRpc,
			EnvVars: eth.PrefixEnvVar(envPrefix, "CELESTIA_DA_RPC"),
		},
		&cli.StringFlag{
			Name:    AuthTokenFlagName,
			Usage:   "authentication token of the data availability client",
			EnvVars: eth.PrefixEnvVar(envPrefix, "CELESTIA_DA_AUTH_TOKEN"),
		},
		&cli.StringFlag{
			Name:    NamespaceFlagName,
			Usage:   "namespace of the data availability client",
			EnvVars: eth.PrefixEnvVar(envPrefix, "CELESTIA_DA_NAMESPACE"),
		},
		&cli.StringFlag{
			Name:     BatcherAddressFlagName,
			Usage:    "Address of celestia Batcher.",
			Required: false,
			EnvVars:  eth.PrefixEnvVar(envPrefix, "CELESTIA_BATCHER_ADDRESS"),
		},
		&cli.StringFlag{
			Name:     BatchInboxAddressFlagName,
			Usage:    "Address of celestia Batch inbox.",
			Required: false,
			EnvVars:  eth.PrefixEnvVar(envPrefix, "CELESTIA_BATCH_INBOX_ADDRESS"),
		},
	}
}

type Config struct {
	DaRpc string
}

func (c Config) Check() error {
	if c.DaRpc == "" {
		c.DaRpc = defaultDaRpc
	}

	if err := Check(c.DaRpc); err != nil {
		return fmt.Errorf("invalid da rpc: %w", err)
	}

	return nil
}

type CLIConfig struct {
	DaRpc     string
	AuthToken string
	Namespace string
	DSConfig  *DataSourceConfig
}

func (c CLIConfig) Check() error {
	if c.DaRpc == "" {
		c.DaRpc = defaultDaRpc
	}

	if err := Check(c.DaRpc); err != nil {
		return fmt.Errorf("invalid da rpc: %w", err)
	}

	return nil
}

func NewCLIConfig() CLIConfig {
	return CLIConfig{
		DaRpc: defaultDaRpc,
	}
}

func ReadCLIConfig(ctx *cli.Context) CLIConfig {

	signer := types.NewCancunSigner(new(big.Int).SetUint64(ctx.Uint64(L1ChainIdFlagName)))

	dsConfig := DataSourceConfig{
		l1Signer:          signer,
		batchInboxAddress: common.HexToAddress(ctx.String(BatchInboxAddressFlagName)),
		batcherAddr:       common.HexToAddress(ctx.String(BatcherAddressFlagName)),
	}
	return CLIConfig{
		DaRpc:     ctx.String(DaRpcFlagName),
		AuthToken: ctx.String(AuthTokenFlagName),
		Namespace: ctx.String(NamespaceFlagName),
		DSConfig:  &dsConfig,
	}
}
