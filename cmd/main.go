package main

import (
	"context"
	"fmt"
	"github.com/eniac-x-labs/rollup-node/common/cliapp"
	"github.com/eniac-x-labs/rollup-node/flags"
	"github.com/eniac-x-labs/rollup-node/x/celestia"
	"github.com/eniac-x-labs/rollup-node/x/eip4844"
	"github.com/ethereum/go-ethereum/params"
	"github.com/urfave/cli/v2"
	"os"
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
			Name:        "eip4844",
			Flags:       flags.Flags,
			Description: "Runs the eip-4844 service",
			Action:      cliapp.LifecycleCmd(eip4844.NewEip4844Rollup),
		},
		{
			Name:        "celestia",
			Flags:       flags.Flags,
			Description: "Runs the celestia service",
			Action:      cliapp.LifecycleCmd(celestia.NewCelestiaRollup),
		},
	}

	ctx := context.Background()
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		fmt.Println("Application failed", "message", err)
	}
}
