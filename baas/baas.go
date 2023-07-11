// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package baas

import (
	mongodbImpl "github.com/byted-apaas/baas-sdk-go/mongodb/impl"
	"github.com/byted-apaas/baas-sdk-go/oss"
	"github.com/byted-apaas/baas-sdk-go/redis"
)

var (
	MongoDB = mongodbImpl.NewMongodb()
	Redis   = redis.NewRedis()
	Oss     = oss.NewOss()
)
