package errors

import "errors"

const (
	UnknownDATypeErrMsg   = "Rollup with unknown da type"
	DANotPreparedErrMsg   = "DA not prepared"
	WrongArgsNumberErrMsg = "Number of args is wrong"
	RollupFailedMsg       = "Rollup into DA failed"
	GetFromDAErrMsg       = "Get from DA failed"
	WrongArgTypeErrMsg    = "Arg with wrong type"
	NilPointerErrMsg      = "got nil pointer"
)

var (
	UnknownDATypeErr   = errors.New(UnknownDATypeErrMsg)
	DANotPreparedErr   = errors.New(DANotPreparedErrMsg)
	WrongArgsNumberErr = errors.New(WrongArgsNumberErrMsg)
	RollupFailedErr    = errors.New(RollupFailedMsg)
	GetFromDAErr       = errors.New(GetFromDAErrMsg)
	WrongArgTypeErr    = errors.New(WrongArgTypeErrMsg)
	NilPointerErr      = errors.New(NilPointerErrMsg)
)
