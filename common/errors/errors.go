package errors

import "errors"

const (
	UnknownDATypeErrMsg = "rollup with unknown da type"
)

var (
	UnknownDATypeErr = errors.New("rollup with unknown da type")
)
