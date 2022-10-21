// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package impl

import (
	"github.com/byted-apaas/baas-sdk-go/mongodb"
)

type Mongodb struct {
}

func NewMongodb() *Mongodb {
	return &Mongodb{}
}

func (m *Mongodb) Table(tableName string) mongodb.ITable {
	return NewTable(tableName)
}
