package oss

import (
	"github.com/byted-apaas/baas-sdk-go/common/structs"
	"github.com/byted-apaas/baas-sdk-go/request/faasinfra"
	"context"
	"io/ioutil"

	cException "github.com/byted-apaas/server-common-go/exceptions"
)

type IOss interface {
	UploadWithContent(ctx context.Context, name string, content []byte, option *structs.Option) (*structs.UploadResult, error)
	UploadWithURL(ctx context.Context, name string, targetUrl string, option *structs.Option) (*structs.UploadResult, error)
	UploadWithPath(ctx context.Context, name string, filePath string, option *structs.Option) (*structs.UploadResult, error)
}

type Oss struct{}

func NewOss() *Oss {
	return &Oss{}
}

func (f *Oss) UploadWithContent(ctx context.Context, name string, content []byte, option *structs.Option) (*structs.UploadResult, error) {
	return faasinfra.UploadWithContent(ctx, name, content, option)
}

func (f *Oss) UploadWithURL(ctx context.Context, name string, targetUrl string, option *structs.Option) (*structs.UploadResult, error) {
	data, err := faasinfra.ReadFromURL(ctx, targetUrl)
	if err != nil {
		return nil, cException.InvalidParamError("fetch data from targetUrl error: %v", err)
	}
	return faasinfra.UploadWithContent(ctx, name, data, option)
}

func (f *Oss) UploadWithPath(ctx context.Context, name string, filePath string, option *structs.Option) (*structs.UploadResult, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, cException.InvalidParamError("read data from filePath error: %v", err)
	}
	return faasinfra.UploadWithContent(ctx, name, data, option)
}
