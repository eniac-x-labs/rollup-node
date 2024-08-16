package flags

import (
	"fmt"

	cli_config "github.com/eniac-x-labs/rollup-node/config/cli-config"
	"github.com/urfave/cli/v2"

	service "github.com/eniac-x-labs/rollup-node/eth-serivce"
	"github.com/eniac-x-labs/rollup-node/metrics"
	"github.com/eniac-x-labs/rollup-node/x/celestia"
	"github.com/eniac-x-labs/rollup-node/x/eip4844"
)

const EnvVarPrefix = "DAPP_ROLLUP"

func prefixEnvVars(name string) []string {
	return service.PrefixEnvVar(EnvVarPrefix, name)
}

var requiredFlags = []cli.Flag{}

var optionalFlags = []cli.Flag{}

var exposedAddress = []cli.Flag{
	&cli.StringFlag{
		Name:    "rpcAddress",
		Usage:   "Listen address for rpc and sdk",
		EnvVars: PrefixEnvVar(EnvVarPrefix, "RPC_ADDRESS"),
	},
	&cli.StringFlag{
		Name:    "apiAddress",
		Usage:   "Listen address for web server",
		EnvVars: PrefixEnvVar(EnvVarPrefix, "API_ADDRESS"),
	},
}

func PrefixEnvVar(prefix, suffix string) []string {
	return []string{prefix + "_" + suffix}
}

func init() {
	optionalFlags = append(optionalFlags, cli_config.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, metrics.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, eip4844.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, celestia.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, exposedAddress...)

	Flags = append(requiredFlags, optionalFlags...)
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func CheckRequired(ctx *cli.Context) error {
	for _, f := range requiredFlags {
		if !ctx.IsSet(f.Names()[0]) {
			return fmt.Errorf("flag %s is required", f.Names()[0])
		}
	}
	return nil
}
