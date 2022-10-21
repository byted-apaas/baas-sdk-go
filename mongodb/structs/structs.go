// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package structs

import "go.mongodb.org/mongo-driver/bson/primitive"

type RecordOnlyId struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
}
