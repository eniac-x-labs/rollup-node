package sdk

import (
	"testing"

	_common "github.com/eniac-x-labs/rollup-node/common"
	"github.com/stretchr/testify/assert"
)

func Test_Sdk(t *testing.T) {
	ast := assert.New(t)
	sdk, err := NewRollupSdk("localhost:9000")
	ast.NoError(err)
	ast.NotNil(sdk)
	t.Log("1")
	data := []byte("rollup data")
	res, err := sdk.RollupWithType(data, _common.EigenDAType)
	ast.NoError(err)
	t.Logf("%+v", res)
	t.Log("2")

	resByte, err := sdk.RetrieveFromDAWithType(_common.EigenDAType, "MWNjNDc5YmVjMTBmNTFkYjVkMTUzNjJiMzg2ZTNmNGU2ZDhlY2E4MmRlZGViOTAyMWNmYWYyZjNkMzI3ZjJhNS0zMTM3MzEzNzM1MzkzMjM4MzIzOTM5MzYzNDM1MzAzNjMxMzYzODJmMzAyZjMz")
	ast.NoError(err)
	t.Logf("%x", resByte)
	t.Log("3")

	res, err = sdk.RollupWithType(data, _common.NearDAType)
	ast.NoError(err)
	t.Logf("%+v", res)
	t.Log("4")

	resByte, err = sdk.RetrieveFromDAWithType(_common.NearDAType, "MWNjNDc5YmVjMTBmNTFkYjVkMTUzNjJiMzg2ZTNmNGU2ZDhlY2E4MmRlZGViOTAyMWNmYWYyZjNkMzI3ZjJhNS0zMTM3MzEzNzM1MzkzMjM4MzIzOTM5MzYzNDM1MzAzNjMxMzYzODJmMzAyZjMz")
	ast.NoError(err)
	t.Logf("%x", resByte)

}
