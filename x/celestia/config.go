package celestia

import (
	_log "github.com/eniac-x-labs/rollup-node/log"
)

// corresponding celestia.toml
type ParseCelestiaConfig struct {
	L1Rpc               string `toml:"l1Rpc"`
	L1ChainID           string `toml:"l1ChainID"`
	PrivateKey          string `toml:"privateKey"`
	DaRpc               string `toml:"daRpc"`
	AuthToken           string `toml:"authToken"`
	Namespace           string `toml:"namespace"`
	EthFallbackDisabled bool   `toml:"ethFallbackDisabled"`

	// data source config
	BatchInboxAddress string `toml:"batchInboxAddress"`
	BatcherAddr       string `toml:"batcherAddr"`

	// BlockTime CLI Config
	BlockTime uint64 `toml:"blockTime"`

	// Metrics CLI Config
	Enable     bool   `toml:"enable"`
	ListenAddr string `toml:"listenAddr"`
	ListenPort int    `toml:"listenPort"`

	// Log CLI Config
	Level      int    `toml:"level"`
	Color      bool   `toml:"color"`
	FormatType string `toml:"formatType"`
}

type CelestiaConfig struct {
	//cfg            *config.CLIConfig
	celestiaConfig CLIConfig
	logger         _log.Logger
}

func ProcessCelestiaConfig(parseConf *ParseCelestiaConfig, logger _log.Logger) (*CelestiaConfig, error) {
	// todo: kezi
	return nil, nil
}
