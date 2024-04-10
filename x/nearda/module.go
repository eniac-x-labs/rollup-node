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
	Network  string // 目前nearDA只支持 "Mainnet", "Testnet", "Localnet"这3个string
	Ns       uint32
}

func NewNearDAClient(nearconf NearADConfig) (*NearDAClient, error) {
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
//func (n *NearDAClient) GetFromDa(frameRefBytes []byte, txIndex uint32) ([]byte, error) {
//	return n.Get(frameRefBytes, txIndex)
//}
//
//func (n *NearDAClient) ForceSubmitData(data []byte) ([]byte, error) {
//	return n.ForceSubmit(data)
//}
