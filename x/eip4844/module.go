package eip4844

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"

	"github.com/eniac-x-labs/rollup-node/client"
	"github.com/eniac-x-labs/rollup-node/config"
	eth "github.com/eniac-x-labs/rollup-node/eth-serivce"
	"github.com/eniac-x-labs/rollup-node/log"
	"github.com/eniac-x-labs/rollup-node/signer"
)

var ErrAlreadyStopped = errors.New("already stopped")

type Eip4844Rollup struct {
	Eip4844Config  CLIConfig
	Config         *config.CLIConfig
	l1BeaconClient *eth.L1BeaconClient
	Log            log.Logger
	ethClients     client.EthClient
	Signer         signer.SignerFn
	From           common.Address
	stopped        atomic.Bool
	driverCtx      context.Context
}

func (e *Eip4844Rollup) Stop(ctx context.Context) error {
	if e.stopped.Load() {
		return ErrAlreadyStopped
	}

	e.Log.Info("Stopping eip4844 rollup service")

	e.stopped.Store(true)
	e.Log.Info("eip4844 rollup service stopped")

	return nil
}

func (e *Eip4844Rollup) Stopped() bool {
	return e.stopped.Load()
}

func NewEip4844Rollup(cliCtx *cli.Context, logger log.Logger) (*Eip4844Rollup, error) {

	cfg, err := config.NewConfig(cliCtx)
	if err != nil {
		return nil, err
	}
	eip4844Config := ReadCLIConfig(cliCtx)

	return Eip4844ServiceFromCLIConfig(cliCtx.Context, cfg, eip4844Config, logger)
}

func Eip4844ServiceFromCLIConfig(ctx context.Context, cfg *config.CLIConfig, eip4844Config CLIConfig, logger log.Logger) (*Eip4844Rollup, error) {
	var e Eip4844Rollup
	if err := e.initFromCLIConfig(ctx, cfg, eip4844Config, logger); err != nil {
		return nil, errors.Join(err, e.Stop(ctx)) // try to clean up our failed initialization attempt
	}
	return &e, nil
}

func (e *Eip4844Rollup) initFromCLIConfig(ctx context.Context, cfg *config.CLIConfig, eip4844Config CLIConfig, logger log.Logger) error {

	e.Config = cfg
	e.Eip4844Config = eip4844Config
	e.Log = logger

	signerFactory, from, err := signer.SignerFactoryFromPrivateKey(eip4844Config.PrivateKey)
	if err != nil {
		log.Error(fmt.Errorf("could not init signer: %w", err).Error())
		return err
	}
	e.Signer = signerFactory(eip4844Config.L1ChainID)
	e.From = from

	bCl := client.NewBasicHTTPClient(eip4844Config.L1BeaconAddr)
	beaconCfg := eth.L1BeaconClientConfig{
		FetchAllSidecars: eip4844Config.ShouldFetchAllSidecars,
	}
	e.l1BeaconClient = eth.NewL1BeaconClient(bCl, beaconCfg)

	e.driverCtx = ctx

	return nil
}

// SendTransaction creates & submits a transaction to the batch inbox address with the given `txData`.
// It currently uses the underlying `txmgr` to handle transaction sending & price management.
// This is a blocking method. It should not be called concurrently.
func (e *Eip4844Rollup) SendTransaction(data []byte) ([]byte, error) {
	// Do the gas estimation offline. A value of 0 will cause the [txmgr] to estimate the gas limit.

	var candidate *eth.TxCandidate
	if e.Eip4844Config.UseBlobs {
		var err error
		if candidate, err = e.blobTxCandidate(data); err != nil {
			// We could potentially fall through and try a calldata tx instead, but this would
			// likely result in the chain spending more in gas fees than it is tuned for, so best
			// to just fail. We do not expect this error to trigger unless there is a serious bug
			// or configuration issue.
			return nil, fmt.Errorf("could not create blob tx candidate: %w", err)
		}
	} else {
		candidate = e.calldataTxCandidate(data)
	}

	intrinsicGas, err := core.IntrinsicGas(data, nil, false, true, true, false)
	if err != nil {
		e.Log.Error("Failed to calculate intrinsic gas", "err", err)
	} else {
		candidate.GasLimit = intrinsicGas
	}

	tx, err := e.craftTx(*candidate)
	if err != nil {
		e.Log.Error("Failed to create a transaction", "err", err)
		return nil, err
	}

	signTx, err := e.Signer(e.driverCtx, e.From, tx)
	if err != nil {
		e.Log.Error("Failed to sign a transaction", "err", err)
		return nil, err
	}

	err = e.ethClients.SendTransaction(e.driverCtx, signTx)
	if err != nil {
		e.Log.Error("Failed to send transaction", "err", err)
		return nil, err
	}

	return signTx.Hash().Bytes(), nil
}

func (e *Eip4844Rollup) blobTxCandidate(data []byte) (*eth.TxCandidate, error) {
	var b eth.Blob
	if err := b.FromData(data); err != nil {
		return nil, fmt.Errorf("data could not be converted to blob: %w", err)
	}
	return &eth.TxCandidate{
		To:    &e.Eip4844Config.DSConfig.batchInboxAddress,
		Blobs: []*eth.Blob{&b},
	}, nil
}

