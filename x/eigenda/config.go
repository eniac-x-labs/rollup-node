package eigenda

import (
	"time"
)

type EigenDAConfig struct {
	// TODO(eigenlayer): Update quorum ID command-line parameters to support passing
	// and arbitrary number of quorum IDs.

	// DaRpc is the HTTP provider URL for the Data Availability node.
	RPC string

	// The total amount of time that the batcher will spend waiting for EigenDA to confirm a blob
	StatusQueryTimeout time.Duration

	// The amount of time to wait between status queries of a newly dispersed blob
	StatusQueryRetryInterval time.Duration
}
