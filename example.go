// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package infra_go

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/byted-apaas/baas-sdk-go/infra"
	cond "github.com/byted-apaas/baas-sdk-go/mongodb/condition"
	cUtils "github.com/byted-apaas/server-common-go/utils"
)

func redisExample() {
	ctx := context.Background()
	redis := infra.Redis

	setRes, err := redis.Set(ctx, "key", "value", time.Second*10).Result()
	if err != nil {
		cUtils.PrintLog(err)
		return
	}
	cUtils.PrintLog(setRes)

	getRes, err := redis.Get(ctx, "key").Result()
	if err != nil {
		cUtils.PrintLog(err)
		return
	}
	cUtils.PrintLog(getRes)
}

func mongodbExample() {
	ctx := context.Background()

	mongodb := infra.MongoDB

	res, err := mongodb.Table("table_name").Create(ctx, cond.M{"name": "zhangsan", "age": 23})
	if err != nil {
		cUtils.PrintLog(err)
		return
	}
	cUtils.PrintLog(res)

	result := bson.M{}
	if err := mongodb.Table("table_name").Where(cond.M{"name": "zhangsan"}).Find(ctx, &result); err != nil {
		cUtils.PrintLog(err)
		return
	}
	cUtils.PrintLog(result)
}

func ossExample() {
	ctx := context.Background()
	oss := infra.Oss
	result, err := oss.UploadWithPath(ctx, "test.go", "./test.go", nil)
	if err != nil {
		cUtils.PrintLog(err)
		return
	}
	cUtils.PrintLog(result)
}
