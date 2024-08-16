package celestia

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/celestiaorg/celestia-openrpc/types/share"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/log"
)

func TestSubmitBlob(t *testing.T) {
	ctx := context.Background()
	logger := log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelDebug, true))
	log.SetDefault(logger)

	namespace, _ := share.NewBlobNamespaceV0([]byte("DappLink"))

	celestiaCfg := &CelestiaConfig{
		celestiaConfig: CLIConfig{
			DaRpc:     "http://localhost:26658",
			AuthToken: "",
			Namespace: namespace.String(),
		},
		logger: logger,
	}
	celestiaRollup, err := NewCelestiaRollupWithConfig(ctx, celestiaCfg)
	require.NoError(t, err)

	data, err := celestiaRollup.SubmitBlob(ctx, []byte("hello DappLink"))
	require.NoError(t, err)

	t.Log("send celestia da transaction success", fmt.Sprintf("height: %v", data))
}

func TestRetrievedBlobs(t *testing.T) {
	ctx := context.Background()
	logger := log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelDebug, true))
	log.SetDefault(logger)

	namespace, _ := share.NewBlobNamespaceV0([]byte("DappLink"))

	celestiaCfg := &CelestiaConfig{
		celestiaConfig: CLIConfig{
			DaRpc:     "http://localhost:26658",
			AuthToken: "",
			Namespace: namespace.String(),
		},
		logger: logger,
	}
	celestiaRollup, err := NewCelestiaRollupWithConfig(ctx, celestiaCfg)
	require.NoError(t, err)

	data, err := celestiaRollup.RetrievedBlobs(ctx, 2075034)
	require.NoError(t, err)

	t.Log("Retrieve celestia transaction success! ", fmt.Sprintf("data: %s", string(data)))
}
