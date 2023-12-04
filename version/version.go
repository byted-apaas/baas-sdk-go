// Package version defines version of server-sdk-go.
package version

import (
	"sync"

	cVersion "github.com/byted-apaas/server-common-go/version"
)

const Version = "v0.0.10"

const SDKName = "byted-apaas/baas-sdk-go"

type BaaSSDKInfo struct{}

func (b *BaaSSDKInfo) GetVersion() string {
	return Version
}

func (b *BaaSSDKInfo) GetSDKName() string {
	return SDKName
}

var (
	baasSDKInfoOnce sync.Once
	baasSDKInfo     cVersion.ISDKInfo
)

func GetBaaSSDKInfo() cVersion.ISDKInfo {
	if baasSDKInfo == nil {
		baasSDKInfoOnce.Do(func() {
			baasSDKInfo = &BaaSSDKInfo{}
		})
	}
	return baasSDKInfo
}
