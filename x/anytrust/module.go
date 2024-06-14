package anytrust

import (
	"context"
	"crypto/ecdsa"
	"strings"

	"github.com/eniac-x-labs/anytrustDA/arbstate"
	"github.com/eniac-x-labs/anytrustDA/das"
	"github.com/eniac-x-labs/anytrustDA/util/signature"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type IAnytrustDA interface {
	WriteDA(ctx context.Context, data []byte, retentionTime uint64) (*arbstate.DataAvailabilityCertificate, error)
	ReadDA(ctx context.Context, hashHex string) ([]byte, error)
}

type AnytrustDACommittee struct {
	das.DataAvailabilityServiceWriter
	das.DataAvailabilityServiceReader
	*das.LifecycleManager
}

func NewAnytrustDAWithCommittee(ctx context.Context, daConfig *das.DataAvailabilityConfig, dataSigner signature.DataSignerFunc) (*AnytrustDACommittee, error) {
	//daWriter, daReader, dasLifecycleManager, err := _das.CreateBatchPosterDAS(ctx, daConfig, dataSigner, l1client, deployInfo.SequencerInbox)

	daWriter, daReader, lifeManager, err := das.CreateAggregatorComponents(ctx, daConfig, dataSigner)
	if err != nil {
		return nil, err
	}
	return &AnytrustDACommittee{
		DataAvailabilityServiceWriter: daWriter,
		DataAvailabilityServiceReader: daReader,
		LifecycleManager:              lifeManager,
	}, nil
}

func (a *AnytrustDACommittee) WriteDA(ctx context.Context, data []byte, retentionTime uint64) (*arbstate.DataAvailabilityCertificate, error) {
	return a.Store(ctx, data, retentionTime, nil)
}

func (a *AnytrustDACommittee) ReadDA(ctx context.Context, hashHex string) ([]byte, error) {
	if strings.HasPrefix(hashHex, "0x") {
		hashHex = hashHex[2:]
	}
	return a.GetByHash(ctx, common.HexToHash(hashHex))
}

type AnytrustDA struct {
	writer das.DataAvailabilityServiceWriter //*das.DASRPCClient
	reader *das.RestfulDasClient
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
	if strings.HasPrefix(hashHex, "0x") {
		hashHex = hashHex[2:]
	}
	return a.reader.GetByHash(ctx, common.HexToHash(hashHex))
}
