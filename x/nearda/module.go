package nearda

import (
	"github.com/ethereum/go-ethereum/log"
	near "github.com/near/rollup-data-availability/gopkg/da-rpc"
)

type NearDAClient struct {
	*near.Config
}

type NearDAConfig struct {
	Account  string `toml:"account"`
	Contract string `toml:"contract"`
	Key      string `toml:"key"`
	Network  string `toml:"network"` // nearDA only support "Mainnet", "Testnet", "Localnet"
	Ns       uint32 `toml:"ns"`
}

type INearDA interface {
	Store(data []byte) ([]byte, error)
	GetFromDA(frameRefBytes []byte, txIndex uint32) ([]byte, error)
}

func NewNearDAClient(nearconf *NearDAConfig) (INearDA, error) {
	conf, err := near.NewConfig(nearconf.Account, nearconf.Contract, nearconf.Key, nearconf.Network, nearconf.Ns)
	if err != nil {
		log.Error("NewConfig failed:", err)
		return nil, err
	}
	return &NearDAClient{conf}, nil
}

//func (n *NearDAClient) SubmitData(candidateHex string, data []byte) ([]byte, error) {
//	return n.Submit(candidateHex, data)
//}
//

func (n *NearDAClient) Store(data []byte) ([]byte, error) {
	return n.ForceSubmit(data)
}
func (n *NearDAClient) GetFromDA(frameRefBytes []byte, txIndex uint32) ([]byte, error) {
	return n.Get(frameRefBytes, txIndex)
}
