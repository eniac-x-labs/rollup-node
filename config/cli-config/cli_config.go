package cli_config

import (
	eth "github.com/eniac-x-labs/rollup-node/eth-serivce"
	"github.com/urfave/cli/v2"
	"math/big"
)

var (
	L1RPCFlagName      = "l1-eth-rpc"
	PrivateKeyFlagName = "private-key"
	L1ChainIdFlagName  = "l1.chain-id"
)

type CLIConfig struct {
	L1Rpc      string
	L1ChainID  *big.Int
	PrivateKey string
}

func CLIFlags(envPrefix string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     L1RPCFlagName,
			Usage:    "The rpc url of l1.",
			Required: true,
			EnvVars:  eth.PrefixEnvVar(envPrefix, "L1_ETH_RPC"),
		},
		&cli.StringFlag{
			Name:     PrivateKeyFlagName,
			Usage:    "The private key to use with the service. Must not be used with mnemonic.",
			Required: true,
			EnvVars:  eth.PrefixEnvVar(envPrefix, "PRIVATE_KEY"),
		},
		&cli.Uint64Flag{
			Name:     L1ChainIdFlagName,
			Usage:    "The chain id of l1.",
			Required: false,
			EnvVars:  eth.PrefixEnvVar(envPrefix, "L1_CHAIN_ID"),
		},
	}
}

// NewConfig parses the Config from the provided flags or environment variables.
func NewConfig(ctx *cli.Context) (*CLIConfig, error) {
	return &CLIConfig{
		L1ChainID:  new(big.Int).SetUint64(ctx.Uint64(L1ChainIdFlagName)),
		PrivateKey: ctx.String(PrivateKeyFlagName),
		L1Rpc:      ctx.String(L1RPCFlagName),
	}, nil
}

func ParseCLIConfig(blockTime uint64) *CLIConfig {
	return &CLIConfig{}
}
