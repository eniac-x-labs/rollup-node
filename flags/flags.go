package flags

import (
	"fmt"

	"github.com/urfave/cli/v2"

	service "github.com/eniac-x-labs/rollup-node/eth-serivce"
	"github.com/eniac-x-labs/rollup-node/log"
	"github.com/eniac-x-labs/rollup-node/metrics"
	"github.com/eniac-x-labs/rollup-node/txmgr"
	"github.com/eniac-x-labs/rollup-node/x/celestia"
	"github.com/eniac-x-labs/rollup-node/x/eip4844"
)

const EnvVarPrefix = "DAPP_ROLLUP"

func prefixEnvVars(name string) []string {
	return service.PrefixEnvVar(EnvVarPrefix, name)
}

var requiredFlags = []cli.Flag{}

var optionalFlags = []cli.Flag{}

func init() {
	optionalFlags = append(optionalFlags, log.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, metrics.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, txmgr.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, eip4844.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, celestia.CLIFlags(EnvVarPrefix)...)

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
