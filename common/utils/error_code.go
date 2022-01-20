/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : error_code.go
*   coder: zemanzeng
*   date : 2021-05-26 15:52:53
*   desc : 错误码封装
*
================================================================*/

package utils

import (
	"fmt"
)

const ECInternalError = 99999

var errorCodes = make(map[int32]string)

// 服务初始化时register
func RegisterErrorCode(code int32, msg string) {
	errorCodes[code] = msg
}

func GetCodeMsg(code int32) string {
	if msg, exist := errorCodes[code]; exist {
		return msg
	}
	return ""
}

// ErrorCode 业务error封装
type ErrorCode struct {
	Code   int32
	Msg    string
	ExtMsg string
}

func (ec *ErrorCode) Error() string {
	if len(ec.Msg) == 0 {
		ec.Msg = GetCodeMsg(ec.Code)
	}

	s := fmt.Sprintf("code:%v", ec.Code)
	if len(ec.Msg) > 0 {
		s = s + " msg:" + ec.Msg
	}
	if len(ec.ExtMsg) > 0 {
		s = s + " ext:" + ec.ExtMsg
	}
	return s
}

func ECError(code int32, msgs ...string) error {
	var msg, extMsg string
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	if len(msgs) > 1 {
		extMsg = msgs[1]
	}
	return &ErrorCode{Code: code, Msg: msg, ExtMsg: extMsg}
}

func ECErrorExt(code int32, extMsg string) error {
	return &ErrorCode{Code: code, Msg: GetCodeMsg(code), ExtMsg: extMsg}
}

func GetErrorCode(err error) int32 {
	ec, ok := err.(*ErrorCode)
	if ok {
		return ec.Code
	}
	return ECInternalError
}

func GetErrorMsg(err error) string {
	if err == nil {
		return ""
	}
	ec, ok := err.(*ErrorCode)
	if ok {
		if len(ec.Msg) > 0 {
			return ec.Msg
		}
		return GetCodeMsg(ec.Code)
	}
	return err.Error()
}

func GetErrorExtMsg(err error) string {
	if err == nil {
		return ""
	}
	if ec, ok := err.(*ErrorCode); ok {
		return ec.ExtMsg
	}
	return ""
}
