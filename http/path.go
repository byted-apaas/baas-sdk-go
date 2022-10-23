// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package http

import (
	"strings"

	cConstants "github.com/byted-apaas/server-common-go/constants"
	cUtils "github.com/byted-apaas/server-common-go/utils"
)

const (
	PathFaaSInfraPathMongodb = "/resource/v3/namespaces/:namespace/db"
	PathFaaSInfraPathRedis   = "/resource/v2/namespaces/:namespace/cache"
	FaaSInfraPathFile        = "/resource/v2/namespaces/:namespace/file"
)

func GetFaaSInfraPathMongodb() string {
	return strings.Replace(PathFaaSInfraPathMongodb, cConstants.ReplaceNamespace, cUtils.GetNamespace(), 1)
}

func GetFaaSInfraPathRedis() string {
	return strings.Replace(PathFaaSInfraPathRedis, cConstants.ReplaceNamespace, cUtils.GetNamespace(), 1)
}

func GetFaaSInfraPathFile() string {
	return strings.Replace(FaaSInfraPathFile, cConstants.ReplaceNamespace, cUtils.GetNamespace(), 1)
}
