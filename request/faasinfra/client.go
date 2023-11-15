package faasinfra

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/byted-apaas/baas-sdk-go/version"
	cConstants "github.com/byted-apaas/server-common-go/constants"
	cHttp "github.com/byted-apaas/server-common-go/http"
	"github.com/byted-apaas/server-common-go/utils"
	cUtils "github.com/byted-apaas/server-common-go/utils"
)

var (
	fsInfraClientOnce sync.Once
	fsInfraClient     *cHttp.HttpClient

	httpClientOnce sync.Once
	httpClient     *http.Client
)

func getFaaSInfraClient() *cHttp.HttpClient {
	fsInfraClientOnce.Do(func() {
		fsInfraClient = &cHttp.HttpClient{
			Type: cHttp.FaaSInfraClient,
			Client: http.Client{
				Transport: &http.Transport{
					DialContext:         cHttp.TimeoutDialer(cConstants.HttpClientDialTimeoutDefault, 0),
					TLSHandshakeTimeout: cConstants.HttpClientTLSTimeoutDefault,
					MaxIdleConns:        1000,
					MaxIdleConnsPerHost: 10,
					IdleConnTimeout:     30 * time.Second,
				},
			},
			FromSDK: version.GetBaaSSDKInfo(),
		}
	})
	if cUtils.EnableMesh() {
		fsInfraClient.MeshClient = &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					unixAddr, err := net.ResolveUnixAddr("unix", utils.GetSocketAddr())
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
	return fsInfraClient
}

func getCommonHttpClient() *http.Client {
	httpClientOnce.Do(func() {
		httpClient = &http.Client{
			Transport: &http.Transport{
				DialContext:         cHttp.TimeoutDialer(cConstants.HttpClientDialTimeoutDefault, 0),
				TLSHandshakeTimeout: cConstants.HttpClientTLSTimeoutDefault,
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     30 * time.Second,
			},
		}
	})
	return httpClient
}
