package config

import (
	"errors"

	"github.com/ethereum/go-ethereum/log"

	"github.com/eniac-x-labs/anytrustDA/das"
	"github.com/eniac-x-labs/anytrustDA/util/signature"
	cli_config "github.com/eniac-x-labs/rollup-node/config/cli-config"
	"github.com/eniac-x-labs/rollup-node/x/anytrust"
	"github.com/eniac-x-labs/rollup-node/x/celestia"
	"github.com/eniac-x-labs/rollup-node/x/eip4844"

	//"github.com/eniac-x-labs/anytrustDA/das"
	//"github.com/eniac-x-labs/anytrustDA/util/signature"
	"github.com/eniac-x-labs/rollup-node/x/eigenda"
	"github.com/eniac-x-labs/rollup-node/x/nearda"
	"github.com/spf13/viper"
)

type RollupConfig struct {
	AnytrustDAConfig *anytrust.AnytrustConfig //*AnytrustConfig
	CelestiaDAConfig *celestia.CelestiaConfig
	CelestiaCLICfg   *cli_config.CLIConfig
	EigenDAConfig    *eigenda.EigenDAConfig
	Eip4844Config    *eip4844.Eip4844Config
	Eip4844CLICfg    *cli_config.CLIConfig
	NearDAConfig     *nearda.NearDAConfig
}

type AnytrustConfig struct {
	DAConfig          *das.DataAvailabilityConfig
	DataSigner        signature.DataSignerFunc
	DataRetentionTime uint64 // second
}

var (
	AnytrustConfigDir  = defaultConfigDir
	AnytrustConfigFile = "anytrust"
	CelestiaConfigDir  = defaultConfigDir
	CelestiaConfigFile = "celestia"
	EigenDAConfigDir   = defaultConfigDir
	EigenDAConfigFile  = "eigenda"
	Eip4844ConfigDir   = defaultConfigDir
	Eip4844ConfigFile  = "eip4844"
	NearDAConfigDir    = defaultConfigDir
	NearDAConfigFile   = "nearda"
	ApiConfigDir       = defaultConfigDir
	ApiConfigFile      = "api"
)

const (
	defaultConfigDir = "./config"
	ConfigType       = "toml"
	EnvVarPrefix     = "ROLLUP"
	AnytrustPrefix   = "anytrust"
	CelestiaPrefix   = "celestia"
	EigenDAPrefix    = "eigenda"
	Eip4844Prefix    = "eip4844"
	NearDAPrefix     = "nearda"
	ApiPrefix        = "api"
)

func NewRollupConfig() *RollupConfig {
	// Anytrust
	anytrustDAConf := &anytrust.AnytrustConfig{}
	if err := PrepareConfig(AnytrustConfigDir, AnytrustConfigFile, anytrustDAConf, AnytrustPrefix, anytrust.AnytrustDAEnvFlags); err != nil {
		log.Error("PrepareConfig failed", "da-type", "AnytrustDA")
	}

	// Celestia
	var celestiaConfig *celestia.CelestiaConfig
	celestiaParseConf := &celestia.ParseCelestiaConfig{}
	if err := PrepareConfig(CelestiaConfigDir, CelestiaConfigFile, celestiaParseConf, CelestiaPrefix, []string{}); err != nil {
		log.Error("PrepareConfig failed", "da-type", "Celestia")
	} else {
		celestiaConfig, err = celestia.ProcessCelestiaConfig(celestiaParseConf, log.Root())
		if err != nil {
			log.Error("Process celestia config failed", "err", err)
		}
	}
	celestiaCliCfg := cli_config.ParseCLIConfig(celestiaParseConf.BlockTime, celestiaParseConf.Enable, celestiaParseConf.ListenAddr, celestiaParseConf.ListenPort)

	// Eigen
	eigendaConf := &eigenda.EigenDAConfig{}
	if err := PrepareConfig(EigenDAConfigDir, EigenDAConfigFile, eigendaConf, EigenDAPrefix, eigenda.EigenDAEnvFlags); err != nil {
		log.Error("PrepareConfig failed", "da-type", "EigenDA")
	}

	// EIP-4844
	var eip4844Config *eip4844.Eip4844Config
	eip4844ParseConf := &eip4844.ParseEip4844Config{}
	if err := PrepareConfig(EigenDAConfigDir, EigenDAConfigFile, eip4844ParseConf, Eip4844Prefix, []string{}); err != nil {
		log.Error("PrepareConfig failed", "da-type", "eip-4844")
	} else {
		eip4844Config, err = eip4844.ProcessEip4844Config(eip4844ParseConf, log.Root())
		if err != nil {
			log.Error("Process eip4844 config failed", "err", err)
		}
	}
	eip4844CliCfg := cli_config.ParseCLIConfig(eip4844ParseConf.BlockTime, eip4844ParseConf.Enable, eip4844ParseConf.ListenAddr, eip4844ParseConf.ListenPort)

	// Near
	neardaConf := &nearda.NearDAConfig{}
	if err := PrepareConfig(NearDAConfigDir, NearDAConfigFile, neardaConf, NearDAPrefix, nearda.NearDAEnvFlags); err != nil {
		log.Error("PrepareConfig failed", "da-type", "NearDA")
	}
	return &RollupConfig{
		AnytrustDAConfig: anytrustDAConf,
		CelestiaDAConfig: celestiaConfig,
		CelestiaCLICfg:   celestiaCliCfg,
		EigenDAConfig:    eigendaConf,
		Eip4844Config:    eip4844Config,
		Eip4844CLICfg:    eip4844CliCfg,
		NearDAConfig:     neardaConf,
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
