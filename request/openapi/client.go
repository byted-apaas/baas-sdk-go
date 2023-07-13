package openapi

import (
	"net/http"
	"sync"
	"time"

	cConstants "github.com/byted-apaas/server-common-go/constants"
	cHttp "github.com/byted-apaas/server-common-go/http"
)

var (
	openapiClientOnce sync.Once
	openapiClient     *cHttp.HttpClient
)

func getOpenapiClient() *cHttp.HttpClient {
	openapiClientOnce.Do(func() {
		openapiClient = &cHttp.HttpClient{
			Type: cHttp.OpenAPIClient,
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
	return openapiClient
}
