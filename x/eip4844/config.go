package eip4844

import (
	"math/big"

	_log "github.com/eniac-x-labs/rollup-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ParseEip4844Config struct {
	L1Rpc                  string `toml:"l1Rpc"`
	PrivateKey             string `toml:"privateKey"`
	L1ChainID              string `toml:"l1ChainID"` // *bigInt
	UseBlobs               bool   `toml:"useBlobs"`
	L1BeaconAddr           string `toml:"l1BeaconAddr"`
	ShouldFetchAllSidecars bool   `toml:"shouldFetchAllSidecars"`

	// data source config
	BatchInboxAddress string `toml:"batchInboxAddress"` // common.Address
	BatcherAddr       string `toml:"batcherAddr"`       // common.Address

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

	L1ChainIdFlagName uint64 `toml:"l1ChainIdFlagName"`
}

type Eip4844Config struct {
	//cfg           *config.CLIConfig
	eip4844Config CLIConfig
	logger        _log.Logger
}

func ProcessEip4844Config(parseConf *ParseEip4844Config, logger _log.Logger) (*Eip4844Config, error) {
	l1ChainID, _ := new(big.Int).SetString(parseConf.L1ChainID, 10)
	signer := types.NewCancunSigner(new(big.Int).SetUint64(parseConf.L1ChainIdFlagName))

	return &Eip4844Config{
		eip4844Config: CLIConfig{
			L1Rpc:      parseConf.L1Rpc,
			PrivateKey: parseConf.PrivateKey,
			L1ChainID:  l1ChainID,
			DSConfig: &DataSourceConfig{
				l1Signer:          signer,
				batchInboxAddress: common.HexToAddress(parseConf.BatchInboxAddress),
				batcherAddr:       common.HexToAddress(parseConf.BatcherAddr),
			},
			UseBlobs:               parseConf.UseBlobs,
			L1BeaconAddr:           parseConf.L1BeaconAddr,
			ShouldFetchAllSidecars: parseConf.ShouldFetchAllSidecars,
		},
		logger: logger,
	}, nil
}
