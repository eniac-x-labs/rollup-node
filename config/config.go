package config

import (
	"github.com/eniac-x-labs/anytrustDA/das"
	"github.com/eniac-x-labs/anytrustDA/util/signature"
	"github.com/eniac-x-labs/rollup-node/x/eigenda"
	"github.com/eniac-x-labs/rollup-node/x/nearda"
	"github.com/urfave/cli/v2"

	"github.com/eniac-x-labs/rollup-node/log"
	"github.com/eniac-x-labs/rollup-node/metrics"
	"github.com/eniac-x-labs/rollup-node/txmgr"
)

var defaultBlockTime = uint64(2)

type CLIConfig struct {
	BlockTime     uint64 `json:"block_time"`
	TxMgrConfig   txmgr.CLIConfig
	MetricsConfig metrics.CLIConfig
	LogConfig     log.CLIConfig
}

// NewConfig parses the Config from the provided flags or environment variables.
func NewConfig(ctx *cli.Context) (*CLIConfig, error) {

	return &CLIConfig{
		BlockTime:     defaultBlockTime,
		TxMgrConfig:   txmgr.ReadCLIConfig(ctx),
		MetricsConfig: metrics.ReadCLIConfig(ctx),
		LogConfig:     log.ReadCLIConfig(ctx),
	}, nil
}

type RollupConfig struct {
	AnytrustDAConfig *AnytrustConfig
	//CelestiaDAConfig
	EigenDAConfig *eigenda.EigenDAConfig
	//Eip4844Config
	NearDAConfig *nearda.NearADConfig
}

type AnytrustConfig struct {
	DAConfig          *das.DataAvailabilityConfig
	DataSigner        signature.DataSignerFunc
	DataRetentionTime uint64 // second
}

func NewRollupConfig() *RollupConfig {
	return &RollupConfig{
		AnytrustDAConfig: PrepareAnytrustConfig(),
	}
}

func PrepareAnytrustConfig() *AnytrustConfig {
	//log.Info("preparing config", "config file dir", dir, "config file name", fileName)
	//
	//// 1. set viper args and read config file into 'fileConf'
	//if len(dir) != 0 && len(fileName) != 0 {
	//	viper.SetConfigName(fileName)
	//	viper.AddConfigPath(dir)
	//	log.Debug("node config", "dir", dir, "file name", fileName)
	//} else {
	//	viper.SetConfigName(ConfigName)
	//	viper.AddConfigPath(ConfigDir)
	//	log.Debug("node config", "dir", ConfigDir, "file name", ConfigName)
	//
	//}
	//
	//viper.SetConfigType(ConfigType)
	//
	//// privKey env flag: DRNG_PRIV_KEY
	//viper.SetEnvPrefix(EnvVarPrefix)
	//viper.BindEnv("PRIV_KEY")
	//
	//// read config file
	//if err := viper.ReadInConfig(); err != nil {
	//	fmt.Println("Error reading config file:", err)
	//	return nil, err
	//}
	//
	//confFromFile := &fileConf{}
	//if err := viper.Unmarshal(confFromFile); err != nil {
	//	fmt.Println("Error parsing config file:", err)
	//	return nil, err
	//}
	return &AnytrustConfig{}
}
