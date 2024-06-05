package eigenda

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
)

var testConf = EigenDAConfig{
	RPC:                      "disperser-holesky.eigenda.xyz:443",
	StatusQueryRetryInterval: 1 * time.Minute,
	StatusQueryTimeout:       5 * time.Second,
}

func Test_EigendaDisperseBlob(t *testing.T) {
	ast := assert.New(t)
	ctx := context.Background()
	log.SetDefault(log.NewLogger(log.NewTerminalHandler(os.Stdout, true)))
	daCli, err := NewEigenDAClient(&testConf)
	ast.NoError(err)
	txData := []byte("hahaha eigenda")
	reqID, err := daCli.DisperseBlob(ctx, txData)
	ast.NoError(err)
	status, blobInfo, err := daCli.GetBlobStatus(ctx, reqID)
	ast.NoError(err)
	t.Logf("status is %d", status)

	blob_index := blobInfo.GetBlobVerificationProof().GetBlobIndex()
	batch_header_hash := blobInfo.GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()
	t.Logf("blob_index = %d, batch_header_hash = %s", blob_index, batch_header_hash)
}

func Test_GetBlobInfo(t *testing.T) {
	ast := assert.New(t)
	ctx := context.Background()
	log.SetDefault(log.NewLogger(log.NewTerminalHandler(os.Stdout, true)))
	daCli, err := NewEigenDAClient(&testConf)
	ast.NoError(err)

	reqIDBase64 := "MzdkZTcxNDIxODgyZTlhMjg4YmUxN2YxNjUyMjFlZTk0OTI5MDNmZWM2M2YxMmY2MzU4YTg2NGQzZGQxZjQxMi0zMTM3MzEzMzMxMzAzMTM4MzkzOTM5MzUzNjM3MzYzMzMwMzYzNTJmMzAyZjMzMzMyZjMxMmYzMzMzMmZlM2IwYzQ0Mjk4ZmMxYzE0OWFmYmY0Yzg5OTZmYjkyNDI3YWU0MWU0NjQ5YjkzNGNhNDk1OTkxYjc4NTJiODU1"
	reqIDByte, err := base64.StdEncoding.DecodeString(reqIDBase64)
	ast.NoError(err)

	status, info, err := daCli.GetBlobStatus(ctx, reqIDByte)
	ast.NoError(err)
	t.Logf("status = %s, info = %v", status.String(), info)
	t.Logf("batchHeaderHash = %s, blobIndex = %d", info.BlobVerificationProof.GetBatchMetadata().GetBatchHeaderHash(), info.GetBlobVerificationProof().GetBlobIndex())
}

func Test_RetrieveBlob(t *testing.T) {
	ast := assert.New(t)
	ctx := context.Background()
	batchHeadHash := "a215dbf1920cc234121ae4ab4ef34b506e30d6755d4719d76a38d4c687a31a19"
	blobIndex := 494

	log.SetDefault(log.NewLogger(log.NewTerminalHandler(os.Stdout, true)))
	dacli, err := NewEigenDAClient(&testConf)
	ast.NoError(err)
	batchHeadHashByte, err := hex.DecodeString(batchHeadHash)
	ast.NoError(err)

	res, err := dacli.RetrieveBlob(ctx, batchHeadHashByte, uint32(blobIndex))
	ast.NoError(err)
	t.Logf("retrieve blob result: %s", res)
}
