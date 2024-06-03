package config

import (
	"testing"

	_das "github.com/eniac-x-labs/anytrustDA/das"

	//_das "github.com/eniac-x-labs/rollup-node/x/anytrust/anytrustDA/das"
	"github.com/stretchr/testify/assert"
)

func Test_PrepareConfig(t *testing.T) {
	ast := assert.New(t)
	conf := &_das.DataAvailabilityConfig{}
	err := PrepareConfig("./", "anytrust", conf, "", nil)
	ast.NoError(err)
	t.Log(conf)
}
