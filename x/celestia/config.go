package celestia

import (
	"github.com/ethereum/go-ethereum/log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

	L1ChainIdFlagName uint64 `toml:"l1ChainIdFlagName"`
}

type CelestiaConfig struct {
	//cfg            *config.CLIConfig
	celestiaConfig CLIConfig
	logger         log.Logger
}

func ProcessCelestiaConfig(parseConf *ParseCelestiaConfig, logger log.Logger) (*CelestiaConfig, error) {
	l1ChainID, _ := new(big.Int).SetString(parseConf.L1ChainID, 10)
	signer := types.NewCancunSigner(new(big.Int).SetUint64(parseConf.L1ChainIdFlagName))

	return &CelestiaConfig{
		celestiaConfig: CLIConfig{
			L1Rpc:               parseConf.L1Rpc,
			L1ChainID:           l1ChainID,
			PrivateKey:          parseConf.PrivateKey,
			DaRpc:               parseConf.DaRpc,
			AuthToken:           parseConf.AuthToken,
			Namespace:           parseConf.Namespace,
			EthFallbackDisabled: parseConf.EthFallbackDisabled,
			DSConfig: &DataSourceConfig{
				l1Signer:          signer,
				batchInboxAddress: common.HexToAddress(parseConf.BatchInboxAddress),
				batcherAddr:       common.HexToAddress(parseConf.BatcherAddr),
			},
		},
		logger: logger,
	}, nil
}
