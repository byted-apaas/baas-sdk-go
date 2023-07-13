package exceptions

import (
	cExceptions "github.com/byted-apaas/server-common-go/exceptions"
)

type BaseError = cExceptions.BaseError

const (
	// ErrCode_RATE_LIMIT_ERROR限流错误码
	ErrCode_Rate_Limit = "k_cf_ec_200009"
)
