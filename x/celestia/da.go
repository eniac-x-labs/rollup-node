package celestia

import (
	"github.com/eniac-x-labs/rollup-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// DerivationVersionCelestia is a byte marker for celestia references submitted
// to the batch inbox address as calldata.
// Mnemonic 0xce = celestia
// version 0xce references are encoded as:
// [8]byte block height ++ [32]byte commitment
// in little-endian encoding.
// see: https://github.com/rollkit/celestia-da/blob/1f2df375fd2fcc59e425a50f7eb950daa5382ef0/celestia.go#L141-L160
const DerivationVersionCelestia = 0xce

type DataSourceConfig struct {
	l1Signer          types.Signer
	batchInboxAddress common.Address
	batcherAddr       common.Address
}

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
