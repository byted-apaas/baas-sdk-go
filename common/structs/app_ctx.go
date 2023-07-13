package structs

import (
	"github.com/byted-apaas/baas-sdk-go/common/constants"
	cHttp "github.com/byted-apaas/server-common-go/http"
)

type AppCtx struct {
	Mode       AppMode
	Env        constants.PlatformEnvType
	Credential *cHttp.AppCredential
}

func (a *AppCtx) IsOpenSDK() bool {
	if a != nil && a.Mode == AppModeOpenSDK {
		return true
	}
	return false
}

func (a *AppCtx) GetEnv() constants.PlatformEnvType {
	if a != nil && a.Env != 0 {
		return a.Env
	}
	return constants.PlatformEnvOnline
}

type AppMode int

const (
	AppModeOpenSDK AppMode = iota + 1
	AppModeFaaSSDK
)
