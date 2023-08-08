package request

import (
	"context"
	"sync"

	"github.com/byted-apaas/baas-sdk-go/common/structs"
	"github.com/byted-apaas/baas-sdk-go/request/openapi"
	"github.com/byted-apaas/baas-sdk-go/tasks"
)

//go:generate mockery --name=IRequestOpenapi --structname=RequestOpenapi --filename=RequestOpenapi.go
type IRequestOpenapi interface {
	InvokeFunctionAsync(ctx context.Context, appCtx *structs.AppCtx, apiName string, params map[string]interface{}) (int64, error)
}

//go:generate mockery --name=IRequestOpenapi --structname=RequestOpenapi --filename=RequestOpenapi.go
type IRequestFaaSInfra interface {
	InvokeFunctionDistributed(ctx context.Context, appCtx *structs.AppCtx, dataset interface{}, handlerFunc string, progressCallbackFunc string, completedCallbackFunc string, options *tasks.Options) (int64, error)
}

var (
	reqHTTP     IRequestOpenapi
	reqHTTPOnce sync.Once
)

func GetInstance(ctx context.Context) IRequestOpenapi {
	return GetHTTPInstance()
}

func GetHTTPInstance() IRequestOpenapi {
	if reqHTTP == nil {
		reqHTTPOnce.Do(func() {
			reqHTTP = &openapi.RequestHttp{}
		})
	}
	return reqHTTP
}
