package rpc

type RollupInter interface {
	RollupWithType(data []byte, daType int) ([]interface{}, error)
	RetrieveFromDAWithType(daType int, args interface{}) ([]byte, error)
}

type DRNGRpcInterface interface {
	Rollup(req RollupRequest, reply *[]interface{}) error
	Retrieve(req RetrieveRequest, reply *[]byte) error
}

//type DAInter interface {
//	Store
//}
