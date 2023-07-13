package utils

import (
	"context"
	"strconv"

	"github.com/byted-apaas/baas-sdk-go/common/structs"
	cConstants "github.com/byted-apaas/server-common-go/constants"
	cExceptions "github.com/byted-apaas/server-common-go/exceptions"
	cHttp "github.com/byted-apaas/server-common-go/http"
	cUtils "github.com/byted-apaas/server-common-go/utils"
	"github.com/bytedance/gopkg/cloud/metainfo"
)

func SetCtx(ctx context.Context, appCtx *structs.AppCtx, method string) context.Context {
	ctx = SetAppConfToCtx(ctx, appCtx)
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = cUtils.SetApiTimeoutMethodToCtx(ctx, method)
	return ctx
}

func SetAppConfToCtx(ctx context.Context, appCtx *structs.AppCtx) context.Context {
	if appCtx == nil || appCtx.Mode != structs.AppModeOpenSDK {
		return ctx
	}

	if ctx == nil {
		ctx = context.Background()
	}

	targetEnv := appCtx.GetEnv()
	conf, _ := cConstants.EnvConfMap[targetEnv.String()]

	ctx = cHttp.SetCredentialToCtx(ctx, appCtx.Credential)
	ctx = context.WithValue(ctx, cConstants.CtxKeyInnerAPIPSM, conf.InnerAPIPSM)
	ctx = context.WithValue(ctx, cConstants.CtxKeyOpenapiDomain, conf.OpenAPIDomain)
	ctx = context.WithValue(ctx, cConstants.CtxKeyFaaSInfraDomain, conf.FaaSInfraDomain)
	ctx = context.WithValue(ctx, cConstants.CtxKeyAGWDomain, conf.InnerAPIDomain)
	boe := conf.BOE
	if len(boe) > 0 {
		ctx = cUtils.SetEnvBoeToCtx(ctx, boe)
	}

	return ctx
}

func GetNamespace(ctx context.Context, appCtx *structs.AppCtx) (string, error) {
	if appCtx != nil && appCtx.Credential != nil && appCtx.Mode == structs.AppModeOpenSDK {
		tenant, err := appCtx.Credential.GetTenantInfo(ctx)
		if err != nil {
			return "", cExceptions.InternalError("GetTenantInfo failed, err: %v", err)
		}
		return tenant.Namespace, nil
	}
	return cUtils.GetNamespace(), nil
}

func GetTenantName(ctx context.Context, appCtx *structs.AppCtx) (string, error) {
	if appCtx != nil && appCtx.Credential != nil && appCtx.Mode == structs.AppModeOpenSDK {
		tenant, err := appCtx.Credential.GetTenantInfo(ctx)
		if err != nil {
			return "", cExceptions.InternalError("GetTenantInfo failed, err: %v", err)
		}
		return tenant.Name, nil
	}
	return cUtils.GetTenantName(), nil
}

func SetUserMetaInfoToContext(ctx context.Context, appCtx *structs.AppCtx) context.Context {
	if appCtx.IsOpenSDK() {
		return ctx
	}
	ctx = metainfo.WithValue(ctx, cConstants.HttpHeaderKeyUser, strconv.FormatInt(cUtils.GetUserIDFromCtx(ctx), 10))
	return ctx
}
