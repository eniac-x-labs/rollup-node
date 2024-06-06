package eip4844

import (
	_log "github.com/eniac-x-labs/rollup-node/log"
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
}

type Eip4844Config struct {
	//cfg           *config.CLIConfig
	eip4844Config CLIConfig
	logger        _log.Logger
}

func ProcessEip4844Config(parseConf *ParseEip4844Config, logger _log.Logger) (*Eip4844Config, error) {
	// todo: kezi
	return nil, nil
}
