package impl

import (
	"context"

	"github.com/byted-apaas/baas-sdk-go/common/structs"
	"github.com/byted-apaas/baas-sdk-go/request"
	"github.com/byted-apaas/baas-sdk-go/request/faasinfra"
	"github.com/byted-apaas/baas-sdk-go/tasks"
	cExceptions "github.com/byted-apaas/server-common-go/exceptions"
	cUtils "github.com/byted-apaas/server-common-go/utils"
)

type Tasks struct {
	appCtx *structs.AppCtx
}

func NewTasks(s *structs.AppCtx) tasks.ITasks {
	return &Tasks{appCtx: s}
}

func (t *Tasks) CreateAsyncTask(ctx context.Context, funcAPIName string, params map[string]interface{}) (int64, error) {
	return request.GetInstance(ctx).InvokeFunctionAsync(ctx, t.appCtx, funcAPIName, params)
}

// CreateDistributedTask
/**
 * 创建分布式任务
 * @param dataset 待处理数据组成的数组
 * @param handlerFunc 用于处理数据集的全局函数的 API name
 * @param progressCallbackFunc 任务进度发生变化时回调的全局函数的 API name，可通过传入 "" 跳过此步骤
 * @param completedCallbackFunc 任务完成时回调的全局函数的 API name，可通过传入 "" 跳过此步骤
 * @param options 高级配置参数，可选
 */
func (t *Tasks) CreateDistributedTask(ctx context.Context, dataset interface{}, handlerFunc string, progressCallbackFunc string, completedCallbackFunc string, options *tasks.Options) (int64, error) {
	if len(cUtils.GetDistributedMaskFromCtx(ctx)) > 0 {
		return 0, cExceptions.InvalidParamError("Does not support creating distributed tasks in handler function and progress callback function")
	}

	if handlerFunc == "" {
		return 0, cExceptions.InvalidParamError("The handlerFunc can not be empty")
	}

	if options == nil {
		options = tasks.NewOptions(5, 5, 1)
	}
	if options.Concurrency == 0 {
		options.Concurrency = 5
	}
	if options.MaxSliceNumber == 0 {
		options.MaxSliceNumber = 5
	}
	if options.ProgressCallbackStep == 0 {
		options.ProgressCallbackStep = 1
	}
	if options.ProgressCallbackStep <= 0 || options.ProgressCallbackStep > 100 {
		return 0, cExceptions.InvalidParamError("parameter option.progressCallbackStep is not between 0 and 100")
	}
	return faasinfra.GetInstance().InvokeFunctionDistributed(ctx, t.appCtx, dataset, handlerFunc, progressCallbackFunc, completedCallbackFunc, options)
}
