package eip4844

import (
	"math/big"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	eth "github.com/eniac-x-labs/rollup-node/eth-serivce"
)

const (
	L1ChainIdFlagName                = "l1.chain-id"
	DataAvailabilityTypeFlagName     = "data-availability-type"
	L1BeaconFlagName                 = "l1.beacon"
	L1BeaconFetchAllSidecarsFlagName = "l1.beacon.fetch-all-sidecars"
	BatcherAddressFlagName           = "eip4844.batcher-address"
	BatchInboxAddressFlagName        = "eip4844.batch-inbox-address"
)

func CLIFlags(envPrefix string) []cli.Flag {
	return []cli.Flag{
		&cli.Uint64Flag{
			Name:     L1ChainIdFlagName,
			Usage:    "The chain id of l1.",
			Required: false,
			EnvVars:  eth.PrefixEnvVar(envPrefix, "L1_CHAIN_ID"),
		},
		&cli.StringFlag{
			Name:    DataAvailabilityTypeFlagName,
			Usage:   "The data availability type to use for submitting batches to the L1.",
			Value:   "blobs",
			EnvVars: eth.PrefixEnvVar(envPrefix, "DATA_AVAILABILITY_TYPE"),
		},
		&cli.StringFlag{
			Name:     L1BeaconFlagName,
			Usage:    "Address of L1 Beacon-node HTTP endpoint to use.",
			Required: false,
			EnvVars:  eth.PrefixEnvVar(envPrefix, "L1_BEACON"),
		},
		&cli.BoolFlag{
			Name:     L1BeaconFetchAllSidecarsFlagName,
			Usage:    "If true, all sidecars are fetched and filtered locally. Workaround for buggy Beacon nodes.",
			Required: false,
			Value:    false,
			EnvVars:  eth.PrefixEnvVar(envPrefix, "L1_BEACON_FETCH_ALL_SIDECARS"),
		},
		&cli.StringFlag{
			Name:     BatcherAddressFlagName,
			Usage:    "Address of eip4844 Batcher.",
			Required: false,
			EnvVars:  eth.PrefixEnvVar(envPrefix, "EIP4844_BATCHER_ADDRESS"),
		},
		&cli.StringFlag{
			Name:     BatchInboxAddressFlagName,
			Usage:    "Address of eip4844 Batch inbox.",
			Required: false,
			EnvVars:  eth.PrefixEnvVar(envPrefix, "EIP4844_BATCH_INBOX_ADDRESS"),
		},
	}
}

type CLIConfig struct {
	L1ChainID              *big.Int
	DSConfig               *DataSourceConfig
	UseBlobs               bool
	L1BeaconAddr           string
	ShouldFetchAllSidecars bool
}

func (c CLIConfig) Check() error {

	return nil
}

func NewCLIConfig() CLIConfig {
	return CLIConfig{}
}

func ReadCLIConfig(ctx *cli.Context) CLIConfig {
	var useBlobs bool
	switch ctx.String(DataAvailabilityTypeFlagName) {
	case "blobs":
		useBlobs = true
	case "calldata":
		useBlobs = false
	}

	signer := types.NewCancunSigner(new(big.Int).SetUint64(ctx.Uint64(L1ChainIdFlagName)))

	dsConfig := DataSourceConfig{
		l1Signer:          signer,
		batchInboxAddress: common.HexToAddress(ctx.String(BatchInboxAddressFlagName)),
		batcherAddr:       common.HexToAddress(ctx.String(BatcherAddressFlagName)),
	}

	return CLIConfig{
		L1ChainID:              new(big.Int).SetUint64(ctx.Uint64(L1ChainIdFlagName)),
		DSConfig:               &dsConfig,
		UseBlobs:               useBlobs,
		L1BeaconAddr:           ctx.String(L1BeaconFlagName),
		ShouldFetchAllSidecars: ctx.Bool(L1BeaconFetchAllSidecarsFlagName),
	}
}
