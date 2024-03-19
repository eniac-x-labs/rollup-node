package signer

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/eniac-x-labs/rollup-node/x/eip4844/log"
)

func PrivateKeySignerFn(key *ecdsa.PrivateKey, chainID *big.Int) bind.SignerFn {
	from := crypto.PubkeyToAddress(key.PublicKey)
	signer := types.LatestSignerForChainID(chainID)
	return func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		if address != from {
			return nil, bind.ErrNotAuthorized
		}
		signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
		if err != nil {
			return nil, err
		}
		return tx.WithSignature(signer, signature)
	}
}

// SignerFn is a generic transaction signing function. It may be a remote signer so it takes a context.
// It also takes the address that should be used to sign the transaction with.
type SignerFn func(context.Context, common.Address, *types.Transaction) (*types.Transaction, error)

// SignerFactory creates a SignerFn that is bound to a specific ChainID
type SignerFactory func(chainID *big.Int) SignerFn

// SignerFactoryFromConfig considers three ways that signers are created & then creates single factory from those config options.
// It can either take a remote signer (via opsigner.CLIConfig) or it can be provided either a mnemonic + derivation path or a private key.
// It prefers the remote signer, then the mnemonic or private key (only one of which can be provided).
func SignerFactoryFromConfig(l log.Logger, privateKey string, signerConfig CLIConfig) (SignerFactory, common.Address, error) {
	var signer SignerFactory
	var fromAddress common.Address
	if signerConfig.Enabled() {
		signerClient, err := NewSignerClientFromConfig(l, signerConfig)
		if err != nil {
			l.Error("Unable to create Signer Client", "error", err)
			return nil, common.Address{}, fmt.Errorf("failed to create the signer client: %w", err)
		}
		fromAddress = common.HexToAddress(signerConfig.Address)
		signer = func(chainID *big.Int) SignerFn {
			return func(ctx context.Context, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
				if !bytes.Equal(address[:], fromAddress[:]) {
					return nil, fmt.Errorf("attempting to sign for %s, expected %s: ", address, signerConfig.Address)
				}
				return signerClient.SignTransaction(ctx, chainID, address, tx)
			}
		}
	} else {
		var privKey *ecdsa.PrivateKey
		var err error

		if privateKey == "" {
			return nil, common.Address{}, fmt.Errorf("failed to create a wallet: %w", err)
		} else {
			privKey, err = crypto.HexToECDSA(strings.TrimPrefix(privateKey, "0x"))
			if err != nil {
				return nil, common.Address{}, fmt.Errorf("failed to parse the private key: %w", err)
			}
		}
		fromAddress = crypto.PubkeyToAddress(privKey.PublicKey)
		signer = func(chainID *big.Int) SignerFn {
			s := PrivateKeySignerFn(privKey, chainID)
			return func(_ context.Context, addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
				return s(addr, tx)
			}
		}
	}

	return signer, fromAddress, nil
}
