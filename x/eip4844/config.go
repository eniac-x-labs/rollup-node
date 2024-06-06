package eip4844

type Eip4844Config struct {
	L1Rpc                  string `toml:"l1Rpc"`
	PrivateKey             string `toml:"privateKey"`
	L1ChainID              string `toml:"l1ChainID"` // *bigInt
	UseBlobs               bool   `toml:"useBlobs"`
	L1BeaconAddr           string `toml:"l1BeaconAddr"`
	ShouldFetchAllSidecars bool   `toml:"shouldFetchAllSidecars"`
	BatchInboxAddress      string `toml:"batchInboxAddress"` // common.Address
	BatcherAddr            string `toml:"batcherAddr"`       // common.Address
}
