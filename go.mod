module github.com/eniac-x-labs/rollup-node

go 1.22.0

replace (
	github.com/eniac-x-labs/anytrustDA => ./x/anytrust/anytrustDA
)
require (
	github.com/eniac-x-labs/anytrustDA v0.0.0
)