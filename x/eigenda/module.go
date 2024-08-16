package eigenda

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/slog"

	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type IEigenDA interface {
	RetrieveBlob(ctx context.Context, BatchHeaderHash []byte, BlobIndex uint32) ([]byte, error)
	DisperseBlob(ctx context.Context, txData []byte) ([]byte, error)
	GetBlobStatus(ctx context.Context, reqID []byte) (disperser.BlobStatus, *disperser.BlobInfo, error)
	DisperseBlobAndGetBlobInfo(ctx context.Context, txData []byte) (*disperser.BlobInfo, error)
}

type EigenDAClient struct {
	DisperserCli disperser.DisperserClient
	EigenDAConfig
	logger log.Logger
}

func NewEigenDAClient(cfg *EigenDAConfig) (IEigenDA, error) {
	config := &tls.Config{}
	credential := credentials.NewTLS(config)
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(credential)}
	conn, err := grpc.Dial(cfg.RPC, dialOptions...)
	if err != nil {
		return nil, err
	}
	daClient := disperser.NewDisperserClient(conn)

	logger := log.Root().With(slog.String("module", "eigenda"))
	return &EigenDAClient{
		DisperserCli: daClient,
		EigenDAConfig: EigenDAConfig{
			RPC:                      cfg.RPC,
			StatusQueryTimeout:       cfg.StatusQueryTimeout,
			StatusQueryRetryInterval: cfg.StatusQueryRetryInterval,
		},
		logger: logger,
	}, nil
}

func (m *EigenDAClient) RetrieveBlob(ctx context.Context, BatchHeaderHash []byte, BlobIndex uint32) ([]byte, error) {
	if m.DisperserCli == nil {
		return nil, errors.New("eigendDA disperserCli is nil")
	}

	reply, err := m.DisperserCli.RetrieveBlob(ctx, &disperser.RetrieveBlobRequest{
		BatchHeaderHash: BatchHeaderHash,
		BlobIndex:       BlobIndex,
	})
	if err != nil {
		return nil, err
	}

	// decode modulo bn254
	decodedData := codec.RemoveEmptyByteFromPaddedBytes(reply.Data)

	return decodedData, nil
}

func (m *EigenDAClient) DisperseBlobAndGetBlobInfo(ctx context.Context, txData []byte) (*disperser.BlobInfo, error) {
	m.logger.Info("Attempting to disperse blob to EigenDA")
	if m.DisperserCli == nil {
		return nil, errors.New("eigendDA disperserCli is nil")
	}

	// encode modulo bn254
	encodedTxData := codec.ConvertByPaddingEmptyByte(txData)

	disperseReq := &disperser.DisperseBlobRequest{
		Data: encodedTxData,
	}
	disperseRes, err := m.DisperserCli.DisperseBlob(ctx, disperseReq)

	if err != nil || disperseRes == nil {
		m.logger.Error("Unable to disperse blob to EigenDA, aborting", "err", err)
		return nil, err
	}
	m.logger.Debug("daClient.DisperseBlob", "disperseRes", disperseRes)
	m.logger.Debug("daClient.DisperseBlob", "disperseRes.Result", disperseRes.Result)
	if disperseRes.Result == disperser.BlobStatus_UNKNOWN ||
		disperseRes.Result == disperser.BlobStatus_FAILED {
		m.logger.Error("Unable to disperse blob to EigenDA, aborting", "err", err)
		return nil, fmt.Errorf("reply status is %d", disperseRes.Result)
	}

	base64RequestID := base64.StdEncoding.EncodeToString(disperseRes.RequestId)

	m.logger.Info("Blob disepersed to EigenDA, now waiting for confirmation", "requestID", base64RequestID)

	var statusRes *disperser.BlobStatusReply
	timeoutTime := time.Now().Add(m.StatusQueryTimeout)
	// Wait before first status check
	time.Sleep(m.StatusQueryRetryInterval)
	for time.Now().Before(timeoutTime) {
		statusRes, err = m.DisperserCli.GetBlobStatus(ctx, &disperser.BlobStatusRequest{
			RequestId: disperseRes.RequestId,
		})
		if err != nil {
			m.logger.Warn("Unable to retrieve blob dispersal status, will retry", "requestID", base64RequestID, "err", err)
		} else if statusRes.Status == disperser.BlobStatus_CONFIRMED || statusRes.Status == disperser.BlobStatus_FINALIZED {
			// TODO(eigenlayer): As long as fault proofs are disabled, we can move on once a blob is confirmed
			// but not yet finalized, without further logic. Once fault proofs are enabled, we will need to update
			// the proposer to wait until the blob associated with an L2 block has been finalized, i.e. the EigenDA
			// contracts on Ethereum have confirmed the full availability of the blob on EigenDA.
			batchHeaderHashHex := fmt.Sprintf("0x%s", hex.EncodeToString(statusRes.Info.BlobVerificationProof.BatchMetadata.BatchHeaderHash))
			m.logger.Info("Successfully dispersed blob to EigenDA", "requestID", base64RequestID, "batchHeaderHash", batchHeaderHashHex)
			return statusRes.Info, nil
		} else if statusRes.Status == disperser.BlobStatus_UNKNOWN ||
			statusRes.Status == disperser.BlobStatus_FAILED {
			m.logger.Error("EigenDA blob dispersal failed in processing", "requestID", base64RequestID, "err", err)
			return nil, fmt.Errorf("eigenDA blob dispersal failed in processing with reply status %d", statusRes.Status)
		} else {
			m.logger.Warn("Still waiting for confirmation from EigenDA", "requestID", base64RequestID)
		}

		// Wait before first status check
		time.Sleep(m.StatusQueryRetryInterval)
	}

	return nil, fmt.Errorf("timed out getting EigenDA status for dispersed blob key: %s", base64RequestID)
}

