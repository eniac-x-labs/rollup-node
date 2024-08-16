package anytrust

type AnytrustConfig struct {
	RpcUrl            string `toml:"rpc_url" mapstructure:"rpc_url"`
	RestfulUrl        string `toml:"restful_url" mapstructure:"restful_url"`
	DataRetentionTime uint64 `toml:"data_retention_time" mapstructure:"data_retention_time"`

	RandomMessageSize int `toml:"random_message_size" mapstructure:"random_message_size"`
	//DASRetentionPeriod time.Duration `toml:"dasRetentionPeriod"`
	SigningKey string `toml:"signing_key" mapstructure:"signing_key"`
	//SigningWallet         string        `toml:"signingWallet"`
	//SigningWalletPassword string        `toml:"signingWalletPassword"`
}

const (
	RpcUrlFlag            = "rpc_url"
	RestfulUrlFlag        = "restful_url"
	DataRetentionTimeFlag = "data_retention_time"
	RandomMessageSizeFlag = "random_message_size"
	//DasRetentionPeriodFlag = "dasRetentionPeriod"
	SigningKeyFlag = "signing_key"
)

// AnytrustDAEnvFlags The env flag is like prefix_flag, with all letters in uppercase.
var AnytrustDAEnvFlags = []string{
	RpcUrlFlag,
	RestfulUrlFlag,
	DataRetentionTimeFlag,
	RandomMessageSizeFlag,
	//DasRetentionPeriodFlag,
	SigningKeyFlag,
}
