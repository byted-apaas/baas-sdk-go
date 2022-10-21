// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package redis

import "time"

// KeepTTL is an option for Set command to keep key's existing TTL.
// For example:
// Set(key, value, redis.KeepTTL)
const KeepTTL = -1

func usePrecise(dur time.Duration) bool {
	return dur < time.Second || dur%time.Second != 0
}

func formatMils(arg time.Duration) int64 {
	if arg > 0 && arg < time.Millisecond {
		return 1
	}
	return int64(arg / time.Millisecond)
}

func formatSecond(arg time.Duration) int64 {
	if arg > 0 && arg < time.Second {
		return 1
	}
	return int64(arg / time.Second)
}

func appendArgs(destination, source []interface{}) []interface{} {
	if 1 == len(source) {
		return appendArg(destination, source[0])
	}

	destination = append(destination, source...)
	return destination
}

func appendArg(destination []interface{}, argument interface{}) []interface{} {
	switch curArg := argument.(type) {
	case map[string]interface{}:
		for k, v := range curArg {
			destination = append(destination, k, v)
		}
		return destination
	case []string:
		for _, s := range curArg {
			destination = append(destination, s)
		}
		return destination
	case []interface{}:
		destination = append(destination, curArg...)
		return destination
	default:
		return append(destination, curArg)
	}
}

func tranIntCmd2BoolCmd(intCmd *IntCmd) *BoolCmd {
	boolCmd := NewBoolCmd(intCmd.client, "expire", intCmd.args...)
	boolCmd.result = intCmd.result
	boolCmd.err = intCmd.err
	if intCmd.val == int64(0) {
		boolCmd.val = false
	} else {
		boolCmd.val = true
	}
	return boolCmd
}

func tranStatusCmd2BoolCmd(intCmd *StatusCmd) *BoolCmd {
	boolCmd := NewBoolCmd(intCmd.client, "expire", intCmd.args...)
	boolCmd.result = intCmd.result
	boolCmd.err = intCmd.err
	if intCmd.val == "" {
		boolCmd.val = false
	} else {
		boolCmd.val = true
	}
	return boolCmd
}
