package nearda

import (
	"github.com/ethereum/go-ethereum/log"
	near "github.com/near/rollup-data-availability/gopkg/da-rpc"
)

type NearDA struct {
	*near.Config
}

type NearADConfig struct {
	Account  string
	Contract string
	Key      string
	Ns       uint32
}

func NewNearDA(nearconf NearADConfig) (*NearDA, error) {
	conf, err := near.NewConfig(nearconf.Account, nearconf.Contract, nearconf.Key, nearconf.Ns)
	if err != nil {
		log.Error("")
		return nil, err
	}
	return &NearDA{conf}, nil
}

func (n *NearDA) Submit() {
	n.Submit()
}
