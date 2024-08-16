package sdk

import (
	"testing"

	_common "github.com/eniac-x-labs/rollup-node/common"
	"github.com/stretchr/testify/assert"
)

func Test_EigenDA(t *testing.T) {
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
	//arg1 := "MWNjNDc5YmVjMTBmNTFkYjVkMTUzNjJiMzg2ZTNmNGU2ZDhlY2E4MmRlZGViOTAyMWNmYWYyZjNkMzI3ZjJhNS0zMTM3MzEzNzM4MzMzNDM1MzAzMDM4MzIzMDM3MzkzNDM0MzYzODJmMzAyZjMzMzMyZjMxMmYzMzMzMmZlM2IwYzQ0Mjk4ZmMxYzE0OWFmYmY0Yzg5OTZmYjkyNDI3YWU0MWU0NjQ5YjkzNGNhNDk1OTkxYjc4NTJiODU1"
	//arg2 := "OGEyYTVjOWI3Njg4MjdkZTVhOTU1MmMzOGEwNDRjNjY5NTljNjhmNmQyZjIxYjUyNjBhZjU0ZDJmODdkYjgyNy0zMTM3MzEzODMwMzkzMTMyMzczNjM2MzgzOTMzMzczNDM2MzkzMTJmMzEyZjMzMzMyZjMwMmYzMzMzMmZlM2IwYzQ0Mjk4ZmMxYzE0OWFmYmY0Yzg5OTZmYjkyNDI3YWU0MWU0NjQ5YjkzNGNhNDk1OTkxYjc4NTJiODU1"
	arg3 := res[0].(string)
	resByte, err := sdk.RetrieveFromDAWithType(_common.EigenDAType, arg3)
	if err != nil {
		t.Log("3")
		t.Log(err.Error())
	} else {
		t.Log("4")
		t.Logf("%s", resByte)
	}
}

func Test_NearDA(t *testing.T) {
	ast := assert.New(t)
	sdk, err := NewRollupSdk("localhost:9000")
	ast.NoError(err)
	ast.NotNil(sdk)
	t.Log("1")
	data := []byte("rollup data")

	res, err := sdk.RollupWithType(data, _common.NearDAType)
	ast.NoError(err)
	t.Logf("%+v", res)
	t.Log("2")

	resByte, err := sdk.RetrieveFromDAWithType(_common.NearDAType, res[0].(string))
	ast.NoError(err)
	t.Logf("%s", resByte)
}

func Test_Anytrust(t *testing.T) {
	ast := assert.New(t)
	sdk, err := NewRollupSdk("localhost:9000")
	ast.NoError(err)
	ast.NotNil(sdk)
	data := []byte("rollup data")

	res, err := sdk.RollupWithType(data, _common.AnytrustType)
	ast.NoError(err)
	t.Log(res[0].(string))

	resRetrieve, err := sdk.RetrieveFromDAWithType(_common.AnytrustType, res[0].(string))
	ast.NoError(err)
	t.Logf("%s", resRetrieve)
}

func Test_AnytrustCommittee(t *testing.T) {
	ast := assert.New(t)
	sdk, err := NewRollupSdk("localhost:9000")
	ast.NoError(err)
	ast.NotNil(sdk)
	data := []byte("rollup data AnytrustCommitteeType")

	res, err := sdk.RollupWithType(data, _common.AnytrustCommitteeType)
	ast.NoError(err)
	t.Log(res[0].(string))

	resRetrieve, err := sdk.RetrieveFromDAWithType(_common.AnytrustCommitteeType, res[0].(string))
	ast.NoError(err)
	t.Logf("%s", resRetrieve)
}
