package eigenda

import (
	"time"
)

type EigenDAConfig struct {
	// TODO(eigenlayer): Update quorum ID command-line parameters to support passing
	// and arbitrary number of quorum IDs.

	// DaRpc is the HTTP provider URL for the Data Availability node.
	RPC string `toml:"rpc"`

	// The total amount of time that the batcher will spend waiting for EigenDA to confirm a blob
	StatusQueryTimeout time.Duration `toml:"status_query_timeout"`

	// The amount of time to wait between status queries of a newly dispersed blob
	StatusQueryRetryInterval time.Duration `toml:"status_query_retry_interval"`
}

const (
	RpcFlag                      = "rpc"
	StatusQueryTimeoutFlag       = "status_query_timeout"
	StatusQueryRetryIntervalFlag = "status_query_retry_interval"
)

// EigenDAEnvFlags The env flag is like prefix_flag, with all letters in uppercase.
var EigenDAEnvFlags = []string{
	RpcFlag,
	StatusQueryTimeoutFlag,
	StatusQueryRetryIntervalFlag,
}
