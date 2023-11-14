package oss

import (
	"bytes"
	cUtils "github.com/byted-apaas/server-common-go/utils"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/byted-apaas/baas-sdk-go/common/constants"
	http2 "github.com/byted-apaas/baas-sdk-go/http"
	cException "github.com/byted-apaas/server-common-go/exceptions"
)

func readFromURL(ctx context.Context, targetURL string) ([]byte, error) {

	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	rsp, err := http2.GetCommonHttpClient().Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("statusCode: %d", rsp.StatusCode)
	}
	return b, err
}

func uploadWithContent(ctx context.Context, name string, content []byte, option *Option) (*UploadResult, error) {
	if !cUtils.IsExternalFaaS() {
		return nil, cException.InvalidParamError("unsupport oss")
	}

	if len(content) > constants.MaxFileSize {
		return nil, cException.InvalidParamError("file too large, exceed %v", constants.MaxFileSize)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(constants.FileFieldName, name)
	if err != nil {
		return nil, err
	}

	if _, err := part.Write(content); err != nil {
		return nil, err
	}

	if option != nil {
		data, err := json.Marshal(option)
		if err != nil {
			return nil, err
		}

		if err := writer.WriteField(constants.FileOptionFieldName, string(data)); err != nil {
			return nil, err
		}

	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	out, err := http2.DoRequestFile(ctx, writer.FormDataContentType(), body)
	if err != nil {
		return nil, err
	}
	res := &fileUploadResult{}
	dest, err := base64.StdEncoding.DecodeString(string(out))
	if err != nil {
		return nil, fmt.Errorf("result decode err: %v", err)
	}
	if err := bson.Unmarshal(dest, &res); err != nil {
		return nil, err
	}

	if res.Data != nil && res.Data.uploadError == nil {
		return &UploadResult{
			URL: res.Data.URL,
		}, nil
	}

	return nil, res.Data.uploadError.error()
}