func (m *EigenDAClient) DisperseBlob(ctx context.Context, txData []byte) ([]byte, error) {
	m.logger.Info("Attempting to disperse blob to EigenDA", "txDataHex", hex.EncodeToString(txData))

	if m.DisperserCli == nil {
		return nil, errors.New("eigendDA disperserCli is nil")
	}

	// encode modulo bn254
	encodedTxData := codec.ConvertByPaddingEmptyByte(txData)

	disperseReq := &disperser.DisperseBlobRequest{
		Data: encodedTxData,
	}

	disperseRes, err := m.DisperserCli.DisperseBlob(ctx, disperseReq)

	if err != nil || disperseRes == nil {
		m.logger.Error("Unable to disperse blob to EigenDA, aborting", "err", err)
		return nil, err
	}

	m.logger.Debug("daClient.DisperseBlob", "disperseRes", disperseRes)
	m.logger.Debug("daClient.DisperseBlob", "disperseRes.Result", disperseRes.Result)
	if disperseRes.Result == disperser.BlobStatus_UNKNOWN ||
		disperseRes.Result == disperser.BlobStatus_FAILED {
		m.logger.Error("Unable to disperse blob to EigenDA, aborting", "err", err)
		return nil, fmt.Errorf("reply status is %d", disperseRes.Result)
	}

	base64RequestID := base64.StdEncoding.EncodeToString(disperseRes.RequestId)

	m.logger.Info("Blob disepersed to EigenDA, now waiting for confirmation", "requestID", base64RequestID)
	return disperseRes.RequestId, nil
}

func (m *EigenDAClient) GetBlobStatus(ctx context.Context, reqID []byte) (disperser.BlobStatus, *disperser.BlobInfo, error) {
	if m.DisperserCli == nil {
		return -1, nil, errors.New("eigendDA disperserCli is nil")
	}

	base64RequestID := base64.StdEncoding.EncodeToString(reqID)
	log.Info("GetBlobStatus", "reqID", reqID, "reqIDBase64", base64RequestID)

	statusRes, err := m.DisperserCli.GetBlobStatus(ctx, &disperser.BlobStatusRequest{
		RequestId: reqID,
	})
	if err != nil {
		m.logger.Warn("Unable to retrieve blob dispersal status, should retry", "requestID", base64RequestID, "err", err)
		return -1, nil, err
	} else if statusRes.Status == disperser.BlobStatus_CONFIRMED || statusRes.Status == disperser.BlobStatus_FINALIZED {
		// TODO(eigenlayer): As long as fault proofs are disabled, we can move on once a blob is confirmed
		// but not yet finalized, without further logic. Once fault proofs are enabled, we will need to update
		// the proposer to wait until the blob associated with an L2 block has been finalized, i.e. the EigenDA
		// contracts on Ethereum have confirmed the full availability of the blob on EigenDA.
		batchHeaderHashHex := fmt.Sprintf("0x%s", hex.EncodeToString(statusRes.Info.BlobVerificationProof.BatchMetadata.BatchHeaderHash))
		m.logger.Info("Successfully dispersed blob to EigenDA", "requestID", base64RequestID, "batchHeaderHash", batchHeaderHashHex)
		return statusRes.Status, statusRes.Info, nil
	} else if statusRes.Status == disperser.BlobStatus_UNKNOWN ||
		statusRes.Status == disperser.BlobStatus_FAILED {
		m.logger.Error("EigenDA blob dispersal failed in processing", "requestID", base64RequestID, "err", err)
		return statusRes.Status, statusRes.Info, fmt.Errorf("eigenDA blob dispersal failed in processing with reply status %d", statusRes.Status)
	}
	m.logger.Warn("Still waiting for confirmation from EigenDA", "requestID", base64RequestID)
	return statusRes.Status, statusRes.Info, nil
}
