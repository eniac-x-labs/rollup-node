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
	defaultBlockTime = uint64(2)

	DaRpcFlagName = "celestia.da.rpc"
	// AuthTokenFlagName defines the flag for the auth token
	AuthTokenFlagName = "celestia.da.auth_token"
	// NamespaceFlagName defines the flag for the namespace
	NamespaceFlagName           = "celestia.da.namespace"
	BatcherAddressFlagName      = "celestia.batcher-address"
	BatchInboxAddressFlagName   = "celestia.batch-inbox-address"
	EthFallbackDisabledFlagName = "celestia.eth_fallback_disabled"
	BlockTime                   = "celestia.block-time"
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
		&cli.BoolFlag{
			Name:    EthFallbackDisabledFlagName,
			Usage:   "disable eth fallback",
			EnvVars: eth.PrefixEnvVar(envPrefix, "CELESTIA_DA_ETH_FALLBACK_DISABLED"),
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
		&cli.Uint64Flag{
			Name:     BlockTime,
			Usage:    "block time of celestia.",
			Required: false,
			Value:    defaultBlockTime,
			EnvVars:  eth.PrefixEnvVar(envPrefix, "CELESTIA_BLOCK_TIME"),
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
	BlockTime           uint64
	DaRpc               string
	AuthToken           string
	Namespace           string
	EthFallbackDisabled bool
	DSConfig            *DataSourceConfig
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

func ReadCLIConfig(ctx *cli.Context, l1ChainId *big.Int) CLIConfig {

	signer := types.NewCancunSigner(l1ChainId)

	dsConfig := DataSourceConfig{
		l1Signer:          signer,
		batchInboxAddress: common.HexToAddress(ctx.String(BatchInboxAddressFlagName)),
		batcherAddr:       common.HexToAddress(ctx.String(BatcherAddressFlagName)),
	}
	return CLIConfig{
		BlockTime:           defaultBlockTime,
		DaRpc:               ctx.String(DaRpcFlagName),
		AuthToken:           ctx.String(AuthTokenFlagName),
		Namespace:           ctx.String(NamespaceFlagName),
		EthFallbackDisabled: ctx.Bool(EthFallbackDisabledFlagName),
		DSConfig:            &dsConfig,
	}
}
