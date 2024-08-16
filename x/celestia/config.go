package celestia

import (
	"github.com/ethereum/go-ethereum/log"
)

// corresponding celestia.toml
type ParseCelestiaConfig struct {
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

	L1ChainIdFlagName uint64 `toml:"l1ChainIdFlagName"`
}

type CelestiaConfig struct {
	//cfg            *config.CLIConfig
	celestiaConfig CLIConfig
	logger         log.Logger
}

func ProcessCelestiaConfig(parseConf *ParseCelestiaConfig, logger log.Logger) (*CelestiaConfig, error) {

	return &CelestiaConfig{
		celestiaConfig: CLIConfig{
			DaRpc:     parseConf.DaRpc,
			AuthToken: parseConf.AuthToken,
			Namespace: parseConf.Namespace,
		},
		logger: logger,
	}, nil
}
