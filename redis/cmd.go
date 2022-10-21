// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package redis

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/byted-apaas/baas-sdk-go/http"
	cExceptions "github.com/byted-apaas/server-common-go/exceptions"
	cUtils "github.com/byted-apaas/server-common-go/utils"
)

type IRedis interface {
	TTL(ctx context.Context, key string) *DurationCmd
	Type(ctx context.Context, key string) *StatusCmd
	Append(ctx context.Context, key, value string) *IntCmd
	GetRange(ctx context.Context, key string, start, end int64) *StringCmd
	GetSet(ctx context.Context, key string, value interface{}) *StringCmd
	Get(ctx context.Context, key string) *StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd
	Del(ctx context.Context, keys ...string) *IntCmd
	Exists(ctx context.Context, keys ...string) *IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *BoolCmd
	ExpireAt(ctx context.Context, key string, tm time.Time) *BoolCmd
	Persist(ctx context.Context, key string) *BoolCmd
	PExpire(ctx context.Context, key string, expiration time.Duration) *BoolCmd
	PExpireAt(ctx context.Context, key string, tm time.Time) *BoolCmd
	PTTL(ctx context.Context, key string) *DurationCmd
	Incr(ctx context.Context, key string) *IntCmd
	Decr(ctx context.Context, key string) *IntCmd
	IncrBy(ctx context.Context, key string, value int64) *IntCmd
	DecrBy(ctx context.Context, key string, value int64) *IntCmd
	IncrByFloat(ctx context.Context, key string, value float64) *FloatCmd
	MGet(ctx context.Context, keys ...string) *SliceCmd
	MSet(ctx context.Context, pairs ...interface{}) *StatusCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd
	SetXX(ctx context.Context, key string, value interface{}) *BoolCmd
	SetRange(ctx context.Context, key string, offset int64, value string) *IntCmd
	StrLen(ctx context.Context, key string) *IntCmd
	GetBit(ctx context.Context, key string, offset int64) *IntCmd
	SetBit(ctx context.Context, key string, offset int64, value int) *IntCmd
	BitCount(ctx context.Context, key string, bitCount *BitCountArgs) *IntCmd
	HDel(ctx context.Context, key string, fields ...string) *IntCmd
	HExists(ctx context.Context, key, field string) *BoolCmd
	HGet(ctx context.Context, key, field string) *StringCmd
	HGetAll(ctx context.Context, key string) *StrStrMapCmd
	HIncrBy(ctx context.Context, key, field string, incr int64) *IntCmd
	HIncrByFloat(ctx context.Context, key, field string, incr float64) *FloatCmd
	HKeys(ctx context.Context, key string) *StrSliceCmd
	HLen(ctx context.Context, key string) *IntCmd
	HMSet(ctx context.Context, key string, pairs ...interface{}) *StatusCmd
	HMGet(ctx context.Context, key string, fields ...string) *SliceCmd
	HSet(ctx context.Context, key string, field string, value interface{}) *BoolCmd
	HSetNX(ctx context.Context, key, field string, value interface{}) *BoolCmd
	HVals(ctx context.Context, key string) *StrSliceCmd
	LIndex(ctx context.Context, key string, index int64) *StringCmd
	LInsert(ctx context.Context, key, op string, pivot, value interface{}) *IntCmd
	LLen(ctx context.Context, key string) *IntCmd
	LPop(ctx context.Context, key string) *StringCmd
	LPush(ctx context.Context, key string, values ...interface{}) *IntCmd
	LPushX(ctx context.Context, key string, values ...interface{}) *IntCmd
	LRange(ctx context.Context, key string, start, stop int64) *StrSliceCmd
	LRem(ctx context.Context, key string, count int64, value interface{}) *IntCmd
	LSet(ctx context.Context, key string, index int64, value interface{}) *StatusCmd
	LTrim(ctx context.Context, key string, start, stop int64) *StatusCmd
	RPop(ctx context.Context, key string) *StringCmd
	RPush(ctx context.Context, key string, values ...interface{}) *IntCmd
	RPushX(ctx context.Context, key string, values ...interface{}) *IntCmd
	SAdd(ctx context.Context, key string, members ...interface{}) *IntCmd
	SCard(ctx context.Context, key string) *IntCmd
	SDiff(ctx context.Context, keys ...string) *StrSliceCmd
	SDiffStore(destination string, ctx context.Context, keys ...string) *IntCmd
	SInter(ctx context.Context, keys ...string) *StrSliceCmd
	SInterStore(destination string, ctx context.Context, keys ...string) *IntCmd
	SIsMember(ctx context.Context, key string, member interface{}) *BoolCmd
	SMembers(ctx context.Context, key string) *StrSliceCmd
	SMove(ctx context.Context, source, destination string, member interface{}) *BoolCmd
	SPop(ctx context.Context, key string) *StringCmd
	SPopN(ctx context.Context, key string, count int64) *StrSliceCmd
	SRandMember(ctx context.Context, key string) *StringCmd
	SRandMemberN(ctx context.Context, key string, count int64) *StrSliceCmd
	SRem(ctx context.Context, key string, members ...interface{}) *IntCmd
	SUnion(ctx context.Context, keys ...string) *StrSliceCmd
	SUnionStore(destination string, ctx context.Context, keys ...string) *IntCmd
	ZAdd(ctx context.Context, key string, members ...*Z) *IntCmd
	ZAddNX(ctx context.Context, key string, members ...*Z) *IntCmd
	ZAddXX(ctx context.Context, key string, members ...*Z) *IntCmd
	ZAddCh(ctx context.Context, key string, members ...*Z) *IntCmd
	ZAddNXCh(ctx context.Context, key string, members ...*Z) *IntCmd
	ZAddXXCh(ctx context.Context, key string, members ...*Z) *IntCmd
	ZIncr(ctx context.Context, key string, member *Z) *FloatCmd
	ZIncrNX(ctx context.Context, key string, member *Z) *FloatCmd
	ZIncrXX(ctx context.Context, key string, member *Z) *FloatCmd
	ZCard(ctx context.Context, key string) *IntCmd
	ZCount(ctx context.Context, key, min, max string) *IntCmd
	ZIncrBy(ctx context.Context, key string, increment float64, member string) *FloatCmd
	ZInterStore(ctx context.Context, destination string, store *ZStore) *IntCmd
	ZRange(ctx context.Context, key string, start, stop int64) *StrSliceCmd
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd
	ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StrSliceCmd
	ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd
	ZRank(ctx context.Context, key, member string) *IntCmd
	ZRem(ctx context.Context, key string, members ...interface{}) *IntCmd
	ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *IntCmd
	ZRemRangeByScore(ctx context.Context, key, min, max string) *IntCmd
	ZRevRange(ctx context.Context, key string, start, stop int64) *StrSliceCmd
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd
	ZRevRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StrSliceCmd
	ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd
	ZRevRank(ctx context.Context, key, member string) *IntCmd
	ZScore(ctx context.Context, key, member string) *FloatCmd
	ZUnionStore(ctx context.Context, dest string, store *ZStore) *IntCmd
	PFAdd(ctx context.Context, key string, els ...interface{}) *IntCmd
	PFCount(ctx context.Context, keys ...string) *IntCmd
	PFMerge(ctx context.Context, dest string, keys ...string) *StatusCmd
}

