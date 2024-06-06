package service

import (
	_core "github.com/eniac-x-labs/rollup-node/core"
)

type HandlerSvc struct {
	_core.RollupInter
}

func New(rollup _core.RollupInter) HandlerSvc {
	return HandlerSvc{rollup}
}
