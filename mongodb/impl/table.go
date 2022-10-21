// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package impl

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/byted-apaas/baas-sdk-go/http"
	"github.com/byted-apaas/baas-sdk-go/mongodb"
	"github.com/byted-apaas/baas-sdk-go/mongodb/structs"
	cExceptions "github.com/byted-apaas/server-common-go/exceptions"
)

type Table struct {
	*MongodbParam
}

func NewTable(tableName string) *Table {
	t := &Table{MongodbParam: NewMongodbParam(tableName)}
	if len(tableName) == 0 {
		t.Err = cExceptions.InvalidParamError("tableName is empty")
	}
	return t
}

func (t *Table) Create(ctx context.Context, record interface{}) (*structs.RecordOnlyId, error) {
	if t.Err != nil {
		return nil, t.Err
	}

	t.SetOp(OpType_Insert)
	t.SetDocs([]interface{}{record})
	return http.Create(ctx, t.MongodbParam)
}

func (t *Table) BatchCreate(ctx context.Context, records interface{}) ([]primitive.ObjectID, error) {
	if t.Err != nil {
		return nil, t.Err
	}

	t.SetOp(OpType_Insert)
	t.SetDocs(records)
	return http.BatchCreate(ctx, t.MongodbParam)
}

func (t *Table) Where(condition interface{}, args ...interface{}) mongodb.IQuery {
	return NewQuery(t.MongodbParam.TableName).Where(condition, args)
}

func (q *Table) GroupBy(field interface{}, alias ...interface{}) mongodb.IAggQuery {
	return NewAggQuery(q.TableName).GroupBy(field, alias...)
}