func (e *Eip4844Rollup) calldataTxCandidate(data []byte) *eth.TxCandidate {
	e.Log.Info("building Calldata transaction candidate", "size", len(data))
	return &eth.TxCandidate{
		To:     &e.Eip4844Config.DSConfig.batchInboxAddress,
		TxData: data,
	}
}

func (e *Eip4844Rollup) DataFromEVMTransactions(ctx context.Context, txHashStr string) (data eth.Data, err error) {
	var datas []eth.Data
	var txs types.Transactions

	tx, header, err := e.getTransactionAndBlockByTxHash(ctx, txHashStr)
	if err != nil {
		log.Error("failed to get transaction and block by tx hash", "tx_hash", txHashStr, "err", err)
		return nil, err
	}
	txs = append(txs, tx)

	_, hashes := dataAndHashesFromTxs(txs, e.Eip4844Config.DSConfig, e.Log)
	if len(hashes) == 0 {
		// there are no blobs to fetch so we can return immediately
		return nil, nil
	}

	ref := eth.L1BlockRef{
		Hash:       header.Hash(),
		Number:     header.Number.Uint64(),
		ParentHash: header.ParentHash,
		Time:       header.Time,
	}

	blobs, err := e.l1BeaconClient.GetBlobs(e.driverCtx, ref, hashes)
	if errors.Is(err, ethereum.NotFound) {
		// If the L1 block was available, then the blobs should be available too. The only
		// exception is if the blob retention window has expired, which we will ultimately handle
		// by failing over to a blob archival service.
		return nil, fmt.Errorf("failed to fetch blobs: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("failed to fetch blobs: %w", err)
	}

	for _, blob := range blobs {
		data, err := blob.ToData()
		if err != nil {
			return nil, fmt.Errorf("decodes the blob into raw byte data failed: %w", err)
		}

		datas = append(datas, data)
	}

	return datas[0], nil
}

func (e *Eip4844Rollup) craftTx(candidate eth.TxCandidate) (*types.Transaction, error) {
	e.Log.Debug("crafting Transaction", "blobs", len(candidate.Blobs), "calldata_size", len(candidate.TxData))

	gasLimit := candidate.GasLimit

	var sidecar *types.BlobTxSidecar
	var blobHashes []common.Hash
	var err error
	if len(candidate.Blobs) > 0 {
		if candidate.To == nil {
			return nil, errors.New("blob txs cannot deploy contracts")
		}
		if sidecar, blobHashes, err = MakeSidecar(candidate.Blobs); err != nil {
			return nil, fmt.Errorf("failed to make sidecar: %w", err)
		}
	}

	var txMessage types.TxData
	if sidecar != nil {
		message := &types.BlobTx{
			To:         *candidate.To,
			Data:       candidate.TxData,
			Gas:        gasLimit,
			BlobHashes: blobHashes,
			Sidecar:    sidecar,
		}
		txMessage = message
	}
	return types.NewTx(txMessage), err
}

// MakeSidecar builds & returns the BlobTxSidecar and corresponding blob hashes from the raw blob
// data.
func MakeSidecar(blobs []*eth.Blob) (*types.BlobTxSidecar, []common.Hash, error) {
	sidecar := &types.BlobTxSidecar{}
	blobHashes := []common.Hash{}
	for i, blob := range blobs {
		rawBlob := *blob.KZGBlob()
		sidecar.Blobs = append(sidecar.Blobs, rawBlob)
		commitment, err := kzg4844.BlobToCommitment(rawBlob)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot compute KZG commitment of blob %d in tx candidate: %w", i, err)
		}
		sidecar.Commitments = append(sidecar.Commitments, commitment)
		proof, err := kzg4844.ComputeBlobProof(rawBlob, commitment)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot compute KZG proof for fast commitment verification of blob %d in tx candidate: %w", i, err)
		}
		sidecar.Proofs = append(sidecar.Proofs, proof)
		blobHashes = append(blobHashes, eth.KZGToVersionedHash(commitment))
	}
	return sidecar, blobHashes, nil
}

func (e *Eip4844Rollup) getTransactionAndBlockByTxHash(ctx context.Context, txHashStr string) (*types.Transaction, *types.Header, error) {
	tx, err := e.ethClients.TxByHash(common.HexToHash(txHashStr))
	if err != nil {
		e.Log.Error("failed to get transaction", "tx_hash", txHashStr)
		return nil, nil, err
	}

	receipt, err := e.ethClients.TxReceiptDetailByHash(tx.Hash())
	if err != nil {
		e.Log.Error("failed to get transaction receipt", "tx_hash", txHashStr)
		return nil, nil, err
	}

	header, err := e.ethClients.HeaderByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		e.Log.Error("failed to get block header by number", "number", receipt.BlockNumber)
		return nil, nil, err
	}

	return tx, header, nil
}
