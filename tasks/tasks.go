package tasks

import "context"

//go:generate mockery --name=ITasks --structname=Tasks --filename=Tasks.go
type ITasks interface {
	// CreateAsyncTask 创建异步任务
	// @param apiName 函数的 apiName
	// @param params 参数
	// @return taskID
	CreateAsyncTask(ctx context.Context, apiName string, params map[string]interface{}) (int64, error)

	// CreateDistributedTask 创建分布式任务
	// @param dataset 待处理数据组成的数组
	// @param handlerFunc 用于处理数据集的全局函数的 API name
	// @param progressCallbackFunc 任务进度发生变化时回调的全局函数的 API name，可通过传入 "" 跳过此步骤
	// @param completedCallbackFunc 任务完成时回调的全局函数的 API name，可通过传入 "" 跳过此步骤
	// @param options 高级配置参数，可选
	CreateDistributedTask(ctx context.Context, dataset interface{}, handlerFunc string, progressCallbackFunc string,
		completedCallbackFunc string, options *Options) (int64, error)
}

type Options struct {
	// 并发数量，默认值为 5，最大可设置值为 10
	Concurrency int64 `json:"concurrency"`
	// 单个子任务的最大数据量，默认值为 5，最大可设置值为 100
	MaxSliceNumber int64 `json:"maxSliceNumber"`
	// 触发进度回调函数的步长，以百分比为单位，每当发生大于等于步长的进度变化时，便触发进度回调函数。默认为每当有进度变化时，便触发进度回调函数（进度变化精度为 1%）
	ProgressCallbackStep int64 `json:"progressCallbackStep"`
}

func NewOptions(concurrency, maxSliceNumber, progressCallbackStep int64) *Options {
	return &Options{concurrency, maxSliceNumber, progressCallbackStep}
}
