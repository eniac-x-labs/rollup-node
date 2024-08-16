package service

type RollupInter interface {
	RollupWithType(data []byte, daType int) ([]interface{}, error)
	RetrieveFromDAWithType(daType int, args interface{}) ([]byte, error)
}

type HandlerSvc struct {
	RollupInter
}

func New(rollup RollupInter) HandlerSvc {
	return HandlerSvc{rollup}
}
