package faasinfra

import (
	"context"
	"net/http"
	"sync"
	"time"

	cConstants "github.com/byted-apaas/server-common-go/constants"
	cHttp "github.com/byted-apaas/server-common-go/http"
)

var (
	fsInfraClientOnce sync.Once
	fsInfraClient     *cHttp.HttpClient
)

func getFaaSInfraClient(ctx context.Context) *cHttp.HttpClient {
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
		}
	})
	return fsInfraClient
}
