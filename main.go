package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	_config "github.com/eniac-x-labs/rollup-node/config"
	_core "github.com/eniac-x-labs/rollup-node/core"
	_rpc "github.com/eniac-x-labs/rollup-node/rpc"
	"github.com/ethereum/go-ethereum/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelDebug, true)))

	var (
		rpcAddress string
		apiAddress string
	)
	flag.StringVar(&rpcAddress, "rpcAddress", "", "listen address for rpc and sdk")
	flag.StringVar(&apiAddress, "apiAddress", "", "listen address for web server")
	flag.Parse()

	rollupModule, err := _core.NewRollupModule(ctx, _config.NewRollupConfig())
	if err != nil {
		log.Error("NewRollupModule failed", "err", err)
		return
	}

	// start rpc for sdk
	if len(rpcAddress) != 0 {
		go _rpc.NewAndStartRollupRpcServer(rpcAddress, rollupModule)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutting down server...")
	cancel()

	fmt.Println("Server gracefully stopped")
}
