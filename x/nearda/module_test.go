package nearda

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

// create contract account: near account create-account sponsor-by-faucet-service wwqcontract.testnet autogenerate-new-keypair save-to-keychain network-config testnet create
// deploy contract： near contract deploy wwqcontract.testnet use-file /Users/wwq/go/github.com/rollup-data-availability/target/wasm32-unknown-unknown/release/near_da_blob_store.wasm with-init-call new json-args {} prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' network-config testnet sign-with-keychain send
// get DA_KEY: near account export-account wwqcontract.testnet using-private-key network-config testnet
const (
	DA_ACCOUNT  = "wwqcontract.testnet"
	DA_CONTRACT = "wwqcontract.testnet"
	// 下面的key 可以从near -> account -> export-account获得
	DA_KEY = "ed25519:4btKLuh9xbrybQUYaJJTeKb1cC35kYtpVxsGByT1H9ixR8PaCoCHHfHq1tEVm4ABG9fckSEDcWcxVzhc3J3C5tNv"
	//DA_KEY      = "ed25519:DXVvU1N8TFdBjm7HiyL5MFLYUDCYewKBhKHbEap6BMcRh5CFNrDtqR1s3QbXQnZAv5Qj4iJHjwGrh6Hfwetxt2p"
	DA_NETWORK = "Testnet"
	DA_NS      = 1
)

var testConf = NearADConfig{
	Account:  DA_ACCOUNT,
	Contract: DA_CONTRACT,
	Key:      DA_KEY,
	Network:  DA_NETWORK,
	Ns:       DA_NS,
}

func Test_NewNearDAClient(t *testing.T) {
	ast := assert.New(t)
	//key := DA_KEY
	//keyBase64 := base64.StdEncoding.EncodeToString([]byte(key))
	//t.Logf("keyBase64: %s", keyBase64)
	//
	//testConf.Key = keyBase64
	//dacli, err := near.NewConfig(DA_ACCOUNT, DA_CONTRACT, DA_KEY, DA_NETWORK, DA_NS)
	dacli, err := NewNearDAClient(&testConf)
	ast.NoError(err)
	ast.NotNil(dacli)

	t.Log("==================== Start Submitting ========================")
	dataStr := "hello nearDA"
	res, err := dacli.Store([]byte(dataStr))
	ast.NoError(err)
	t.Logf("ForceSubmit result: %s", res)
	t.Logf("len of submit result: %d", len(res))

	t.Log("==================== get from da ========================")
	txid := binary.BigEndian.Uint32(res[:32])
	t.Logf("txid: %d", txid)

	getRes, err := dacli.GetFromDA(res, txid)
	ast.NoError(err)
	t.Logf("getRes: %s", getRes)
}
