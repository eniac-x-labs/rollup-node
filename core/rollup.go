package core

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"

	_errors "github.com/eniac-x-labs/rollup-node/common/errors"
	_config "github.com/eniac-x-labs/rollup-node/config"
	"github.com/eniac-x-labs/rollup-node/x/anytrust"
	"github.com/eniac-x-labs/rollup-node/x/eigenda"
	"github.com/eniac-x-labs/rollup-node/x/nearda"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type RollupModule struct {
	ctx context.Context

	RollupConfig *_config.RollupConfig

	anytrustDA *anytrust.AnytrustDA
	//celestiaDA
	eigenDA eigenda.IEigenDA
	//eip4844
	nearDA nearda.INearDA
}

func NewRollupModule(ctx context.Context, conf *_config.RollupConfig) (RollupInter, error) {

	anytrustDA, err := anytrust.NewAnytrustDA(ctx, conf.AnytrustDAConfig.DAConfig, conf.AnytrustDAConfig.DataSigner)
	if err != nil {
		log.Error("NewAnytrustDA failed", "err", err)
	}

	eigenda, err := eigenda.NewEigenDAClient(conf.EigenDAConfig)
	if err != nil {
		log.Error("NewEigenDA failed", "err", err)
	}

	nearDA, err := nearda.NewNearDAClient(conf.NearDAConfig)
	if err != nil {
		log.Error("NewNearDA failed", "err", err)
	}

	return &RollupModule{
		ctx:          ctx,
		RollupConfig: conf,

		anytrustDA: anytrustDA,
		//celestiaDA: nil,
		eigenDA: eigenda,
		//eip4844:    nil,
		nearDA: nearDA,
	}, nil
}

func (r *RollupModule) RollupWithType(data []byte, daType int) ([]interface{}, error) {
	res := make([]interface{}, 0)
	switch daType {
	case AnytrustType:
		if r.anytrustDA == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "anytrustDA")
			return nil, _errors.DANotPreparedErr
		}

		daCert, err := r.anytrustDA.Store(r.ctx, data, r.RollupConfig.AnytrustDAConfig.DataRetentionTime, nil)
		if err != nil {
			log.Error(_errors.RollupFailedMsg, "da-type", "anytrustDA", "err", err)
			return nil, err
		}
		log.Debug("eigenDA stored data", "daCert.DataHash", fmt.Sprintf("%x", daCert.DataHash))

		res = append(res, daCert)
		return res, nil
	case CelestiaType:

	case EigenDAType:
		if r.eigenDA == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "eigenDA")
			return nil, _errors.DANotPreparedErr
		}

		reqID, err := r.eigenDA.DisperseBlob(r.ctx, data)
		if err != nil {
			log.Error(_errors.RollupFailedMsg, "da-type", "eigenDA", "err", err)
		}
		reqIDBase64 := base64.StdEncoding.EncodeToString(reqID)
		log.Debug("eigenDA stored data", "reqIDBase64", reqIDBase64)

		res = append(res, reqIDBase64)
		return res, nil
	case Eip4844Type:
	case NearDAType:
		if r.nearDA == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "nearDA")
			return nil, _errors.DANotPreparedErr
		}

		txIDByte, err := r.nearDA.Store(data)
		if err != nil {
			log.Error(_errors.RollupFailedMsg, "da-type", "nearDA", "err", err)
			return nil, err
		}
		txid := binary.BigEndian.Uint32(txIDByte[:32])
		log.Debug("nearDA stored data", "txID", txid)
		res = append(res, txid)
		return res, nil
	default:
		log.Error("rollup with unknown da type", "daType", daType, "expected", "[0,4]")
	}
	return nil, _errors.UnknownDATypeErr
}

func (r *RollupModule) GetFromDAWithType(daType int, args ...interface{}) ([]byte, error) {
	switch daType {
	case AnytrustType:
		if r.anytrustDA == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "anytrustDA")
			return nil, _errors.DANotPreparedErr
		}
		if len(args) != 1 {
			log.Error(_errors.WrongArgsNumberErrMsg, "da-type", "anytrustDA", "got", len(args), "expected", 1)
			return nil, _errors.WrongArgsNumberErr
		}

		hashHex := args[0].(string)
		r.anytrustDA.DataAvailabilityServiceReader.String()
		res, err := r.anytrustDA.GetByHash(r.ctx, common.HexToHash(hashHex))
		if err != nil {
			log.Error(_errors.GetFromDAErrMsg, "err", err, "hashHex", hashHex, "da-type", "anytrustDA")
			return nil, err
		}

		log.Debug("get from anytrustDA successfully", "hashHex", hashHex)
		return res, nil
	case CelestiaType:
	case EigenDAType:
		if r.eigenDA == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "eigenDA")
			return nil, _errors.DANotPreparedErr
		}
		if len(args) != 1 {
			log.Error(_errors.WrongArgsNumberErrMsg, "da-type", "eigenDA", "got", len(args), "expected", 1)
			return nil, _errors.WrongArgsNumberErr
		}
		reqIDBase64 := args[0].(string)
		log.Debug("request get from eigenDA", "reqID", reqIDBase64)

		reqIDByte, err := base64.StdEncoding.DecodeString(reqIDBase64)
		log.Error("decode base64 reqID into string failed", "err", err, "reqIDBase64", reqIDBase64, "da-type", "eigenDA")

		status, info, err := r.eigenDA.GetBlobStatus(r.ctx, reqIDByte)
		if err != nil {
			log.Error(_errors.GetFromDAErrMsg, "err", err, "reqIDBase64", reqIDBase64, "da-type", "eigenDA")
			return nil, err
		}
		log.Debug("get from eigenda", "status", status.String(), "reqIDBase64", reqIDBase64)

		batchHeaderHash, blobIndex := info.BlobVerificationProof.GetBatchMetadata().GetBatchHeaderHash(), info.GetBlobVerificationProof().GetBlobIndex()
		res, err := r.eigenDA.RetrieveBlob(r.ctx, batchHeaderHash, blobIndex)
		if err != nil {
			log.Error(_errors.GetFromDAErrMsg, "da-type", "eigenDA")
			return nil, err
		}

		log.Debug("get from eigenda successfully", "reqIDBase64", reqIDBase64)
		return res, nil
	case Eip4844Type:
	case NearDAType:
		if r.nearDA == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "nearDA")
			return nil, _errors.DANotPreparedErr
		}
		if len(args) != 1 {
			log.Error(_errors.WrongArgsNumberErrMsg, "da-type", "nearDA", "got", len(args), "expected", 1)
			return nil, _errors.WrongArgsNumberErr
		}
		frameRefBytes := args[0].([]byte)
		if len(frameRefBytes) < 32 {
			log.Error("nearda arg length incorrect", "length", len(frameRefBytes), "want", "larger than 32")
			return nil, errors.New(fmt.Sprintf("nearda arg length incorrect, expected: larger than 32, got: %d", len(frameRefBytes)))
		}

		result, err := r.nearDA.GetFromDA(frameRefBytes, binary.BigEndian.Uint32(frameRefBytes[:32]))
		if err != nil {
			log.Error(_errors.GetFromDAErrMsg, "da-type", "nearDA", "err", err)
			return nil, err
		}

		log.Debug("get from da successfully")
		return result, nil
	default:
		log.Error("RetrieveFromDAWithType got unknown da type", "daType", daType, "expected", "[0,4]")
	}
	return nil, _errors.UnknownDATypeErr
}
