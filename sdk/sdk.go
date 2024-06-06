package sdk

import (
	"net/rpc"

	_rpc "github.com/eniac-x-labs/rollup-node/rpc"
	"github.com/ethereum/go-ethereum/log"
)

type RollupSDK struct {
	*rpc.Client
}

func NewRollupSdk(addr string) (_rpc.RollupInter, error) {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Error("rpc Dial failed", "err", err)
		return nil, err
	}
	return &RollupSDK{client}, nil
}

func (s *RollupSDK) RollupWithType(data []byte, daType int) ([]interface{}, error) {
	var res []interface{}
	err := s.Call("RollupRpcServer.Rollup", _rpc.RollupRequest{
		DAType: daType,
		Data:   data,
	}, &res)
	return res, err
}

func (s *RollupSDK) RetrieveFromDAWithType(daType int, args interface{}) ([]byte, error) {
	var res []byte
	err := s.Call("RollupRpcServer.Retrieve", _rpc.RetrieveRequest{
		DAType: daType,
		Args:   args,
	}, &res)
	return res, err
}
