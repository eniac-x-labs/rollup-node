package eip4844

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	cli_config "github.com/eniac-x-labs/rollup-node/config/cli-config"
)

func TestSendBlobTransaction(t *testing.T) {
	ctx := context.Background()
	logger := log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelDebug, true))
	log.SetDefault(logger)

	cliCfg := &cli_config.CLIConfig{
		L1Rpc:      "https://1rpc.io/sepolia",
		L1ChainID:  big.NewInt(11155111),
		PrivateKey: "",
	}

	signer := types.NewCancunSigner(cliCfg.L1ChainID)
	eip4844Cfg := &Eip4844Config{
		eip4844Config: CLIConfig{
			DSConfig: &DataSourceConfig{
				l1Signer:          signer,
				batchInboxAddress: common.HexToAddress("0x4F34C922fB0D80c7d79Ac25e497d90d7efa513C2"),
				batcherAddr:       common.HexToAddress("0x2822E13eF080475e8CaBe39b3dc65c6dbe9b083a"),
			},
			UseBlobs:               true,
			L1BeaconAddr:           "",
			ShouldFetchAllSidecars: false,
		},
		logger: logger,
	}
	eip4844Rollup, err := NewEip4844WithConfig(ctx, cliCfg, eip4844Cfg)
	require.NoError(t, err)

	data, err := eip4844Rollup.SendTransaction(ctx, []byte("hello dappLink"))
	require.NoError(t, err)

	t.Log("send eip4844 transaction success", fmt.Sprintf("tx_hash: 0x%s", hex.EncodeToString(data)))
}

func TestRetrieveBlobTransaction(t *testing.T) {
	ctx := context.Background()
	logger := log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelDebug, true))
	log.SetDefault(logger)

	cliCfg := &cli_config.CLIConfig{
		L1Rpc:      "https://ethereum-sepolia-rpc.publicnode.com",
		L1ChainID:  big.NewInt(11155111),
		PrivateKey: "",
	}

	signer := types.NewCancunSigner(cliCfg.L1ChainID)
	eip4844Cfg := &Eip4844Config{
		eip4844Config: CLIConfig{
			DSConfig: &DataSourceConfig{
				l1Signer:          signer,
				batchInboxAddress: common.HexToAddress("0x4F34C922fB0D80c7d79Ac25e497d90d7efa513C2"),
				batcherAddr:       common.HexToAddress("0x2822E13eF080475e8CaBe39b3dc65c6dbe9b083a"),
			},
			UseBlobs:               true,
			L1BeaconAddr:           "",
			ShouldFetchAllSidecars: false,
		},
		logger: logger,
	}
	eip4844Rollup, err := NewEip4844WithConfig(ctx, cliCfg, eip4844Cfg)
	require.NoError(t, err)

	data, err := eip4844Rollup.DataFromEVMTransactions(ctx, "0x0412ee533cb3243fa938b7ee61bed48d77cb05b8ab0181e8cde0ec0c8f54f774")
	require.NoError(t, err)

	t.Log("Retrieve eip4844 transaction success! ", fmt.Sprintf("data: %s", string(data)))
}
