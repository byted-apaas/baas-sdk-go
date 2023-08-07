package common

import (
	"context"

	cExceptions "github.com/byted-apaas/server-common-go/exceptions"
	cHttp "github.com/byted-apaas/server-common-go/http"
	cUtils "github.com/byted-apaas/server-common-go/utils"
)

func BuildInvokeParamsStr(ctx context.Context, funcAPIName string, params interface{}, needPermission bool) (map[string]interface{}, error) {
	sysParams, bizParams, err := BuildInvokeParamAndContext(ctx, params, funcAPIName, needPermission)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"apiAlias":    funcAPIName,
		"params":      bizParams,
		"context":     sysParams,
		"triggerType": "workflow",
	}, nil
}

func BuildInvokeParamAndContext(ctx context.Context, params interface{}, funcAPIName string, needPermission bool) (string, string, error) {
	sysParams, _ := cUtils.JsonMarshalBytes(BuildInvokeSysParams(ctx, params, funcAPIName, needPermission))

	bizParams, err := cUtils.JsonMarshalBytes(params)
	if err != nil {
		return "", "", cExceptions.InvalidParamError("Marshal params failed, err: %+v", err)
	}

	return string(sysParams), string(bizParams), nil
}

func BuildInvokeSysParams(ctx context.Context, params interface{}, funcAPIName string, needPermission bool) map[string]interface{} {
	v := map[string]interface{}{
		"triggertaskid":             cUtils.GetTriggerTaskIDFromCtx(ctx),
		"x-kunlun-distributed-mask": cUtils.GetDistributedMaskFromCtx(ctx),
		"x-kunlun-loop-masks":       cUtils.GetLoopMaskFromCtx(ctx),
	}

	if needPermission {
		v["permission"] = cHttp.CalcParamsNeedPermission(ctx, funcAPIName, "input", params)
	}
	return v
}