const Nil = ErrorRedis("redis: nil")

type ErrorRedis string

func (e ErrorRedis) Error() string { return string(e) }

type result struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *result) bind(v interface{}) {
	r.Data = v
}

type baseCmd struct {
	client *Redis
	name   string
	args   []interface{}
	err    error
	result *result
}

func (c *baseCmd) Err() error {
	return c.err
}

type redisArgumentList struct {
	Cmd  string        `json:"cmd"`
	Args []interface{} `json:"args"`
}

// RedisCmdExecution Request
func (c *baseCmd) request(ctx context.Context) {
	data, extra, e := http.DoRequestRedis(ctx, redisArgumentList{Cmd: c.name, Args: c.args})
	if e != nil {
		c.err = cExceptions.ErrWrap(e)
		return
	}

	if e := json.Unmarshal(data, c.result); e != nil {
		c.err = cExceptions.InternalError("[Redis] JsonUnmarshal failed(%v), err: %v", cUtils.GetLogIDFromExtra(extra), e)
		return
	}

	if http.HasError(c.result.Code) {
		if http.IsSysError(c.result.Code) {
			c.err = cExceptions.InternalError("[Redis] call remote failed, err: %v ([%v] %v)", c.result.Msg, c.result.Code, cUtils.GetLogIDFromExtra(extra))
		} else {
			c.err = cExceptions.InvalidParamError("[Redis] call remote failed, err: %v ([%v] %v)", c.result.Msg, c.result.Code, cUtils.GetLogIDFromExtra(extra))
		}
		return
	}

	if c.result.Data == nil {
		c.err = Nil
		return
	}
}

type StringCmd struct {
	baseCmd
	val string
}

func NewStringCmd(client *Redis, name string, args ...interface{}) *StringCmd {
	cmd := &StringCmd{
		baseCmd: baseCmd{
			client: client,
			name:   name,
			args:   args,
			result: &result{},
		},
	}
	cmd.result.bind(&cmd.val)
	return cmd
}

func (c *StringCmd) Val() string {
	return c.val
}

func (c *StringCmd) Result() (string, error) {
	return c.val, c.err
}

func (c *StringCmd) Int() (int, error) {
	if c.err != nil {
		return 0, c.err
	}
	return strconv.Atoi(c.val)
}

func (c *StringCmd) Int64() (int64, error) {
	if c.err != nil {
		return 0, c.err
	}
	return strconv.ParseInt(c.val, 10, 64)
}

func (c *StringCmd) Uint64() (uint64, error) {
	if c.err != nil {
		return 0, c.err
	}
	return strconv.ParseUint(c.val, 10, 64)
}

func (c *StringCmd) Float64() (float64, error) {
	if c.err != nil {
		return 0, c.err
	}
	return strconv.ParseFloat(c.val, 64)
}

func (c *StringCmd) Time() (time.Time, error) {
	if c.err != nil {
		return time.Time{}, c.err
	}
	return time.Parse(time.RFC3339Nano, c.val)
}

