package anytrust

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"io"
	"strings"
	"time"

	"github.com/eniac-x-labs/rollup-node/x/anytrust/anytrustDA/arbstate"
	"github.com/eniac-x-labs/rollup-node/x/anytrust/anytrustDA/das"
	"github.com/eniac-x-labs/rollup-node/x/anytrust/anytrustDA/util/signature"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type IAnytrustDA interface {
	WriteDA(ctx context.Context, data []byte, retentionTime uint64) (*arbstate.DataAvailabilityCertificate, error)
	ReadDA(ctx context.Context, hashHex string) ([]byte, error)
}

//	type AnytrustDA struct {
//		writer das.DataAvailabilityServiceWriter
//		reader das.DataAvailabilityServiceReader
//		*das.LifecycleManager
//	}
type AnytrustDA struct {
	writer das.DataAvailabilityServiceWriter //*das.DASRPCClient
	reader *das.RestfulDasClient
}

type AnytrustConfig struct {
	RpcUrl            string `toml:"rpcUrl"`
	RestfulUrl        string `toml:"restfulUrl"`
	DataRetentionTime uint64 `toml:"dataRetentionTime"`

	RandomMessageSize  int           `toml:"randomMessageSize"`
	DASRetentionPeriod time.Duration `toml:"dasRetentionPeriod"`
	SigningKey         string        `toml:"signingKey"`
	//SigningWallet         string        `toml:"signingWallet"`
	//SigningWalletPassword string        `toml:"signingWalletPassword"`
}

func NewAnytrustDA(config *AnytrustConfig) (IAnytrustDA, error) {
	rpcClient, err := das.NewDASRPCClient(config.RpcUrl)
	if err != nil {
		return nil, err
	}

	var dasClient das.DataAvailabilityServiceWriter = rpcClient
	if config.SigningKey != "" {
		var privateKey *ecdsa.PrivateKey
		if config.SigningKey[:2] == "0x" {
			privateKey, err = crypto.HexToECDSA(config.SigningKey[2:])
			if err != nil {
				return nil, err
			}
		} else {
			privateKey, err = crypto.LoadECDSA(config.SigningKey)
			if err != nil {
				return nil, err
			}
		}
		signer := signature.DataSignerFromPrivateKey(privateKey)

		dasClient, err = das.NewStoreSigningDAS(dasClient, signer)
		if err != nil {
			return nil, err
		}
	}
	//} else if config.SigningWallet != "" {
	//	walletConf := &genericconf.WalletConfig{
	//		Pathname:      config.SigningWallet,
	//		Password:      config.SigningWalletPassword,
	//		PrivateKey:    "",
	//		Account:       "",
	//		OnlyCreateKey: false,
	//	}
	//	_, signer, err := util.OpenWallet("datool", walletConf, nil)
	//	if err != nil {
	//		return err
	//	}
	//	dasClient, err = das.NewStoreSigningDAS(dasClient, signer)
	//	if err != nil {
	//		return err
	//	}
	//}

	reader, err := das.NewRestfulDasClientFromURL(config.RestfulUrl)
	if err != nil {
		return nil, err
	}
	return &AnytrustDA{
		writer: dasClient,
		reader: reader,
	}, nil
}

func (a *AnytrustDA) WriteDA(ctx context.Context, data []byte, retentionTime uint64) (*arbstate.DataAvailabilityCertificate, error) {
	return a.writer.Store(ctx, data, retentionTime, nil)
}

func (a *AnytrustDA) ReadDA(ctx context.Context, hashHex string) ([]byte, error) {
	var (
		decodedHash []byte
		err         error
	)
	if strings.HasPrefix(hashHex, "0x") {
		decodedHash, err = hexutil.Decode(hashHex)
		if err != nil {
			return nil, err
		}
	} else {
		hashDecoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(hashHex)))
		decodedHash, err = io.ReadAll(hashDecoder)
		if err != nil {
			return nil, err
		}
	}

	return a.reader.GetByHash(ctx, common.BytesToHash(decodedHash))
}
