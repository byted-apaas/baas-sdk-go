package faasinfra

import (
	"context"
	"reflect"
	"strconv"
	"sync"

	"github.com/tidwall/gjson"

	"github.com/byted-apaas/baas-sdk-go/common/structs"
	"github.com/byted-apaas/baas-sdk-go/common/utils"
	"github.com/byted-apaas/baas-sdk-go/request"
	"github.com/byted-apaas/baas-sdk-go/request/common"
	"github.com/byted-apaas/baas-sdk-go/tasks"
	cConstants "github.com/byted-apaas/server-common-go/constants"
	cExceptions "github.com/byted-apaas/server-common-go/exceptions"
	cHttp "github.com/byted-apaas/server-common-go/http"
	cUtils "github.com/byted-apaas/server-common-go/utils"
)

type requestFaaSInfra struct{}

var (
	reqFaaSInfra     request.IRequestFaaSInfra
	reqFaaSInfraOnce sync.Once
)

func GetInstance() request.IRequestFaaSInfra {
	if reqFaaSInfra == nil {
		reqFaaSInfraOnce.Do(func() {
			reqFaaSInfra = &requestFaaSInfra{}
		})
	}
	return reqFaaSInfra
}

func (r *requestFaaSInfra) InvokeFunction(ctx context.Context, appCtx *structs.AppCtx, apiName string, params map[string]interface{}, result interface{}) error {
	body, err := common.BuildInvokeParamsStr(ctx, apiName, params)
	if err != nil {
		return err
	}

	namespace, err := utils.GetNamespace(ctx, appCtx)
	if err != nil {
		return err
	}
	tenantName, err := utils.GetTenantName(ctx, appCtx)
	if err != nil {
		return err
	}
	headers := map[string][]string{
		cConstants.HttpHeaderKeyTenant: {tenantName},
		cConstants.HttpHeaderKeyUser:   {strconv.FormatInt(cUtils.GetUserIDFromCtx(ctx), 10)},
	}

	data, err := errorWrapper(getFaaSInfraClient(ctx).PostJson(utils.SetAppConfToCtx(ctx, appCtx), GetPathInvokeFunction(namespace), headers, body, cHttp.AppTokenMiddleware))
	if err != nil {
		return err
	}

	resultRaw := gjson.GetBytes(data, "data").Raw
	err = cUtils.JsonUnmarshalBytes([]byte(resultRaw), result)
	if err != nil {
		return cExceptions.InternalError("InvokeFunction failed, err: %v", err)
	}
	return nil
}

func (r *requestFaaSInfra) InvokeFunctionAsync(ctx context.Context, appCtx *structs.AppCtx, apiName string, params map[string]interface{}) (int64, error) {
	body, err := common.BuildInvokeParamsStr(ctx, apiName, params)
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

	data, err := errorWrapper(getFaaSInfraClient(ctx).PostJson(utils.SetAppConfToCtx(ctx, appCtx), GetPathInvokeFunctionAsync(namespace), headers, body, cHttp.AppTokenMiddleware))
	if err != nil {
		return 0, err
	}

	return gjson.GetBytes(data, "task_id").Int(), nil
}

func (r *requestFaaSInfra) InvokeFunctionDistributed(ctx context.Context, appCtx *structs.AppCtx, dataset interface{}, handlerFunc string, progressCallbackFunc string, completedCallbackFunc string, options *tasks.Options) (int64, error) {
	v := reflect.ValueOf(dataset)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice {
		return 0, cExceptions.InvalidParamError("The type of dataset should be slice, but %s", v.Kind())
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

	lookMask := cUtils.GetLoopMaskFromCtx(ctx)
	if (handlerFunc != "" && cUtils.StrInStrs(lookMask, handlerFunc)) ||
		(progressCallbackFunc != "" && cUtils.StrInStrs(lookMask, progressCallbackFunc)) ||
		(completedCallbackFunc != "" && cUtils.StrInStrs(lookMask, completedCallbackFunc)) {
		return 0, cExceptions.InvalidParamError("Distributed task execution forms a loop.")
	}

	body := map[string]interface{}{
		"domainName":            tenantName,
		"namespace":             namespace,
		"userId":                cUtils.GetUserIDFromCtx(ctx),
		"dataset":               []interface{}{dataset},
		"handlerFunc":           handlerFunc,
		"progressCallbackFunc":  progressCallbackFunc,
		"completedCallbackFunc": completedCallbackFunc,
		"options":               options,
		"x-kunlun-loop-masks":   lookMask,
	}

	data, err := errorWrapper(getFaaSInfraClient(ctx).PostJson(utils.SetAppConfToCtx(ctx, appCtx), GetPathInvokeFunctionDistributed(namespace), headers, body, cHttp.AppTokenMiddleware))
	if err != nil {
		return 0, err
	}

	var taskID int64
	err = cUtils.JsonUnmarshalBytes(data, &taskID)
	if err != nil {
		return 0, cExceptions.InternalError("Unmarshal result failed, err: %v", err)
	}

	return taskID, nil
}
