package openapi

import (
	cUtils "github.com/byted-apaas/server-common-go/utils"
)

const (
	PathInvokeFunctionAsync = "/faasAsyncTask/v1/namespaces/:namespace/asyncTask/CreateAsyncTask"
)

func GetPathInvokeFunctionAsync(namespace string) string {
	return cUtils.NewPathReplace(PathInvokeFunctionAsync).Namespace(namespace).Path()
}
