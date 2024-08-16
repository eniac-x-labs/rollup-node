package main

import (
	"context"
	"fmt"
	"os"

	"github.com/eniac-x-labs/rollup-node/common/cliapp"
	"github.com/eniac-x-labs/rollup-node/core"
	"github.com/eniac-x-labs/rollup-node/flags"
	"github.com/ethereum/go-ethereum/params"
	"github.com/urfave/cli/v2"
)

var (
	GitCommit = ""
	GitDate   = ""
)

func main() {
	app := cli.NewApp()
	app.Version = params.VersionWithCommit(GitCommit, GitDate)
	app.Usage = "Rollup Node Service"
	app.Description = "Service for generating and submitting L2 tx batches to L1"

	app.Commands = []*cli.Command{
		{
			Name:        "rollup-node",
			Flags:       flags.Flags,
			Description: "Runs the rollup node service",
			Action:      cliapp.LifecycleCmd(core.RunRollupModuleForCLI),
		},
	}

	ctx := context.Background()
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		fmt.Println("Application failed", "message", err)
	}
}
