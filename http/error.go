// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package http

import (
	cExceptions "github.com/byted-apaas/server-common-go/exceptions"
	cUtils "github.com/byted-apaas/server-common-go/utils"
	"github.com/tidwall/gjson"
)

const (
	FaaSInfraSuccessCodeSuccess     = "0"
	FaaSInfraFailCodeInternalError  = "k_ec_000001"
	FaaSInfraFailCodeTokenExpire    = "k_ident_013000"
	FaaSInfraFailCodeIllegalToken   = "k_ident_013001"
	FaaSInfraFailCodeMissingToken   = "k_fs_ec_100001"
	FaaSInfraFailCodeRateLimitError = "k_fs_ec_000004"
)

func HasError(errCode string) bool {
	return errCode != FaaSInfraSuccessCodeSuccess
}

func IsSysError(errCode string) bool {
	return errCode == FaaSInfraFailCodeInternalError ||
		errCode == FaaSInfraFailCodeTokenExpire ||
		errCode == FaaSInfraFailCodeIllegalToken ||
		errCode == FaaSInfraFailCodeMissingToken ||
		errCode == FaaSInfraFailCodeRateLimitError
}

func errorWrapper(body []byte, extra map[string]interface{}, err error) ([]byte, error) {
	if err != nil {
		return nil, cExceptions.ErrWrap(err)
	}
	code := gjson.GetBytes(body, "code").String()
	msg := gjson.GetBytes(body, "msg").String()
	if !HasError(code) {
		data := gjson.GetBytes(body, "data")
		if data.Type == gjson.String {
			return []byte(data.Str), nil
		}
		return []byte(data.Raw), nil
	} else if IsSysError(code) {
		return nil, cExceptions.InternalError("%v ([%v] %v)", msg, code, cUtils.GetLogIDFromExtra(extra))
	} else {
		return nil, cExceptions.InvalidParamError("%v ([%v] %v)", msg, code, cUtils.GetLogIDFromExtra(extra))
	}
}
