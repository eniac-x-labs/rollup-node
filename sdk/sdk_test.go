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
	t.Log(res[0].(string))
	t.Log("2")

	resByte, err := sdk.RetrieveFromDAWithType(_common.EigenDAType, "MWNjNDc5YmVjMTBmNTFkYjVkMTUzNjJiMzg2ZTNmNGU2ZDhlY2E4MmRlZGViOTAyMWNmYWYyZjNkMzI3ZjJhNS0zMTM3MzEzNzM4MzMzNDM1MzAzMDM4MzIzMDM3MzkzNDM0MzYzODJmMzAyZjMzMzMyZjMxMmYzMzMzMmZlM2IwYzQ0Mjk4ZmMxYzE0OWFmYmY0Yzg5OTZmYjkyNDI3YWU0MWU0NjQ5YjkzNGNhNDk1OTkxYjc4NTJiODU1")
	ast.NoError(err)
	t.Logf("%s", resByte)
	t.Log("3")

	res, err = sdk.RollupWithType(data, _common.NearDAType)
	ast.NoError(err)
	t.Logf("%+v", res)
	t.Log("4")

	resByte, err = sdk.RetrieveFromDAWithType(_common.NearDAType, res[0].(string))
	ast.NoError(err)
	t.Logf("%s", resByte)

}
