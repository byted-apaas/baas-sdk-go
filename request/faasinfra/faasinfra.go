package faasinfra

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"sync"

	"github.com/byted-apaas/baas-sdk-go/common/constants"
	"github.com/byted-apaas/baas-sdk-go/common/structs"
	"github.com/byted-apaas/baas-sdk-go/common/utils"
	"github.com/byted-apaas/baas-sdk-go/mongodb/structs/inner"
	"github.com/byted-apaas/baas-sdk-go/request"
	"github.com/byted-apaas/baas-sdk-go/tasks"
	cConstants "github.com/byted-apaas/server-common-go/constants"
	cExceptions "github.com/byted-apaas/server-common-go/exceptions"
	cHttp "github.com/byted-apaas/server-common-go/http"
	cUtils "github.com/byted-apaas/server-common-go/utils"
)

type requestFaaSInfra struct{}

var (
	reqFaaSInfra     request.IRequestFaaSInfra
	reqFaaSInfraOnce sync.Once
)

func GetInstance() request.IRequestFaaSInfra {
	if reqFaaSInfra == nil {
		reqFaaSInfraOnce.Do(func() {
			reqFaaSInfra = &requestFaaSInfra{}
		})
	}
	return reqFaaSInfra
}

func (r *requestFaaSInfra) InvokeFunctionDistributed(ctx context.Context, appCtx *structs.AppCtx, dataset interface{}, handlerFunc string, progressCallbackFunc string, completedCallbackFunc string, options *tasks.Options) (int64, error) {
	v := reflect.ValueOf(dataset)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice {
		return 0, cExceptions.InvalidParamError("The type of dataset should be slice, but %s", v.Kind())
	}

	namespace, err := utils.GetNamespace(ctx, appCtx)
	if err != nil {
		return 0, err
	}
	tenantName, err := utils.GetTenantName(ctx, appCtx)
	if err != nil {
		return 0, err
	}
	headers := map[string][]string{
		cConstants.HttpHeaderKeyTenant: {tenantName},
		cConstants.HttpHeaderKeyUser:   {strconv.FormatInt(cUtils.GetUserIDFromCtx(ctx), 10)},
	}

	lookMask := cUtils.GetLoopMaskFromCtx(ctx)
	if (handlerFunc != "" && cUtils.StrInStrs(lookMask, handlerFunc)) ||
		(progressCallbackFunc != "" && cUtils.StrInStrs(lookMask, progressCallbackFunc)) ||
		(completedCallbackFunc != "" && cUtils.StrInStrs(lookMask, completedCallbackFunc)) {
		return 0, cExceptions.InvalidParamError("Distributed task execution forms a loop.")
	}

	body := map[string]interface{}{
		"domainName":            tenantName,
		"namespace":             namespace,
		"userId":                cUtils.GetUserIDFromCtx(ctx),
		"dataset":               []interface{}{dataset},
		"handlerFunc":           handlerFunc,
		"progressCallbackFunc":  progressCallbackFunc,
		"completedCallbackFunc": completedCallbackFunc,
		"options":               options,
		"x-kunlun-loop-masks":   lookMask,
	}

	data, err := cUtils.ErrorWrapper(getFaaSInfraClient().PostJson(utils.SetAppConfToCtx(ctx, appCtx), GetPathInvokeFunctionDistributed(namespace), headers, body, cHttp.AppTokenMiddleware))
	if err != nil {
		return 0, err
	}

	var taskID int64
	err = cUtils.JsonUnmarshalBytes(data, &taskID)
	if err != nil {
		return 0, cExceptions.InternalError("Unmarshal result failed, err: %v", err)
	}

	return taskID, nil
}

func DoRequestRedis(ctx context.Context, param interface{}) ([]byte, map[string]interface{}, error) {
	ctx = cUtils.SetApiTimeoutMethodToCtx(ctx, cConstants.RequestRedis)

	data, extra, err := getFaaSInfraClient().PostJson(ctx, GetFaaSInfraPathRedis(), nil, param, cHttp.AppTokenMiddleware, cHttp.TenantAndUserMiddleware, cHttp.ServiceIDMiddleware)
	return data, extra, err
}

func DoRequestFile(ctx context.Context, contentType string, body *bytes.Buffer) ([]byte, error) {
	return cUtils.ErrorWrapper(getFaaSInfraClient().PostFormData(ctx, GetFaaSInfraPathFile(), map[string][]string{
		cConstants.HttpHeaderKeyContentType: {contentType},
	}, body, cHttp.AppTokenMiddleware, cHttp.TenantAndUserMiddleware, cHttp.ServiceIDMiddleware))
}

func doRequestMongodb(ctx context.Context, param interface{}) ([]byte, error) {
	ctx = cUtils.SetApiTimeoutMethodToCtx(ctx, cConstants.RequestMongodb)

	data, err := cUtils.ErrorWrapper(getFaaSInfraClient().PostBson(ctx, GetFaaSInfraPathMongodb(), nil, param, cHttp.AppTokenMiddleware, cHttp.TenantAndUserMiddleware, cHttp.ServiceIDMiddleware))
	if err != nil {
		return data, err
	}
	return base64.StdEncoding.DecodeString(string(data))
}

