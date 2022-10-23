package oss

import (
	"context"
	"io/ioutil"

	cException "github.com/byted-apaas/server-common-go/exceptions"
)

type IOss interface {
	UploadWithContent(ctx context.Context, name string, content []byte, option *Option) (*UploadResult, error)
	UploadWithURL(ctx context.Context, name string, targetUrl string, option *Option) (*UploadResult, error)
	UploadWithPath(ctx context.Context, name string, filePath string, option *Option) (*UploadResult, error)
}

type Oss struct{}

func NewOss() *Oss {
	return &Oss{}
}

func (f *Oss) UploadWithContent(ctx context.Context, name string, content []byte, option *Option) (*UploadResult, error) {
	return uploadWithContent(ctx, name, content, option)
}

func (f *Oss) UploadWithURL(ctx context.Context, name string, targetUrl string, option *Option) (*UploadResult, error) {
	data, err := readFromURL(ctx, targetUrl)
	if err != nil {
		return nil, cException.InvalidParamError("fetch data from targetUrl error: %v", err)
	}
	return uploadWithContent(ctx, name, data, option)
}

func (f *Oss) UploadWithPath(ctx context.Context, name string, filePath string, option *Option) (*UploadResult, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, cException.InvalidParamError("read data from filePath error: %v", err)
	}
	return uploadWithContent(ctx, name, data, option)
}
