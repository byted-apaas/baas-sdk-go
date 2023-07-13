package openapi

import (
	"github.com/tidwall/gjson"

	"github.com/byted-apaas/baas-sdk-go/common/exceptions"
	cExceptions "github.com/byted-apaas/server-common-go/exceptions"
	cUtils "github.com/byted-apaas/server-common-go/utils"
)

const (
	SCFileDownload = ""
	SCSuccess      = "0"

	ECInternalError  = "k_ec_000001"
	ECNoTenantID     = "k_ec_000002"
	ECNoUserID       = "k_ec_000003"
	ECUnknownError   = "k_ec_000004"
	ECOpUnknownError = "k_op_ec_00001"
	ECSystemBusy     = "k_op_ec_20001"
	ECSystemError    = "k_op_ec_20002"
	ECRateLimitError = "k_op_ec_20003"
	ECTokenExpire    = "k_ident_013000"
	ECIllegalToken   = "k_ident_013001"
	ECMissingToken   = "k_op_ec_10205"
)

func errorWrapper(body []byte, extra map[string]interface{}, err error) ([]byte, error) {
	if err != nil {
		return nil, cExceptions.ErrWrap(err)
	}

	code := gjson.GetBytes(body, "code").String()
	msg := gjson.GetBytes(body, "msg").String()
	switch code {
	case SCFileDownload:
		return body, nil
	case SCSuccess:
		data := gjson.GetBytes(body, "data")
		if data.Type == gjson.String {
			return []byte(data.Str), nil
		}
		return []byte(data.Raw), nil
	case ECRateLimitError:
		return nil, cExceptions.NewErrWithCode(exceptions.ErrCode_Rate_Limit, "%v (logID: %v)", msg, cUtils.GetLogIDFromExtra(extra))
	case ECInternalError, ECNoTenantID, ECNoUserID, ECUnknownError,
		ECOpUnknownError, ECSystemBusy, ECSystemError,
		ECTokenExpire, ECIllegalToken, ECMissingToken:
		return nil, cExceptions.InternalError("%v ([%v] %v)", msg, code, cUtils.GetLogIDFromExtra(extra))
	default:
		return nil, cExceptions.NewErrWithCode(code, "%v (logID: %v)", msg, cUtils.GetLogIDFromExtra(extra))
	}
}
