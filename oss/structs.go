package oss

import (
	"errors"
	"fmt"
)

type Option struct {
	Type string `json:"type,omitempty"` // http content type
	//Region string `json:"region"` // region of storage
}

type UploadResult struct {
	URL string `json:"url,omitempty"`
}

type uploadError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e uploadError) error() error {
	if e.Code != 0 {
		if len(e.Message) == 0 {
			e.Message = "upload file fail"
		}
		return errors.New(fmt.Sprintf(`code: %d, message:"%s"`, e.Code, e.Message))
	}
	return nil
}

type fileUploadResult struct {
	Data *struct {
		URL string `json:"url,omitempty" bson:"url,omitempty"`
		*uploadError
	} `json:"data" bson:"data"`
}