type StatusCmd struct {
	baseCmd
	val string
}

func NewStatusCmd(client *Redis, name string, args ...interface{}) *StatusCmd {
	cmd := &StatusCmd{
		baseCmd: baseCmd{
			client: client,
			name:   name,
			args:   args,
			result: &result{},
		},
	}
	cmd.result.bind(&cmd.val)
	return cmd
}

func (c *StatusCmd) Val() string {
	return c.val
}

func (c *StatusCmd) Result() (string, error) {
	return c.val, c.err
}

type IntCmd struct {
	baseCmd
	val int64
}

func NewIntCmd(client *Redis, name string, args ...interface{}) *IntCmd {
	cmd := &IntCmd{
		baseCmd: baseCmd{
			client: client,
			name:   name,
			args:   args,
			result: &result{},
		},
	}
	cmd.result.bind(&cmd.val)
	return cmd
}

func (c *IntCmd) Val() int64 {
	return c.val
}

func (c *IntCmd) Result() (int64, error) {
	return c.val, c.err
}

func (c *IntCmd) Uint64() (uint64, error) {
	return uint64(c.val), c.err
}

type DurationCmd struct {
	baseCmd
	val time.Duration
}

func NewDurationCmd(client *Redis, name string, args ...interface{}) *DurationCmd {
	cmd := &DurationCmd{
		baseCmd: baseCmd{
			client: client,
			name:   name,
			args:   args,
			result: &result{},
		},
	}
	cmd.result.bind(&cmd.val)
	return cmd
}

func (c *DurationCmd) Val() time.Duration {
	return c.val
}

func (c *DurationCmd) Result() (time.Duration, error) {
	return c.val, c.err
}

type SliceCmd struct {
	baseCmd
	val []interface{}
}

func NewSliceCmd(client *Redis, name string, args ...interface{}) *SliceCmd {
	cmd := &SliceCmd{
		baseCmd: baseCmd{
			client: client,
			name:   name,
			args:   args,
			result: &result{},
		},
	}
	cmd.result.bind(&cmd.val)
	return cmd
}

func (c *SliceCmd) Val() []interface{} {
	return c.val
}

func (c *SliceCmd) Result() ([]interface{}, error) {
	return c.val, c.err
}

type FloatCmd struct {
	baseCmd
	val float64
}

func NewFloatCmd(client *Redis, name string, args ...interface{}) *FloatCmd {
	cmd := &FloatCmd{
		baseCmd: baseCmd{
			client: client,
			name:   name,
			args:   args,
			result: &result{},
		},
	}
	cmd.result.bind(&cmd.val)
	return cmd
}

func (c *FloatCmd) Val() float64 {
	val, err := c.Result()
	if err != nil {
		return 0
	}
	return val
}

func (c *FloatCmd) Result() (float64, error) {
	return c.val, c.err
}

type BoolCmd struct {
	baseCmd
	val bool
}

func NewBoolCmd(client *Redis, name string, args ...interface{}) *BoolCmd {
	cmd := &BoolCmd{
		baseCmd: baseCmd{
			client: client,
			name:   name,
			args:   args,
			result: &result{},
		},
	}
	cmd.result.bind(&cmd.val)
	return cmd
}

func (c *BoolCmd) Val() bool {
	return c.val
}

func (c *BoolCmd) Result() (bool, error) {
	return c.val, c.err
}

type StrStrMapCmd struct {
	baseCmd
	val map[string]string
}

func NewStrStrMapCmd(client *Redis, name string, args ...interface{}) *StrStrMapCmd {
	cmd := &StrStrMapCmd{
		baseCmd: baseCmd{
			client: client,
			name:   name,
			args:   args,
			result: &result{},
		},
	}
	cmd.result.bind(&cmd.val)
	return cmd
}

func (c *StrStrMapCmd) Val() map[string]string {
	return c.val
}

func (c *StrStrMapCmd) Result() (map[string]string, error) {
	return c.val, c.err
}

type StrSliceCmd struct {
	baseCmd
	val []string
}

func NewStrSliceCmd(client *Redis, name string, args ...interface{}) *StrSliceCmd {
	cmd := &StrSliceCmd{
		baseCmd: baseCmd{
			client: client,
			name:   name,
			args:   args,
			result: &result{},
		},
	}
	cmd.result.bind(&cmd.val)
	return cmd
}

func (c *StrSliceCmd) Val() []string {
	return c.val
}

func (c *StrSliceCmd) Result() ([]string, error) {
	return c.val, c.err
}

type ZSliceCmd struct {
	baseCmd
	val ZSlice
}

func NewZSliceCmd(client *Redis, name string, args ...interface{}) *ZSliceCmd {
	cmd := &ZSliceCmd{
		baseCmd: baseCmd{
			client: client,
			name:   name,
			args:   args,
			result: &result{},
		},
	}
	cmd.result.bind(&cmd.val)
	return cmd
}

func (c *ZSliceCmd) Val() []Z {
	return c.val
}

func (c *ZSliceCmd) Result() ([]Z, error) {
	return c.val, c.err
}
