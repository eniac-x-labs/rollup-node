package core

import (
	_errors "github.com/eniac-x-labs/rollup-node/common/errors"
	"github.com/eniac-x-labs/rollup-node/x/anytrust"
	"github.com/eniac-x-labs/rollup-node/x/eigenda"
	"github.com/eniac-x-labs/rollup-node/x/nearda"
	"github.com/ethereum/go-ethereum/log"
)

type RollupModule struct {
	anytrustDA *anytrust.AnytrustDA
	celestiaDA
	eigenDA *eigenda.EigenDAClient
	eip4844
	nearDA *nearda.NearDAClient
}

func NewRollupModule() (RollupInter, error) {
	anytrustDA, err := anytrust.NewAnytrustDA()
	if err != nil {

	}
	return &RollupModule{
		anytrustDA: anytrustDA,
		celestiaDA: nil,
		eigenDA:    nil,
		eip4844:    nil,
		nearDA:     nil,
	}, nil
}

func (r *RollupModule) RollupWithType(data []byte, daType int) ([]interface{}, error) {
	switch daType {
	case AnytrustType:
	case CelestiaType:
	case EigenDAType:
	case Eip4844Type:
	case NearDAType:
	default:
		log.Error("rollup with unknown da type", "daType", daType, "expected", "[0,4]")
		return nil, _errors.UnknownDATypeErr
	}
}
