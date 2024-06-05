package config

import (
	"github.com/urfave/cli/v2"

	"github.com/eniac-x-labs/rollup-node/log"
	"github.com/eniac-x-labs/rollup-node/metrics"
)

var defaultBlockTime = uint64(2)

type CLIConfig struct {
	BlockTime     uint64 `json:"block_time"`
	MetricsConfig metrics.CLIConfig
	LogConfig     log.CLIConfig
}

// NewConfig parses the Config from the provided flags or environment variables.
func NewConfig(ctx *cli.Context) (*CLIConfig, error) {

	return &CLIConfig{
		BlockTime:     defaultBlockTime,
		MetricsConfig: metrics.ReadCLIConfig(ctx),
		LogConfig:     log.ReadCLIConfig(ctx),
	}, nil
}
