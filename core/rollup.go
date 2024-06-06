package core

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/eniac-x-labs/anytrustDA/das"
	_common "github.com/eniac-x-labs/rollup-node/common"
	"github.com/eniac-x-labs/rollup-node/common/cliapp"
	_errors "github.com/eniac-x-labs/rollup-node/common/errors"
	"github.com/eniac-x-labs/rollup-node/config"
	_log "github.com/eniac-x-labs/rollup-node/log"
	_rpc "github.com/eniac-x-labs/rollup-node/rpc"
	"github.com/eniac-x-labs/rollup-node/x/anytrust"
	"github.com/eniac-x-labs/rollup-node/x/celestia"
	"github.com/eniac-x-labs/rollup-node/x/eigenda"
	"github.com/eniac-x-labs/rollup-node/x/eip4844"
	"github.com/eniac-x-labs/rollup-node/x/nearda"
	"github.com/urfave/cli/v2"

	_config "github.com/eniac-x-labs/rollup-node/config"

	"github.com/ethereum/go-ethereum/log"
)

var ErrAlreadyStopped = errors.New("already stopped")

type RollupModule struct {
	ctx context.Context

	RollupConfig *_config.RollupConfig

	anytrustDA anytrust.IAnytrustDA
	celestiaDA *celestia.CelestiaRollup
	eigenDA    eigenda.IEigenDA
	eip4844    *eip4844.Eip4844Rollup
	nearDA     nearda.INearDA
	stopped    atomic.Bool
	Log        _log.Logger
}

func (r *RollupModule) Start(ctx context.Context) error {
	return nil
}

func (r *RollupModule) Stop(ctx context.Context) error {
	if r.stopped.Load() {
		return ErrAlreadyStopped
	}
	r.Log.Info("Stopping rollup node service")

	r.stopped.Store(true)
	r.Log.Info("rollup node service stopped")
	return nil
}

func (r *RollupModule) Stopped() bool {
	return r.stopped.Load()
}

func RunRollupModuleForCLI(cliCtx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	rollupModule, err := NewRollupModule_2(cliCtx)
	if err != nil {
		log.Error("Run Rollup Module failed", "err", err)
		return nil, err
	}

	rpcAddress := cliCtx.String("rpcAddress")
	if len(rpcAddress) != 0 {
		go _rpc.NewAndStartRollupRpcServer(cliCtx.Context, rpcAddress, rollupModule)
	}

	return rollupModule, nil
}

func NewRollupModuleWithConfig(ctx context.Context, conf *_config.RollupConfig) (_rpc.RollupInter, error) {
	anytrustDA, err := anytrust.NewAnytrustDA(conf.AnytrustDAConfig)
	if err != nil {
		log.Error("NewAnytrustDA failed", "err", err)
	}

	//celestiaDA, err := celestia.NewCelestiaRollup()
	eigenda, err := eigenda.NewEigenDAClient(conf.EigenDAConfig)
	if err != nil {
		log.Error("NewEigenDA failed", "err", err)
	}

	//eip4844, err := eip4844.NewEip4844Rollup(cliCtx, logger)
	//if err != nil {
	//
	//}

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
		//eip4844: eip4844,
		nearDA: nearDA,
	}, nil
}

