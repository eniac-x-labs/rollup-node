module github.com/eniac-x-labs/rollup-node

go 1.21.1

toolchain go1.21.10

replace (
	//github.com/eniac-x-labs/anytrustDA => ./x/anytrust/anytrustDA
	github.com/near/rollup-data-availability => github.com/near/rollup-data-availability v0.2.4-0.20240327150603-b004cd803551
)

require (
	github.com/Layr-Labs/eigenda v0.6.1
	github.com/Layr-Labs/eigenda/api v0.6.1
	//github.com/eniac-x-labs/anytrustDA v0.0.0
	github.com/ethereum/go-ethereum v1.13.14
	github.com/fsnotify/fsnotify v1.6.0
	github.com/go-stack/stack v1.8.1
	github.com/holiman/uint256 v1.2.4
	github.com/near/rollup-data-availability v0.0.0-00010101000000-000000000000
	github.com/prometheus/client_golang v1.19.0
	github.com/rollkit/go-da v0.5.0
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.9.0
	github.com/urfave/cli/v2 v2.27.1
	golang.org/x/sync v0.6.0
	golang.org/x/term v0.16.0
	google.golang.org/grpc v1.62.1
)

require (
	github.com/DataDog/zstd v1.5.2 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/VictoriaMetrics/fastcache v1.12.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.10.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cockroachdb/errors v1.9.1 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v0.0.0-20230928194634-aa077af62593 // indirect
	github.com/cockroachdb/redact v1.1.3 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.12.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/crate-crypto/go-ipa v0.0.0-20231025140028-3c0104f4b233 // indirect
	github.com/crate-crypto/go-kzg-4844 v0.7.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set/v2 v2.1.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/ethereum/c-kzg-4844 v0.4.0 // indirect
	github.com/filecoin-project/go-jsonrpc v0.3.1 // indirect
	github.com/gballet/go-verkle v0.1.1-0.20231031103413-a67434b50f46 // indirect
	github.com/getsentry/sentry-go v0.18.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/ipfs/go-log/v2 v2.0.8 // indirect
	github.com/klauspost/compress v1.16.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/magiconair/properties v1.8.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/afero v1.1.2 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/supranational/blst v0.3.11 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	go.opencensus.io v0.22.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa // indirect
	golang.org/x/mod v0.16.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.15.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)
