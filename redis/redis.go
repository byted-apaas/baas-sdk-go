// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package redis

import (
	"context"
	"time"
)

type Redis struct{}

func NewRedis() *Redis {
	return &Redis{}
}

func (c *Redis) TTL(ctx context.Context, key string) *DurationCmd {
	cmd := NewDurationCmd(c, "ttl", key)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) Type(ctx context.Context, key string) *StatusCmd {
	cmd := NewStatusCmd(c, "type", key)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) Append(ctx context.Context, key, value string) *IntCmd {
	cmd := NewIntCmd(c, "append", key, value)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) GetRange(ctx context.Context, key string, start, end int64) *StringCmd {
	cmd := NewStringCmd(c, "getrange", key, start, end)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) GetSet(ctx context.Context, key string, value interface{}) *StringCmd {
	cmd := NewStringCmd(c, "getset", key, value)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) Get(ctx context.Context, key string) *StringCmd {
	cmd := NewStringCmd(c, "get", key)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd {
	args := []interface{}{key, value}
	if expiration > 0 {
		if usePrecise(expiration) {
			args = append(args, "px", formatMils(expiration))
		} else {
			args = append(args, "ex", formatSecond(expiration))
		}
	}
	cmd := NewStatusCmd(c, "set", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) Del(ctx context.Context, keys ...string) *IntCmd {
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = key
	}
	cmd := NewIntCmd(c, "del", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) Exists(ctx context.Context, keys ...string) *IntCmd {
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = key
	}
	cmd := NewIntCmd(c, "exists", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) Expire(ctx context.Context, key string, expiration time.Duration) *BoolCmd {
	cmd := NewIntCmd(c, "expire", key, formatSecond(expiration))
	cmd.request(ctx)
	return tranIntCmd2BoolCmd(cmd)
}

func (c *Redis) ExpireAt(ctx context.Context, key string, tm time.Time) *BoolCmd {
	cmd := NewIntCmd(c, "expireat", key, tm.Unix())
	cmd.request(ctx)
	return tranIntCmd2BoolCmd(cmd)
}

func (c *Redis) Persist(ctx context.Context, key string) *BoolCmd {
	cmd := NewIntCmd(c, "persist", key)
	cmd.request(ctx)
	return tranIntCmd2BoolCmd(cmd)
}

func (c *Redis) PExpire(ctx context.Context, key string, expiration time.Duration) *BoolCmd {
	cmd := NewIntCmd(c, "pexpire", key, formatMils(expiration))
	cmd.request(ctx)
	return tranIntCmd2BoolCmd(cmd)
}

func (c *Redis) PExpireAt(ctx context.Context, key string, tm time.Time) *BoolCmd {
	cmd := NewIntCmd(c, "pexpireat", key, tm.UnixNano()/int64(time.Millisecond))
	cmd.request(ctx)
	return tranIntCmd2BoolCmd(cmd)
}

func (c *Redis) PTTL(ctx context.Context, key string) *DurationCmd {
	cmd := NewDurationCmd(c, "pttl", key)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) Incr(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd(c, "incr", key)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) Decr(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd(c, "decr", key)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) IncrBy(ctx context.Context, key string, value int64) *IntCmd {
	cmd := NewIntCmd(c, "incrby", key, value)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) DecrBy(ctx context.Context, key string, value int64) *IntCmd {
	cmd := NewIntCmd(c, "decrby", key, value)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) IncrByFloat(ctx context.Context, key string, value float64) *FloatCmd {
	cmd := NewFloatCmd(c, "incrbyfloat", key, value)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) MGet(ctx context.Context, keys ...string) *SliceCmd {
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = key
	}
	cmd := NewSliceCmd(c, "mget", args...)
	cmd.request(ctx)
	return cmd
}

// MSet pairs can be map or indefinite parameter
// eg:
// 1: map[string]interface{}{"k1":"v1","k2":"v2"}
// 2: "k1", "v1", "k2", "v2"
func (c *Redis) MSet(ctx context.Context, pairs ...interface{}) *StatusCmd {
	cmd := NewStatusCmd(c, "mset", pairs...)
	cmd.request(ctx)
	return cmd
}

// SetNX is short for "SET If Not Exists".
func (c *Redis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd {
	var statusCmd *StatusCmd
	args := []interface{}{key, value}
	if expiration <= 0 {
		intCmd := NewIntCmd(c, "setnx", key, value)
		intCmd.request(ctx)
		return tranIntCmd2BoolCmd(intCmd)
	}
	if usePrecise(expiration) {
		args = append(args, "px", formatMils(expiration), "nx")
	} else {
		args = append(args, "ex", formatSecond(expiration), "nx")
	}
	statusCmd = NewStatusCmd(c, "set", args...)
	statusCmd.request(ctx)
	return tranStatusCmd2BoolCmd(statusCmd)
}

// SetXX is short for "SET If Exists".
func (c *Redis) SetXX(ctx context.Context, key string, value interface{}) *BoolCmd {
	cmd := NewStatusCmd(c, "set", key, value, "xx")
	cmd.request(ctx)
	return tranStatusCmd2BoolCmd(cmd)
}

func (c *Redis) SetRange(ctx context.Context, key string, offset int64, value string) *IntCmd {
	cmd := NewIntCmd(c, "setrange", key, offset, value)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) StrLen(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd(c, "strlen", key)
	cmd.request(ctx)
	return cmd
}

//------------------------------------------------------------------------------
// Bit

func (c *Redis) GetBit(ctx context.Context, key string, offset int64) *IntCmd {
	cmd := NewIntCmd(c, "getbit", key, offset)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) SetBit(ctx context.Context, key string, offset int64, value int) *IntCmd {
	cmd := NewIntCmd(c, "setbit", key, offset, value)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) BitCount(ctx context.Context, key string, bitCount *BitCountArgs) *IntCmd {
	args := []interface{}{key}
	if bitCount != nil {
		args = append(args, bitCount.Start, bitCount.End)
	}
	cmd := NewIntCmd(c, "bitcount", args...)
	cmd.request(ctx)
	return cmd
}

//------------------------------------------------------------------------------
// Hash
func (c *Redis) HDel(ctx context.Context, key string, fields ...string) *IntCmd {
	args := make([]interface{}, 1+len(fields))
	args[0] = key
	for i, field := range fields {
		args[i+1] = field
	}
	cmd := NewIntCmd(c, "hdel", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) HExists(ctx context.Context, key, field string) *BoolCmd {
	cmd := NewIntCmd(c, "hexists", key, field)
	cmd.request(ctx)
	return tranIntCmd2BoolCmd(cmd)
}

func (c *Redis) HGet(ctx context.Context, key, field string) *StringCmd {
	cmd := NewStringCmd(c, "hget", key, field)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) HGetAll(ctx context.Context, key string) *StrStrMapCmd {
	cmd := NewStrStrMapCmd(c, "hgetall", key)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) HIncrBy(ctx context.Context, key, field string, incrI64 int64) *IntCmd {
	cmd := NewIntCmd(c, "hincrby", key, field, incrI64)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) HIncrByFloat(ctx context.Context, key, field string, incrF64 float64) *FloatCmd {
	cmd := NewFloatCmd(c, "hincrbyfloat", key, field, incrF64)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) HKeys(ctx context.Context, key string) *StrSliceCmd {
	cmd := NewStrSliceCmd(c, "hkeys", key)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) HLen(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd(c, "hlen", key)
	cmd.request(ctx)
	return cmd
}

// HMSet pairs can be map or indefinite parameter
// eg:
// 1: map[string]interface{}{"k1":"v1","k2":"v2"}
// 2: "k1", "v1", "k2", "v2"
func (c *Redis) HMSet(ctx context.Context, key string, pairs ...interface{}) *StatusCmd {
	args := make([]interface{}, 1, 1+len(pairs))
	args[0] = key
	args = appendArgs(args, pairs)
	cmd := NewStatusCmd(c, "hmset", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) HMGet(ctx context.Context, key string, fields ...string) *SliceCmd {
	args := make([]interface{}, 1+len(fields))
	args[0] = key
	for i, field := range fields {
		args[i+1] = field
	}
	cmd := NewSliceCmd(c, "hmget", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) HSet(ctx context.Context, key string, field string, value interface{}) *BoolCmd {
	cmd := NewIntCmd(c, "hset", key, field, value)
	cmd.request(ctx)
	return tranIntCmd2BoolCmd(cmd)
}

func (c *Redis) HSetNX(ctx context.Context, key, field string, value interface{}) *BoolCmd {
	cmd := NewIntCmd(c, "hsetnx", key, field, value)
	cmd.request(ctx)
	return tranIntCmd2BoolCmd(cmd)
}

func (c *Redis) HVals(ctx context.Context, key string) *StrSliceCmd {
	cmd := NewStrSliceCmd(c, "hvals", key)
	cmd.request(ctx)
	return cmd
}

//------------------------------------------------------------------------------
// List
func (c *Redis) LIndex(ctx context.Context, key string, index int64) *StringCmd {
	cmd := NewStringCmd(c, "lindex", key, index)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *IntCmd {
	cmd := NewIntCmd(c, "linsert", key, op, pivot, value)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) LLen(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd(c, "llen", key)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) LPop(ctx context.Context, key string) *StringCmd {
	cmd := NewStringCmd(c, "lpop", key)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) LPush(ctx context.Context, key string, values ...interface{}) *IntCmd {
	args := make([]interface{}, 1, 1+len(values))
	args[0] = key
	args = appendArgs(args, values)
	cmd := NewIntCmd(c, "lpush", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) LPushX(ctx context.Context, key string, values ...interface{}) *IntCmd {
	args := make([]interface{}, 1, 1+len(values))
	args[0] = key
	args = appendArgs(args, values)
	cmd := NewIntCmd(c, "lpushx", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) LRange(ctx context.Context, key string, start, stop int64) *StrSliceCmd {
	cmd := NewStrSliceCmd(c, "lrange", key, start, stop)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) LRem(ctx context.Context, key string, count int64, value interface{}) *IntCmd {
	cmd := NewIntCmd(c, "lrem", key, count, value)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) LSet(ctx context.Context, key string, index int64, value interface{}) *StatusCmd {
	cmd := NewStatusCmd(c, "lset", key, index, value)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) LTrim(ctx context.Context, key string, start, stop int64) *StatusCmd {
	cmd := NewStatusCmd(c, "ltrim", key, start, stop)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) RPop(ctx context.Context, key string) *StringCmd {
	cmd := NewStringCmd(c, "rpop", key)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) RPush(ctx context.Context, key string, values ...interface{}) *IntCmd {
	args := make([]interface{}, 1, 1+len(values))
	args[0] = key
	args = append(args, values...)
	cmd := NewIntCmd(c, "rpush", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) RPushX(ctx context.Context, key string, values ...interface{}) *IntCmd {
	args := make([]interface{}, 1, 1+len(values))
	args[0] = key
	args = append(args, values...)
	cmd := NewIntCmd(c, "rpushx", args...)
	cmd.request(ctx)
	return cmd
}

//------------------------------------------------------------------------------
// Set

func (c *Redis) SAdd(ctx context.Context, key string, members ...interface{}) *IntCmd {
	args := make([]interface{}, 1, 1+len(members))
	args[0] = key
	args = append(args, members...)
	cmd := NewIntCmd(c, "sadd", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) SCard(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd(c, "scard", key)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) SDiff(ctx context.Context, keys ...string) *StrSliceCmd {
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = key
	}
	cmd := NewStrSliceCmd(c, "sdiff", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) SDiffStore(destination string, ctx context.Context, keys ...string) *IntCmd {
	args := make([]interface{}, 1+len(keys))
	args[0] = destination
	for i, key := range keys {
		args[1+i] = key
	}
	cmd := NewIntCmd(c, "sdiffstore", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) SInter(ctx context.Context, keys ...string) *StrSliceCmd {
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = key
	}
	cmd := NewStrSliceCmd(c, "sinter", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) SInterStore(destination string, ctx context.Context, keys ...string) *IntCmd {
	args := make([]interface{}, 1+len(keys))
	args[0] = destination
	for i, key := range keys {
		args[1+i] = key
	}
	cmd := NewIntCmd(c, "sinterstore", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) SIsMember(ctx context.Context, key string, member interface{}) *BoolCmd {
	cmd := NewIntCmd(c, "sismember", key, member)
	cmd.request(ctx)
	return tranIntCmd2BoolCmd(cmd)
}

// SMembers `SMEMBERS key` command output as a slice.
func (c *Redis) SMembers(ctx context.Context, key string) *StrSliceCmd {
	cmd := NewStrSliceCmd(c, "smembers", key)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) SMove(ctx context.Context, source, destination string, member interface{}) *BoolCmd {
	cmd := NewIntCmd(c, "smove", source, destination, member)
	cmd.request(ctx)
	return tranIntCmd2BoolCmd(cmd)
}

// SPop `SPOP key` command.
func (c *Redis) SPop(ctx context.Context, key string) *StringCmd {
	cmd := NewStringCmd(c, "spop", key)
	cmd.request(ctx)
	return cmd
}

// SPOP `SPOP ctx context.Context, key count` command.
func (c *Redis) SPopN(ctx context.Context, key string, count int64) *StrSliceCmd {
	cmd := NewStrSliceCmd(c, "spop", key, count)
	cmd.request(ctx)
	return cmd
}

// SRandMember `SRANDMEMBER key` command.
func (c *Redis) SRandMember(ctx context.Context, key string) *StringCmd {
	cmd := NewStringCmd(c, "srandmember", key)
	cmd.request(ctx)
	return cmd
}

// SRandMemberN `SRANDMEMBER ctx context.Context, key count` command.
func (c *Redis) SRandMemberN(ctx context.Context, key string, count int64) *StrSliceCmd {
	cmd := NewStrSliceCmd(c, "srandmember", key, count)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) SRem(ctx context.Context, key string, members ...interface{}) *IntCmd {
	args := make([]interface{}, 1, 1+len(members))
	args[0] = key
	args = append(args, members...)
	cmd := NewIntCmd(c, "srem", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) SUnion(ctx context.Context, keys ...string) *StrSliceCmd {
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = key
	}
	cmd := NewStrSliceCmd(c, "sunion", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) SUnionStore(destination string, ctx context.Context, keys ...string) *IntCmd {
	args := make([]interface{}, 1+len(keys))
	args[0] = destination
	for i, key := range keys {
		args[1+i] = key
	}
	cmd := NewIntCmd(c, "sunionstore", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) executeZSetIntCmd(ctx context.Context, name string, key string, members ...*Z) *IntCmd {
	l := make([]interface{}, 2*len(members)+1)
	l[0] = key
	for i := range members {
		l[2*i+1] = members[i].Score
		l[2*i+2] = members[i].Member
	}

	cmd := NewIntCmd(c, name, l...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) executeZSetFloatCmd(ctx context.Context, name string, key string, members ...*Z) *FloatCmd {
	l := make([]interface{}, 2*len(members)+1)
	l[0] = key
	for i := range members {
		l[i+1] = members[i].Score
		l[i+2] = members[i].Member
	}

	cmd := NewFloatCmd(c, name, l...)
	cmd.request(ctx)
	return cmd
}

// ZAdd `ZADD ctx context.Context, key score member [score member ...]` command.
func (c *Redis) ZAdd(ctx context.Context, key string, members ...*Z) *IntCmd {
	return c.executeZSetIntCmd(ctx, "zadd", key, members...)
}

// ZAddNX `ZADD ctx context.Context, key NX score member [score member ...]` command.
func (c *Redis) ZAddNX(ctx context.Context, key string, members ...*Z) *IntCmd {
	return c.executeZSetIntCmd(ctx, "zaddnx", key, members...)
}

// ZAddXX `ZADD ctx context.Context, key XX score member [score member ...]` command.
func (c *Redis) ZAddXX(ctx context.Context, key string, members ...*Z) *IntCmd {
	return c.executeZSetIntCmd(ctx, "zaddxx", key, members...)
}

// ZAddCh `ZADD ctx context.Context, key CH score member [score member ...]` command.
func (c *Redis) ZAddCh(ctx context.Context, key string, members ...*Z) *IntCmd {
	return c.executeZSetIntCmd(ctx, "zaddch", key, members...)
}

// ZAddNXCh `ZADD ctx context.Context, key NX CH score member [score member ...]` command.
func (c *Redis) ZAddNXCh(ctx context.Context, key string, members ...*Z) *IntCmd {
	return c.executeZSetIntCmd(ctx, "zaddnxch", key, members...)
}

// ZAddXXCh `ZADD ctx context.Context, key XX CH score member [score member ...]` command.
func (c *Redis) ZAddXXCh(ctx context.Context, key string, members ...*Z) *IntCmd {
	return c.executeZSetIntCmd(ctx, "zaddxxch", key, members...)
}

// ZIncr `ZADD ctx context.Context, key INCR score member` command.
func (c *Redis) ZIncr(ctx context.Context, key string, member *Z) *FloatCmd {
	return c.executeZSetFloatCmd(ctx, "zincr", key, member)
}

// ZIncrNX `ZADD ctx context.Context, key NX INCR score member` command.
func (c *Redis) ZIncrNX(ctx context.Context, key string, member *Z) *FloatCmd {
	return c.executeZSetFloatCmd(ctx, "zincrnx", key, member)
}

// ZIncrXX `ZADD ctx context.Context, key XX INCR score member` command.
func (c *Redis) ZIncrXX(ctx context.Context, key string, member *Z) *FloatCmd {
	return c.executeZSetFloatCmd(ctx, "zincrxx", key, member)
}

func (c *Redis) ZCard(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd(c, "zcard", key)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZCount(ctx context.Context, key, min, max string) *IntCmd {
	cmd := NewIntCmd(c, "zcount", key, min, max)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZIncrBy(ctx context.Context, key string, increment float64, member string) *FloatCmd {
	cmd := NewFloatCmd(c, "zincrby", key, increment, member)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZInterStore(ctx context.Context, destination string, store *ZStore) *IntCmd {
	args := make([]interface{}, 2+len(store.Keys))
	args[0] = destination
	args[1] = len(store.Keys)
	for i, key := range store.Keys {
		args[2+i] = key
	}
	if len(store.Weights) > 0 {
		args = append(args, "weights")
		for _, weight := range store.Weights {
			args = append(args, weight)
		}
	}
	if store.Aggregate != "" {
		args = append(args, "aggregate", store.Aggregate)
	}
	cmd := NewIntCmd(c, "zinterstore", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) zRange(ctx context.Context, key string, start, stop int64, withScores bool) *StrSliceCmd {
	args := []interface{}{key, start, stop}
	if withScores {
		args = append(args, "withscores")
	}
	cmd := NewStrSliceCmd(c, "zrange", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZRange(ctx context.Context, key string, start, stop int64) *StrSliceCmd {
	return c.zRange(ctx, key, start, stop, false)
}

func (c *Redis) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd {
	cmd := NewZSliceCmd(c, "zrange", key, start, stop, "withscores")
	cmd.request(ctx)
	return cmd
}

func (c *Redis) zRangeBy(ctx context.Context, zcmd string, key string, opt *ZRangeBy, withScores bool) *StrSliceCmd {
	args := []interface{}{key, opt.Min, opt.Max}
	if withScores {
		args = append(args, "withscores")
	}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(args, "limit", opt.Offset, opt.Count)
	}
	cmd := NewStrSliceCmd(c, zcmd, args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StrSliceCmd {
	return c.zRangeBy(ctx, "zrangebyscore", key, opt, false)
}

func (c *Redis) ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd {
	args := []interface{}{key, opt.Min, opt.Max, "withscores"}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(args, "limit", opt.Offset, opt.Count)
	}
	cmd := NewZSliceCmd(c, "zrangebyscore", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZRank(ctx context.Context, key, member string) *IntCmd {
	cmd := NewIntCmd(c, "zrank", key, member)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZRem(ctx context.Context, key string, members ...interface{}) *IntCmd {
	args := make([]interface{}, 2, 2+len(members))
	args[0] = key
	args = appendArgs(args, members)
	cmd := NewIntCmd(c, "zrem", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *IntCmd {
	cmd := NewIntCmd(c, "zremrangebyrank", key, start, stop)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZRemRangeByScore(ctx context.Context, key, min, max string) *IntCmd {
	cmd := NewIntCmd(c, "zremrangebyscore", key, min, max)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZRevRange(ctx context.Context, key string, start, stop int64) *StrSliceCmd {
	cmd := NewStrSliceCmd(c, "zrevrange", key, start, stop)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd {
	cmd := NewZSliceCmd(c, "zrevrange", key, start, stop, "withscores")
	cmd.request(ctx)
	return cmd
}

func (c *Redis) zRevRangeBy(ctx context.Context, zcmd, key string, opt *ZRangeBy) *StrSliceCmd {
	args := []interface{}{key, opt.Max, opt.Min}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(args, "limit", opt.Offset, opt.Count)
	}
	cmd := NewStrSliceCmd(c, zcmd, args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZRevRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StrSliceCmd {
	return c.zRevRangeBy(ctx, "zrevrangebyscore", key, opt)
}

func (c *Redis) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd {
	args := []interface{}{key, opt.Min, opt.Max, "withscores"}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(args, "limit", opt.Offset, opt.Count)
	}
	cmd := NewZSliceCmd(c, "zrevrangebyscore", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZRevRank(ctx context.Context, key, member string) *IntCmd {
	cmd := NewIntCmd(c, "zrevrank", key, member)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZScore(ctx context.Context, key, member string) *FloatCmd {
	cmd := NewFloatCmd(c, "zscore", key, member)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) ZUnionStore(ctx context.Context, dest string, store *ZStore) *IntCmd {
	args := make([]interface{}, 2+len(store.Keys))
	args[0] = dest
	args[1] = len(store.Keys)
	for i, key := range store.Keys {
		args[2+i] = key
	}
	if len(store.Weights) > 0 {
		args = append(args, "weights")
		for _, weight := range store.Weights {
			args = append(args, weight)
		}
	}
	if store.Aggregate != "" {
		args = append(args, "aggregate", store.Aggregate)
	}

	cmd := NewIntCmd(c, "zunionstore", args...)
	cmd.request(ctx)
	return cmd
}

//------------------------------------------------------------------------------
// HyperLogLog

func (c *Redis) PFAdd(ctx context.Context, key string, els ...interface{}) *IntCmd {
	args := make([]interface{}, 1, 1+len(els))
	args[0] = key
	args = appendArgs(args, els)
	cmd := NewIntCmd(c, "pfadd", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) PFCount(ctx context.Context, keys ...string) *IntCmd {
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = key
	}
	cmd := NewIntCmd(c, "pfcount", args...)
	cmd.request(ctx)
	return cmd
}

func (c *Redis) PFMerge(ctx context.Context, dest string, keys ...string) *StatusCmd {
	args := make([]interface{}, 1+len(keys))
	args[0] = dest
	for i, key := range keys {
		args[1+i] = key
	}
	cmd := NewStatusCmd(c, "pfmerge", args...)
	cmd.request(ctx)
	return cmd
}
