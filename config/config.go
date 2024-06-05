package config

import (
	"errors"

	//"github.com/eniac-x-labs/anytrustDA/das"
	//"github.com/eniac-x-labs/anytrustDA/util/signature"
	"github.com/eniac-x-labs/rollup-node/x/eigenda"
	"github.com/eniac-x-labs/rollup-node/x/nearda"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"

	"github.com/eniac-x-labs/rollup-node/log"
	"github.com/eniac-x-labs/rollup-node/metrics"
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

type RollupConfig struct {
	AnytrustDAConfig *AnytrustConfig
	//CelestiaDAConfig
	EigenDAConfig *eigenda.EigenDAConfig
	//Eip4844Config
	NearDAConfig *nearda.NearDAConfig
}

type AnytrustConfig struct {
	//DAConfig          *das.DataAvailabilityConfig
	//DataSigner        signature.DataSignerFunc
	//DataRetentionTime uint64 // second
}

var (
	AnytrustConfigDir  = defaultConfigDir
	AnytrustConfigFile = "anyrtust"
	EigenDAConfigDir   = defaultConfigDir
	EigenDAConfigFile  = "eigenda"
	NearDAConfigDir    = defaultConfigDir
	NearDAConfigFile   = "nearda"
)

const (
	defaultConfigDir = "./config"
	ConfigType       = "toml"
	EnvVarPrefix     = "ROLLUP"
	AnytrustPrefix   = "anytrust"
	EigenDAPrefix    = "eigenda"
	NearDAPrefix     = "nearda"
)

func NewRollupConfig() *RollupConfig {
	//anytrustDAConf := &das.DataAvailabilityConfig{}
	anytrustDAConf := &struct{}{}
	if err := PrepareConfig(AnytrustConfigDir, AnytrustConfigFile, anytrustDAConf, AnytrustPrefix, []string{}); err != nil {
		log.Error("PrepareConfig failed", "da-type", "AnytrustDA")
	}

	eigendaConf := &eigenda.EigenDAConfig{}
	if err := PrepareConfig(EigenDAConfigDir, EigenDAConfigFile, eigendaConf, EigenDAPrefix, []string{}); err != nil {
		log.Error("PrepareConfig failed", "da-type", "EigenDA")
	}

	neardaConf := &nearda.NearDAConfig{}
	if err := PrepareConfig(NearDAConfigDir, NearDAConfigFile, neardaConf, NearDAPrefix, []string{}); err != nil {
		log.Error("PrepareConfig failed", "da-type", "NearDA")
	}
	return &RollupConfig{
		AnytrustDAConfig: &AnytrustConfig{
			//DAConfig:          anytrustDAConf,
			//DataRetentionTime: uint64(anytrustDAConf.RestAggregator.SyncToStorage.RetentionPeriod.Seconds()),
		},
		EigenDAConfig: eigendaConf,
		NearDAConfig:  neardaConf,
	}
}

func PrepareConfig(dir, fileName string, target interface{}, envPrefix string, envFlags []string) error {
	log.Debug("Preparing config", "dir", dir, "file_name", fileName, "config_type", ConfigType)

	// 1. set viper args and read config file into 'fileConf'
	if len(dir) == 0 || len(fileName) == 0 {
		log.Error("config dir or file name with empty string")
		return errors.New("config dir or file name with empty string")
	}

	viper.SetConfigName(fileName)
	viper.AddConfigPath(dir)
	viper.SetConfigType(ConfigType)

	// env flag: envFlags_f
	if envFlags != nil && len(envFlags) != 0 {
		viper.SetEnvPrefix(envPrefix)
		for _, f := range envFlags {
			viper.BindEnv(f)
		}
	}

	// read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Error("reading config file failed", "err", err)
		return err
	}

	if err := viper.Unmarshal(target); err != nil {
		log.Error("unmarshal config file failed", "err", err)
		return err
	}

	return nil
}
