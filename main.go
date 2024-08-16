package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/eniac-x-labs/rollup-node/api"

	_config "github.com/eniac-x-labs/rollup-node/config"
	_core "github.com/eniac-x-labs/rollup-node/core"
	_rpc "github.com/eniac-x-labs/rollup-node/rpc"
	"github.com/ethereum/go-ethereum/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	logger := log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelDebug, true))
	log.SetDefault(logger)

	var (
		rpcAddress string
		apiAddress string
	)
	flag.StringVar(&rpcAddress, "rpcAddress", "", "listen address for rpc and sdk")
	flag.StringVar(&apiAddress, "apiAddress", "", "listen address for web server")
	flag.Parse()

	if len(rpcAddress) == 0 && len(apiAddress) == 0 {
		flag.Usage()
	}

	rollupModule, err := _core.NewRollupModuleWithConfig(ctx, _config.NewRollupConfig())
	if err != nil {
		log.Error("NewRollupModule failed", "err", err)
		return
	}

	// start rpc for sdk
	//var wg sync.WaitGroup
	if len(rpcAddress) != 0 {
		//wg.Add(1)
		go _rpc.NewAndStartRollupRpcServer(ctx, rpcAddress, rollupModule)
	}

	err = api.NewApi(ctx, logger, apiAddress, rollupModule)
	if err != nil {
		log.Error("NewApi failed", "err", err)
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutting down server...")
	cancel()

	//wg.Wait()
	fmt.Println("Server gracefully stopped")
}
