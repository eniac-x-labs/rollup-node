package signer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

func SignerFactoryFromPrivateKey(privateKey string) (SignerFactory, common.Address, error) {
	var signer SignerFactory
	var fromAddress common.Address
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

	return signer, fromAddress, nil
}
