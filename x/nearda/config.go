package nearda

type NearDAConfig struct {
	Account  string `toml:"account"`
	Contract string `toml:"contract"`
	Key      string `toml:"key"`
	Network  string `toml:"network"` // nearDA only support "Mainnet", "Testnet", "Localnet"
	Ns       uint32 `toml:"ns"`
}

const (
	AccountFlag  = "account"
	ContractFlag = "contract"
	KeyFlag      = "key"
	NetworkFlag  = "network"
	NsFlag       = "ns"
)

// NearDAEnvFlags The env flag is like prefix_flag, with all letters in uppercase.
var NearDAEnvFlags = []string{
	AccountFlag,
	ContractFlag,
	KeyFlag,
	NetworkFlag,
	NsFlag,
}
