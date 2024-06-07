package eip4844

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	eth "github.com/eniac-x-labs/rollup-node/eth-serivce"
)

type blobOrCalldata struct {
	// union type. exactly one of calldata or blob should be non-nil
	blob     *eth.Blob
	calldata *eth.Data
}

type DataSourceConfig struct {
	l1Signer          types.Signer
	batchInboxAddress common.Address
	batcherAddr       common.Address
}

func dataAndHashesFromTxs(txs types.Transactions, config *DataSourceConfig, log log.Logger) ([]blobOrCalldata, []eth.IndexedBlobHash) {
	data := []blobOrCalldata{}
	var hashes []eth.IndexedBlobHash
	blobIndex := 0 // index of each blob in the block's blob sidecar
	for _, tx := range txs {
		// skip any non-batcher transactions
		if !isValidBatchTx(tx, config.l1Signer, config.batchInboxAddress, config.batcherAddr) {
			blobIndex += len(tx.BlobHashes())
			continue
		}
		// handle non-blob batcher transactions by extracting their calldata
		if tx.Type() != types.BlobTxType {
			calldata := eth.Data(tx.Data())
			data = append(data, blobOrCalldata{nil, &calldata})
			continue
		}
		// handle blob batcher transactions by extracting their blob hashes, ignoring any calldata.
		if len(tx.Data()) > 0 {
			log.Warn("blob tx has calldata, which will be ignored", "txhash", tx.Hash())
		}
		for _, h := range tx.BlobHashes() {
			idh := eth.IndexedBlobHash{
				Index: uint64(blobIndex),
				Hash:  h,
			}
			hashes = append(hashes, idh)
			data = append(data, blobOrCalldata{nil, nil}) // will fill in blob pointers after we download them below
			blobIndex += 1
		}
	}
	return data, hashes
}

// isValidBatchTx returns true if:
//  1. the transaction has a To() address that matches the batch inbox address, and
//  2. the transaction has a valid signature from the batcher address
func isValidBatchTx(tx *types.Transaction, l1Signer types.Signer, batchInboxAddr, batcherAddr common.Address) bool {
	to := tx.To()
	if to == nil || *to != batchInboxAddr {
		return false
	}
	seqDataSubmitter, err := l1Signer.Sender(tx) // optimization: only derive sender if To is correct
	if err != nil {
		log.Warn("tx in inbox with invalid signature", "hash", tx.Hash(), "err", err)
		return false
	}
	// some random L1 user might have sent a transaction to our batch inbox, ignore them
	if seqDataSubmitter != batcherAddr {
		log.Warn("tx in inbox with unauthorized submitter", "addr", seqDataSubmitter, "hash", tx.Hash(), "err", err)
		return false
	}
	return true
}
