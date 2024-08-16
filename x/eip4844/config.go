package eip4844

import (
	"github.com/ethereum/go-ethereum/log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ParseEip4844Config struct {
	UseBlobs               bool   `toml:"useBlobs"`
	L1BeaconAddr           string `toml:"l1BeaconAddr"`
	ShouldFetchAllSidecars bool   `toml:"shouldFetchAllSidecars"`

	// data source config
	BatchInboxAddress string `toml:"batchInboxAddress"` // common.Address
	BatcherAddr       string `toml:"batcherAddr"`       // common.Address

	// BlockTime CLI Config
	BlockTime uint64 `toml:"blockTime"`

	// Log CLI Config
	Level      int    `toml:"level"`
	Color      bool   `toml:"color"`
	FormatType string `toml:"formatType"`

	L1ChainIdFlagName uint64 `toml:"l1ChainIdFlagName"`
}

type Eip4844Config struct {
	//cfg           *config.CLIConfig
	eip4844Config CLIConfig
	logger        log.Logger
}

func ProcessEip4844Config(parseConf *ParseEip4844Config, logger log.Logger) (*Eip4844Config, error) {
	signer := types.NewCancunSigner(new(big.Int).SetUint64(parseConf.L1ChainIdFlagName))

	return &Eip4844Config{
		eip4844Config: CLIConfig{
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