func NewRollupModule_2(cliCtx *cli.Context) (*RollupModule, error) {
	cfg, err := config.NewConfig(cliCtx) // celestia & eip4844
	if err != nil {
		return nil, err
	}

	logger := _log.NewLogger(_log.AppOut(cliCtx), cfg.LogConfig).New("rollup-node")
	_log.SetGlobalLogHandler(logger.GetHandler())

	celestiaDA, err := celestia.NewCelestiaRollup(cliCtx, logger)
	if err != nil {
		log.Error("NewCelestiaRollup failed", "err", err)
		return nil, err
	}
	eip4844, err := eip4844.NewEip4844Rollup(cliCtx, logger)
	if err != nil {
		log.Error("NewEip4844Rollup failed", "err", err)
		return nil, err
	}

	conf := _config.NewRollupConfig() // config for anytrust & eigenda & nearda
	anytrustDA, err := anytrust.NewAnytrustDA(conf.AnytrustDAConfig)
	if err != nil {
		log.Error("NewAnytrustRollup failed", "err", err)
		return nil, err
	}

	eigenDA, err := eigenda.NewEigenDAClient(conf.EigenDAConfig)
	if err != nil {
		log.Error("NewEigenRollup failed", "err", err)
		return nil, err
	}

	nearDA, err := nearda.NewNearDAClient(conf.NearDAConfig)
	if err != nil {
		log.Error("NewNearRollup failed", "err", err)
		return nil, err
	}

	return &RollupModule{
		anytrustDA: anytrustDA,
		celestiaDA: celestiaDA,
		eigenDA:    eigenDA,
		eip4844:    eip4844,
		nearDA:     nearDA,
		Log:        logger,
	}, nil

}

func (r *RollupModule) RollupWithType(data []byte, daType int) ([]interface{}, error) {
	res := make([]interface{}, 0)
	switch daType {
	case _common.AnytrustType:
		if r.anytrustDA == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "anytrustDA")
			return nil, _errors.DANotPreparedErr
		}
		daCert, err := r.anytrustDA.WriteDA(r.ctx, data, r.RollupConfig.AnytrustDAConfig.DataRetentionTime)
		if err != nil {
			log.Error(_errors.RollupFailedMsg, "da-type", "anytrustDA", "err", err)
			return nil, err
		}
		log.Debug("eigenDA stored data", "daCert.DataHash", fmt.Sprintf("%x", daCert.DataHash))

		das.Serialize(daCert)
		res = append(res, daCert.DataHash)
		res = append(res, das.Serialize(daCert))
		return res, nil

	case CelestiaType:
		if r.celestiaDA == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "celestiaDA")
			return nil, _errors.DANotPreparedErr
		}

		txHash, err := r.celestiaDA.SendTransaction(r.ctx, data)
		if err != nil {
			log.Error(_errors.RollupFailedMsg, "da-type", "celestiaDA", "err", err)
			return nil, err
		}
		txHashStr := fmt.Sprintf("0x%s", hex.EncodeToString(txHash))
		log.Debug("celestiaDA stored data", "txHash", txHashStr)

		res = append(res, txHashStr)
		return res, nil
	case _common.EigenDAType:
		if r.eigenDA == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "eigenDA")
			return nil, _errors.DANotPreparedErr
		}

		reqID, err := r.eigenDA.DisperseBlob(r.ctx, data)
		if err != nil {
			log.Error(_errors.RollupFailedMsg, "da-type", "eigenDA", "err", err)
			return nil, err
		}
		reqIDBase64 := base64.StdEncoding.EncodeToString(reqID)
		log.Debug("eigenDA stored data", "reqIDBase64", reqIDBase64)

		res = append(res, reqIDBase64)
		return res, nil
	case _common.Eip4844Type:
		if r.eip4844 == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "eip4844")
			return nil, _errors.DANotPreparedErr
		}

		txHash, err := r.eip4844.SendTransaction(data)
		if err != nil {
			log.Error(_errors.RollupFailedMsg, "da-type", "eip4844", "err", err)
			return nil, err
		}
		txHashStr := fmt.Sprintf("0x%s", hex.EncodeToString(txHash))
		log.Debug("eip4844 stored data", "txHash", txHashStr)
		res = append(res, txHashStr)
		return res, nil
	case _common.NearDAType:
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

