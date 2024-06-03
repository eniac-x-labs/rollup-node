package nearda

import (
	"github.com/ethereum/go-ethereum/log"
	near "github.com/near/rollup-data-availability/gopkg/da-rpc"
)

type NearDAClient struct {
	*near.Config
}

type NearADConfig struct {
	Account  string
	Contract string
	Key      string
	Network  string // nearDA only support "Mainnet", "Testnet", "Localnet"
	Ns       uint32
}

type INearDA interface {
	Store(data []byte) ([]byte, error)
	GetFromDA(frameRefBytes []byte, txIndex uint32) ([]byte, error)
}

func NewNearDAClient(nearconf *NearADConfig) (INearDA, error) {
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
