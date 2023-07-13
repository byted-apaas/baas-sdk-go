package faasinfra

import (
	cUtils "github.com/byted-apaas/server-common-go/utils"
)

const (
	PathInvokeFunction            = "/cloudfunction/v1/namespaces/:namespace/function/invokeSync"
	PathInvokeFunctionAsync       = "/faasAsyncTask/v1/namespaces/:namespace/asyncTask/CreateAsyncTask"
	PathInvokeFunctionDistributed = "/distributedTask/v1/namespaces/:namespace/create"
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