func BatchCreate(ctx context.Context, param interface{}) ([]primitive.ObjectID, error) {
	data, err := doRequestMongodb(ctx, param)
	if err != nil {
		return nil, err
	}

	var result inner.BatchCreateResult
	err = bson.Unmarshal(data, &result)
	if err != nil {
		return nil, cExceptions.InternalError("BatchCreate failed, err: %v", err)
	}

	return result.IDs, nil
}

func Create(ctx context.Context, param interface{}) (*structs.RecordOnlyId, error) {
	ids, err := BatchCreate(ctx, param)
	if err != nil {
		return nil, err
	}

	if len(ids) > 0 {
		return &structs.RecordOnlyId{ID: ids[0]}, nil
	}

	return nil, nil
}

func Find(ctx context.Context, param, results interface{}) error {
	resultsVal := reflect.ValueOf(results)
	if resultsVal.Kind() != reflect.Ptr {
		return fmt.Errorf("[Find] results argument must be a pointer to a slice, but was a %s", resultsVal.Kind())
	}

	data, err := doRequestMongodb(ctx, param)
	if err != nil {
		return err
	}

	res := &inner.RawResult{}
	res.Bind(results)

	err = bson.Unmarshal(data, res)
	if err != nil {
		return cExceptions.InternalError("[Find] Unmarshal failed, err: %v", err)
	}

	return err
}

func FindOne(ctx context.Context, param, result interface{}) error {

	resultsVal := reflect.ValueOf(result)
	if resultsVal.Kind() != reflect.Ptr {
		return fmt.Errorf("[FindOne] results argument must be a pointer to a slice, but was a %s", resultsVal.Kind())
	}

	data, err := doRequestMongodb(ctx, param)
	if err != nil {
		return err
	}

	res := &inner.RawResult{}
	res.Bind(result)

	err = bson.Unmarshal(data, res)
	if err != nil {
		return cExceptions.InternalError("[FindOne] Unmarshal failed, err: %v", err)
	}
	return nil
}

func Count(ctx context.Context, param interface{}) (int64, error) {
	data, err := doRequestMongodb(ctx, param)
	if err != nil {
		return 0, err
	}

	result := &inner.CountResult{}

	err = bson.Unmarshal(data, &result)
	if err != nil {
		return 0, cExceptions.InternalError("[Count] Unmarshal failed, err: %v", err)
	}

	return result.Data.Count, nil
}

func Distinct(ctx context.Context, param interface{}, results interface{}) error {
	resultsVal := reflect.ValueOf(results)
	if resultsVal.Kind() != reflect.Ptr {
		return fmt.Errorf("[Distinct] results argument must be a pointer to a slice, but was a %s", resultsVal.Kind())
	}

	data, err := doRequestMongodb(ctx, param)
	if err != nil {
		return err
	}

	res := &inner.RawResult{}
	res.Bind(results)

	err = bson.Unmarshal(data, &res)
	if err != nil {
		return cExceptions.InternalError("[Distinct] Unmarshal failed, err: %v", err)
	}

	return nil
}

func Update(ctx context.Context, param interface{}) error {
	_, err := doRequestMongodb(ctx, param)
	if err != nil {
		return err
	}

	return nil
}

func Delete(ctx context.Context, param interface{}) error {
	_, err := doRequestMongodb(ctx, param)
	if err != nil {
		return err
	}

	return nil
}

func ReadFromURL(ctx context.Context, targetURL string) ([]byte, error) {

	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	rsp, err := getCommonHttpClient().Do(req.WithContext(ctx))
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

func UploadWithContent(ctx context.Context, name string, content []byte, option *structs.Option) (*structs.UploadResult, error) {
	if !cUtils.IsExternalFaaS() {
		return nil, cExceptions.InvalidParamError("unsupport oss")
	}

	if len(content) > constants.MaxFileSize {
		return nil, cExceptions.InvalidParamError("file too large, exceed %v", constants.MaxFileSize)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(constants.FileFieldName, name)
	if err != nil {
		return nil, err
	}

	if _, err = part.Write(content); err != nil {
		return nil, err
	}

	if option != nil {
		var data []byte
		data, err = json.Marshal(option)
		if err != nil {
			return nil, err
		}

		if err = writer.WriteField(constants.FileOptionFieldName, string(data)); err != nil {
			return nil, err
		}

	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	out, err := DoRequestFile(ctx, writer.FormDataContentType(), body)
	if err != nil {
		return nil, err
	}
	res := &structs.FileUploadResult{}
	dest, err := base64.StdEncoding.DecodeString(string(out))
	if err != nil {
		return nil, fmt.Errorf("result decode err: %v", err)
	}
	if err = bson.Unmarshal(dest, &res); err != nil {
		return nil, err
	}

	if res.Data != nil && res.Data.UploadError == nil {
		return &structs.UploadResult{
			URL: res.Data.URL,
		}, nil
	}

	return nil, res.Data.UploadError.Error()
}
