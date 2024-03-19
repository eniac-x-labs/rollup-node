package config

import (
	"fmt"
	"github.com/eniac-x-labs/rollup-node/x/eip4844/log"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/common"

	"github.com/eniac-x-labs/rollup-node/metrics"
	"github.com/eniac-x-labs/rollup-node/x/eip4844/flags"
	"github.com/eniac-x-labs/rollup-node/x/eip4844/txmgr"
)

type CLIConfig struct {
	NetworkTimeout         time.Duration
	MaxPendingTransactions uint64
	BatchInboxAddress      common.Address
	UseBlobs               bool
	TxMgrConfig            txmgr.CLIConfig
	MetricsConfig          metrics.CLIConfig
	LogConfig              log.CLIConfig
	DataAvailabilityType   string
}

// NewConfig parses the Config from the provided flags or environment variables.
func NewConfig(ctx *cli.Context) (*CLIConfig, error) {
	var useBlobs bool
	switch ctx.String(flags.DataAvailabilityTypeFlag.Name) {
	case "blobs":
		useBlobs = true
	case "calldata":
		useBlobs = false
	default:
		return nil, fmt.Errorf("unknown data availability type: %v", ctx.String(flags.DataAvailabilityTypeFlag.Name))
	}
	return &CLIConfig{
		UseBlobs:      useBlobs,
		TxMgrConfig:   txmgr.ReadCLIConfig(ctx),
		MetricsConfig: metrics.ReadCLIConfig(ctx),
		LogConfig:     log.ReadCLIConfig(ctx),
	}, nil
}
