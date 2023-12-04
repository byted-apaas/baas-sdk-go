package faasinfra

import (
	cConstants "github.com/byted-apaas/server-common-go/constants"
	cUtils "github.com/byted-apaas/server-common-go/utils"
	"strings"
)

const (
	PathInvokeFunction            = "/cloudfunction/v1/namespaces/:namespace/function/invokeSync"
	PathInvokeFunctionAsync       = "/faasAsyncTask/v1/namespaces/:namespace/asyncTask/CreateAsyncTask"
	PathInvokeFunctionDistributed = "/distributedTask/v1/namespaces/:namespace/create"
	PathFaaSInfraPathMongodb      = "/resource/v3/namespaces/:namespace/db"
	PathFaaSInfraPathRedis        = "/resource/v2/namespaces/:namespace/cache"
	FaaSInfraPathFile             = "/resource/v2/namespaces/:namespace/file"
)

func GetPathInvokeFunction(namespace string) string {
	return cUtils.NewPathReplace(PathInvokeFunction).Namespace(namespace).Path()
}

func GetPathInvokeFunctionAsync(namespace string) string {
	return cUtils.NewPathReplace(PathInvokeFunctionAsync).Namespace(namespace).Path()
}

func GetPathInvokeFunctionDistributed(namespace string) string {
	return cUtils.NewPathReplace(PathInvokeFunctionDistributed).Namespace(namespace).Path()
}

func GetFaaSInfraPathMongodb() string {
	return strings.Replace(PathFaaSInfraPathMongodb, cConstants.ReplaceNamespace, cUtils.GetNamespace(), 1)
}

func GetFaaSInfraPathRedis() string {
	return strings.Replace(PathFaaSInfraPathRedis, cConstants.ReplaceNamespace, cUtils.GetNamespace(), 1)
}

func GetFaaSInfraPathFile() string {
	return strings.Replace(FaaSInfraPathFile, cConstants.ReplaceNamespace, cUtils.GetNamespace(), 1)
}
