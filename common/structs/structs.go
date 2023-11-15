// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package structs

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecordOnlyId struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
}

type Option struct {
	Type string `json:"type,omitempty"` // http content type
	//Region string `json:"region"` // region of storage
}

type UploadResult struct {
	URL string `json:"url,omitempty"`
}

type UploadError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e UploadError) Error() error {
	if e.Code != 0 {
		if len(e.Message) == 0 {
			e.Message = "upload file fail"
		}
		return errors.New(fmt.Sprintf(`code: %d, message:"%s"`, e.Code, e.Message))
	}
	return nil
}

type FileUploadResult struct {
	Data *struct {
		URL string `json:"url,omitempty" bson:"url,omitempty"`
		*UploadError
	} `json:"data" bson:"data"`
}
