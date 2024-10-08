package metrics

import (
	"errors"
	"math"

	"github.com/urfave/cli/v2"
)

const (
	EnabledFlagName    = "metrics.enabled"
	ListenAddrFlagName = "metrics.addr"
	PortFlagName       = "metrics.port"
	defaultListenAddr  = "0.0.0.0"
	defaultListenPort  = 7300
)

func DefaultCLIConfig() CLIConfig {
	return CLIConfig{
		Enabled:    false,
		ListenAddr: defaultListenAddr,
		ListenPort: defaultListenPort,
	}
}

func PrefixEnvVar(prefix, suffix string) []string {
	return []string{prefix + "_" + suffix}
}

func CLIFlags(envPrefix string) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    EnabledFlagName,
			Usage:   "Enable the metrics server",
			EnvVars: PrefixEnvVar(envPrefix, "METRICS_ENABLED"),
		},
		&cli.StringFlag{
			Name:    ListenAddrFlagName,
			Usage:   "Metrics listening address",
			Value:   defaultListenAddr, // TODO(CLI-4159): Switch to 127.0.0.1
			EnvVars: PrefixEnvVar(envPrefix, "METRICS_ADDR"),
		},
		&cli.IntFlag{
			Name:    PortFlagName,
			Usage:   "Metrics listening port",
			Value:   defaultListenPort,
			EnvVars: PrefixEnvVar(envPrefix, "METRICS_PORT"),
		},
	}
}

type CLIConfig struct {
	Enabled    bool   `toml:"enable"`
	ListenAddr string `toml:"listenAddr"`
	ListenPort int    `toml:"listenPort"`
}

func (m CLIConfig) Check() error {
	if !m.Enabled {
		return nil
	}

	if m.ListenPort < 0 || m.ListenPort > math.MaxUint16 {
		return errors.New("invalid metrics port")
	}

	return nil
}

func ReadCLIConfig(ctx *cli.Context) CLIConfig {
	return CLIConfig{
		Enabled:    ctx.Bool(EnabledFlagName),
		ListenAddr: ctx.String(ListenAddrFlagName),
		ListenPort: ctx.Int(PortFlagName),
	}
}
