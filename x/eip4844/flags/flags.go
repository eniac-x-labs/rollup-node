package flags

import (
	"fmt"
	"github.com/eniac-x-labs/rollup-node/x/eip4844/log"

	"github.com/urfave/cli/v2"

	"github.com/eniac-x-labs/rollup-node/metrics"
	service "github.com/eniac-x-labs/rollup-node/x/eip4844/eth-serivce"
	"github.com/eniac-x-labs/rollup-node/x/eip4844/txmgr"
)

const EnvVarPrefix = "DAPP-ROLLUP"

func prefixEnvVars(name string) []string {
	return service.PrefixEnvVar(EnvVarPrefix, name)
}

var (
	DataAvailabilityTypeFlag = &cli.BoolFlag{
		Name:    "data-availability-type",
		Usage:   "The data availability type to use for submitting batches to the L1.",
		Value:   false,
		EnvVars: prefixEnvVars("DATA_AVAILABILITY_TYPE"),
	}
)

var requiredFlags = []cli.Flag{
	DataAvailabilityTypeFlag,
}

var optionalFlags = []cli.Flag{}

func init() {
	optionalFlags = append(optionalFlags, log.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, metrics.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, txmgr.CLIFlags(EnvVarPrefix)...)

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
