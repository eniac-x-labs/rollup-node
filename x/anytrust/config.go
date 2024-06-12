package anytrust

type AnytrustConfig struct {
	RpcUrl            string `toml:"rpcUrl"`
	RestfulUrl        string `toml:"restfulUrl"`
	DataRetentionTime uint64 `toml:"dataRetentionTime"`

	RandomMessageSize int `toml:"randomMessageSize"`
	//DASRetentionPeriod time.Duration `toml:"dasRetentionPeriod"`
	SigningKey string `toml:"signingKey"`
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

// AnytrustDAEnvFlags The env flag is like prefix_flag, with all letters in uppercase.
var AnytrustDAEnvFlags = []string{
	RpcUrlFlag,
	RestfulUrlFlag,
	DataRetentionTimeFlag,
	RandomMessageSizeFlag,
	DasRetentionPeriodFlag,
	SigningKeyFlag,
}
