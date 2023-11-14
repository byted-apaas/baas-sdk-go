package innerapi

type RequestRpc struct{}

//func (r *RequestRpc) pre(ctx context.Context, appCtx *structs.AppCtx, method string) (context.Context, context.CancelFunc, string, error) {
//	var err error
//	ctx, err = cHttp.RebuildRpcCtx(utils.SetCtx(ctx, appCtx, method))
//	if err != nil {
//		return nil, nil, "", err
//	}
//
//	namespace, err := utils.GetNamespace(ctx, appCtx)
//	if err != nil {
//		return nil, nil, "", err
//	}
//
//	if namespace == "" {
//		return nil, nil, "", cExceptions.InternalError("namespace is empty")
//	}
//
//	var cancel context.CancelFunc
//	ctx, cancel = cHttp.GetTimeoutCtx(ctx)
//	return ctx, cancel, namespace, nil
//}
//
//func (r *RequestRpc) post(ctx context.Context, err error, baseResp *base.BaseResp, baseReq *base.Base) error {
//	var logid string
//	if baseReq != nil {
//		logid = baseReq.LogID
//	}
//
//	if err != nil {
//		return cExceptions.InternalError("Call InnerAPI failed: %+v, logid: %s", err, logid)
//	}
//
//	if baseResp == nil {
//		return cExceptions.InternalError("Call InnerAPI resp is empty, logid: %s", logid)
//	}
//
//	if baseResp.KStatusCode != "" {
//		msg := baseResp.KStatusMessage
//		if baseResp.StatusMessage != "" {
//			msg = baseResp.StatusMessage
//		}
//		return cExceptions.NewErrWithCodeV2(baseResp.KStatusCode, msg, logid)
//	}
//	return nil
//}
//
//func (r *RequestRpc) InvokeFunctionAsync(ctx context.Context, appCtx *structs.AppCtx, apiName string, params map[string]interface{}) (int64, error) {
//	sysParams, bizParams, err := reqCommon.BuildInvokeParamAndContext(ctx, params, apiName, appCtx == nil || appCtx.Credential == nil || appCtx.Mode != structs.AppModeOpenSDK)
//	if err != nil {
//		return 0, err
//	}
//
//	req := cloudfunction.NewCreateAsyncTaskRequest()
//
//	ctx, cancel, _, err := r.pre(ctx, appCtx, cConstants.CreateAsyncTask)
//	if err != nil {
//		return 0, err
//	}
//	defer cancel()
//
//	ctx = utils.SetUserMetaInfoToContext(ctx, appCtx)
//	namespace, err := utils.GetNamespace(ctx, appCtx)
//	if err != nil {
//		return 0, err
//	}
//	req.Namespace = namespace
//	req.APIAlias = apiName
//	req.Context = sysParams
//	req.TriggerType = "workflow"
//	req.Params = cUtils.StringPtr(bizParams)
//
//	cli, err := cHttp.GetInnerAPICli(ctx)
//	if err != nil {
//		return 0, err
//	}
//
//	resp, err := cli.CreateAsyncTask(ctx, req)
//
//	var baseResp *base.BaseResp
//	if resp != nil {
//		baseResp = resp.BaseResp
//	}
//	if err = r.post(ctx, err, baseResp, req.Base); err != nil {
//		return 0, err
//	}
//
//	return resp.TaskID, nil
//}
