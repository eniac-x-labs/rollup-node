package celestia

import (
	"errors"
	"fmt"
	"math/big"
	"net"

	"github.com/urfave/cli/v2"

	eth "github.com/eniac-x-labs/rollup-node/eth-serivce"
)

const (
	DaRpcFlagName = "celestia.da.rpc"
	// AuthTokenFlagName defines the flag for the auth token
	AuthTokenFlagName = "celestia.da.auth_token"
	// NamespaceFlagName defines the flag for the namespace
	NamespaceFlagName = "celestia.da.namespace"
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

	return CLIConfig{
		DaRpc:     ctx.String(DaRpcFlagName),
		AuthToken: ctx.String(AuthTokenFlagName),
		Namespace: ctx.String(NamespaceFlagName),
	}
}
