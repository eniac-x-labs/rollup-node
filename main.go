package main

import (
	"context"

	_config "github.com/eniac-x-labs/rollup-node/config"
	_core "github.com/eniac-x-labs/rollup-node/core"
)

func main() {
	ctx, cancle := context.WithCancel(context.Background())
	conf := _config.NewRollupConfig()
	rollup, err := _core.NewRollupModule(ctx, conf)
}
