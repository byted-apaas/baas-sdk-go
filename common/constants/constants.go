// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package constants

import (
	"fmt"

	cConstants "github.com/byted-apaas/server-common-go/constants"
)

const (
	// MaxFileSize is the limitation size of file
	MaxFileSize = 30 * 1024 * 1024

	// FileFieldName Form-Data fieldKey
	FileFieldName       = "fileFieldName"
	FileOptionFieldName = "fileOption"
)

type PlatformEnvType int

const (
	PlatformEnvDEV PlatformEnvType = iota + 1
	PlatformEnvUAT
	PlatformEnvLR
	PlatformEnvPRE
	PlatformEnvOnline
)

func (p PlatformEnvType) String() string {
	switch p {
	case PlatformEnvUAT:
		return cConstants.EnvTypeStaging
	case PlatformEnvLR:
		return cConstants.EnvTypeLr
	case PlatformEnvPRE:
		return cConstants.EnvTypeGray
	case PlatformEnvOnline:
		return cConstants.EnvTypeOnline
	}
	fmt.Printf("invalid platform env type %d", p)
	return ""
}

type OperationType int
