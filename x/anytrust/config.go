package anytrust

import "time"

type AnytrustConfig struct {
	RpcUrl            string `toml:"rpcUrl"`
	RestfulUrl        string `toml:"restfulUrl"`
	DataRetentionTime uint64 `toml:"dataRetentionTime"`

	RandomMessageSize  int           `toml:"randomMessageSize"`
	DASRetentionPeriod time.Duration `toml:"dasRetentionPeriod"`
	SigningKey         string        `toml:"signingKey"`
	//SigningWallet         string        `toml:"signingWallet"`
	//SigningWalletPassword string        `toml:"signingWalletPassword"`
}

const (
	RpcUrlFlag             = "rpcUrl"
	RestfulUrlFlag         = "restfulUrl"
	DataRetentionTimeFlag  = "dataRetentionTime"
	RandomMessageSizeFlag  = "randomMessageSize"
	DasRetentionPeriodFlag = "dasRetentionPeriod"
	SigningKeyFlag         = "signingKey"
)

var AnytrustDAEnvFlags = []string{
	RpcUrlFlag,
	RestfulUrlFlag,
	DataRetentionTimeFlag,
	RandomMessageSizeFlag,
	DasRetentionPeriodFlag,
	SigningKeyFlag,
}
