package openapi

import (
	cUtils "github.com/byted-apaas/server-common-go/utils"
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/byted-apaas/baas-sdk-go/version"
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
			FromSDK: version.GetBaaSSDKInfo(),
		}
	})
	if cUtils.EnableMesh() {
		openapiClient.MeshClient = &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					unixAddr, err := net.ResolveUnixAddr("unix", cUtils.GetSocketAddr())
					if err != nil {
						return nil, err
					}
					return net.DialUnix("unix", nil, unixAddr)
				},
				TLSHandshakeTimeout: cConstants.HttpClientTLSTimeoutDefault,
				MaxIdleConns:        1000,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     60 * time.Second,
			},
		}
	}
	return openapiClient
}
