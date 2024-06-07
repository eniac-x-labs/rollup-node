package cli_config

import (
	"github.com/eniac-x-labs/rollup-node/metrics"
	"github.com/urfave/cli/v2"
)

var defaultBlockTime = uint64(2)

type CLIConfig struct {
	BlockTime     uint64 `json:"block_time"`
	MetricsConfig metrics.CLIConfig
}

// NewConfig parses the Config from the provided flags or environment variables.
func NewConfig(ctx *cli.Context) (*CLIConfig, error) {

	return &CLIConfig{
		BlockTime:     defaultBlockTime,
		MetricsConfig: metrics.ReadCLIConfig(ctx),
	}, nil
}

func ParseCLIConfig(blockTime uint64, enableMetrics bool, listenAddr string, listenPort int) *CLIConfig {
	return &CLIConfig{
		BlockTime: blockTime,
		MetricsConfig: metrics.CLIConfig{
			Enabled:    enableMetrics,
			ListenAddr: listenAddr,
			ListenPort: listenPort,
		},
	}
}
