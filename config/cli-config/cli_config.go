package cli_config

import (
	"github.com/eniac-x-labs/rollup-node/log"
	"github.com/eniac-x-labs/rollup-node/metrics"
	"github.com/urfave/cli/v2"
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

func ParseCLIConfig(blockTime uint64, enableMetrics bool, listenAddr string, listenPort int, logLevel int, color bool, formatType string) *CLIConfig {
	return &CLIConfig{
		BlockTime: blockTime,
		MetricsConfig: metrics.CLIConfig{
			Enabled:    enableMetrics,
			ListenAddr: listenAddr,
			ListenPort: listenPort,
		},
		LogConfig: log.CLIConfig{
			Level:  log.Lvl(logLevel),
			Color:  color,
			Format: log.FormatType(formatType),
		},
	}
}