func (r *RollupModule) RetrieveFromDAWithType(daType int, args interface{}) ([]byte, error) {
	switch daType {
	case _common.AnytrustType:
		if r.anytrustDA == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "anytrustDA")
			return nil, _errors.DANotPreparedErr
		}
		hashHex := args.(string)
		res, err := r.anytrustDA.ReadDA(r.ctx, hashHex)
		if err != nil {
			log.Error(_errors.GetFromDAErrMsg, "err", err, "hashHex", hashHex, "da-type", "anytrustDA")
			return nil, err
		}

		log.Debug("get from anytrustDA successfully", "hashHex", hashHex)
		return res, nil
	case _common.CelestiaType:
		if r.celestiaDA == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "celestiaDA")
			return nil, _errors.DANotPreparedErr
		}
		reqTxHashStr := args.(string)
		log.Debug("request get from celestiaDA", "reqTxHashStr", reqTxHashStr)

		res, err := r.celestiaDA.DataFromEVMTransactions(reqTxHashStr)
		if err != nil {
			log.Error(_errors.GetFromDAErrMsg, "err", err, "reqTxHashStr", reqTxHashStr, "da-type", "celestiaDA")
			return nil, err
		}

		log.Debug("get from celestiaDA successfully", "reqTxHashStr", reqTxHashStr)
		return res, nil

	case _common.EigenDAType:
		if r.eigenDA == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "eigenDA")
			return nil, _errors.DANotPreparedErr
		}
		reqIDBase64, ok := args.(string)
		if !ok {
			log.Error("args[0] is not string type")
			return nil, _errors.WrongArgTypeErr
		}
		log.Debug("request get from eigenDA", "reqID", reqIDBase64)

		reqIDByte, err := base64.StdEncoding.DecodeString(reqIDBase64)
		log.Error("decode base64 reqID into string failed", "err", err, "reqIDBase64", reqIDBase64, "da-type", "eigenDA")

		status, info, err := r.eigenDA.GetBlobStatus(r.ctx, reqIDByte)
		if err != nil {
			log.Error(_errors.GetFromDAErrMsg, "err", err, "reqIDBase64", reqIDBase64, "da-type", "eigenDA")
			return nil, err
		}
		log.Debug("get from eigenDA", "status", status.String(), "reqIDBase64", reqIDBase64)

		batchHeaderHash, blobIndex := info.BlobVerificationProof.GetBatchMetadata().GetBatchHeaderHash(), info.GetBlobVerificationProof().GetBlobIndex()
		res, err := r.eigenDA.RetrieveBlob(r.ctx, batchHeaderHash, blobIndex)
		if err != nil {
			log.Error(_errors.GetFromDAErrMsg, "da-type", "eigenDA")
			return nil, err
		}

		log.Debug("get from eigenDA successfully", "reqIDBase64", reqIDBase64)
		return res, nil
	case _common.Eip4844Type:
		if r.eip4844 == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "eip4844")
			return nil, _errors.DANotPreparedErr
		}
		reqTxHashStr := args.(string)
		log.Debug("request get from eip4844", "reqTxHashStr", reqTxHashStr)

		res, err := r.eip4844.DataFromEVMTransactions(r.ctx, reqTxHashStr)
		if err != nil {
			log.Error(_errors.GetFromDAErrMsg, "err", err, "reqTxHashStr", reqTxHashStr, "da-type", "eip4844")
			return nil, err
		}

		log.Debug("get from eip4844 successfully", "reqTxHashStr", reqTxHashStr)
		return res, nil
	case _common.NearDAType:
		if r.nearDA == nil {
			log.Error(_errors.DANotPreparedErrMsg, "da-type", "nearDA")
			return nil, _errors.DANotPreparedErr
		}
		frameRefBytes := args.([]byte)
		if len(frameRefBytes) < 32 {
			log.Error("nearda arg length incorrect", "length", len(frameRefBytes), "want", "larger than 32")
			return nil, errors.New(fmt.Sprintf("nearda arg length incorrect, expected: larger than 32, got: %d", len(frameRefBytes)))
		}

		result, err := r.nearDA.GetFromDA(frameRefBytes, binary.BigEndian.Uint32(frameRefBytes[:32]))
		if err != nil {
			log.Error(_errors.GetFromDAErrMsg, "da-type", "nearDA", "err", err)
			return nil, err
		}

		log.Debug("get from nearDA successfully")
		return result, nil
	default:
		log.Error("RetrieveFromDAWithType got unknown da type", "daType", daType, "expected", "[0,4]")
	}
	return nil, _errors.UnknownDATypeErr
}
