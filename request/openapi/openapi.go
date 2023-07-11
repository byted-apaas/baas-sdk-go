package openapi

import (
	"context"
	"strconv"

	"github.com/tidwall/gjson"

	"github.com/byted-apaas/baas-sdk-go/common/structs"
	"github.com/byted-apaas/baas-sdk-go/common/utils"
	reqCommon "github.com/byted-apaas/baas-sdk-go/request/common"
	cConstants "github.com/byted-apaas/server-common-go/constants"
	cHttp "github.com/byted-apaas/server-common-go/http"
	cUtils "github.com/byted-apaas/server-common-go/utils"
)

type RequestHttp struct{}

func (r *RequestHttp) InvokeFunctionAsync(ctx context.Context, appCtx *structs.AppCtx, apiName string, params map[string]interface{}) (int64, error) {
	body, err := reqCommon.BuildInvokeParamsStr(ctx, apiName, params)
	if err != nil {
		return 0, err
	}

	namespace, err := utils.GetNamespace(ctx, appCtx)
	if err != nil {
		return 0, err
	}
	tenantName, err := utils.GetTenantName(ctx, appCtx)
	if err != nil {
		return 0, err
	}
	headers := map[string][]string{
		cConstants.HttpHeaderKeyTenant: {tenantName},
		cConstants.HttpHeaderKeyUser:   {strconv.FormatInt(cUtils.GetUserIDFromCtx(ctx), 10)},
	}

	data, err := errorWrapper(getOpenapiClient().PostJson(utils.SetAppConfToCtx(ctx, appCtx), GetPathInvokeFunctionAsync(namespace), headers, body, cHttp.AppTokenMiddleware))
	if err != nil {
		return 0, err
	}

	return gjson.GetBytes(data, "task_id").Int(), nil
}
