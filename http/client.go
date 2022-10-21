// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package http

import (
	"context"
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	cConstants "github.com/byted-apaas/server-common-go/constants"
	cHttp "github.com/byted-apaas/server-common-go/http"
	cUtils "github.com/byted-apaas/server-common-go/utils"
)

var (
	fsInfraOnce   sync.Once
	fsInfraClient *cHttp.HttpClient
)

func getFaaSInfraClient() *cHttp.HttpClient {
	fsInfraOnce.Do(func() {
		fsInfraClient = &cHttp.HttpClient{
			Type: cHttp.FaaSInfraClient,
			Client: http.Client{
				Transport: &http.Transport{
					DialContext:         cHttp.TimeoutDialer(cConstants.HttpClientDialTimeoutDefault, 0),
					TLSHandshakeTimeout: cConstants.HttpClientTLSTimeoutDefault,
					MaxIdleConns:        1000,
					MaxIdleConnsPerHost: 10,
					IdleConnTimeout:     60 * time.Second,
				},
			},
		}
	})
	return fsInfraClient
}

func doRequestMongodb(ctx context.Context, param interface{}) ([]byte, error) {
	ctx = cUtils.SetApiTimeoutMethodToCtx(ctx, cConstants.RequestMongodb)

	data, err := errorWrapper(getFaaSInfraClient().PostBson(ctx, GetFaaSInfraPathMongodb(), nil, param, cHttp.AppTokenMiddleware, cHttp.TenantAndUserMiddleware, cHttp.ServiceIDMiddleware))
	if err != nil {
		return data, err
	}
	return base64.StdEncoding.DecodeString(string(data))
}

func DoRequestRedis(ctx context.Context, param interface{}) ([]byte, map[string]interface{}, error) {
	ctx = cUtils.SetApiTimeoutMethodToCtx(ctx, cConstants.RequestRedis)

	data, extra, err := getFaaSInfraClient().PostJson(ctx, GetFaaSInfraPathRedis(), nil, param, cHttp.AppTokenMiddleware, cHttp.TenantAndUserMiddleware, cHttp.ServiceIDMiddleware)
	return data, extra, err
}
