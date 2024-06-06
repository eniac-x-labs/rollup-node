package rpc

import (
	"context"
	"net"
	"net/rpc"

	"github.com/ethereum/go-ethereum/log"
)

type RollupRequest struct {
	DAType int
	Data   []byte
}

type RetrieveRequest struct {
	DAType int
	Args   interface{}
}

type RollupRpcServer struct {
	RollupInter
}

func NewAndStartRollupRpcServer(ctx context.Context, address string, rollup RollupInter) {
	if err := rpc.Register(&RollupRpcServer{
		rollup,
	}); err != nil {
		log.Error("RpcServer Register failed", "err", err)
		return
	}
	log.Debug("RpcServer Register finished")

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Error("RpcServer Listen failed", "err", err, "address", address)
		return
	}
	log.Debug("RpcServer listen address finished", "address", address)

	for {
		select {
		case <-ctx.Done():
			listener.Close()
			log.Info("rollup rpc listener closed successfully")
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Error("RpcServer listener.Accept failed", "err", err)
			}

			go rpc.ServeConn(conn)
		}
	}
}

func (s *RollupRpcServer) Rollup(req RollupRequest, reply *[]interface{}) error {
	var err error
	*reply, err = s.RollupWithType(req.Data, req.DAType)
	if err != nil {
		return err
	}

	return nil
}

func (s *RollupRpcServer) Retrieve(req RetrieveRequest, reply *[]byte) error {
	var err error
	*reply, err = s.RetrieveFromDAWithType(req.DAType, req.Args)
	if err != nil {
		return err
	}
	return nil
}
