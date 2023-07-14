package redisx

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var CTL *Control

// redis控制器
type Control struct {
	Config *Config            // redis配置
	Cmd    redis.Cmdable      // redis接口方法
	Ctx    context.Context    // 上下文
	Cancel context.CancelFunc // Cancel方法
}

// 初始化redis控制器
func InitRedisCTL(conf *Config) {
	if conf.Host == "" {
		return
	}
	var err error
	var msg = conf.Format()
	if CTL == nil {
		CTL = conf.NewRedisCTL()
		_, err = CTL.Ping().Result()
		if err != nil {
			log.Error("初始化redis连接-失败！", msg)
			log.Error("error : ", err)
		} else {
			log.Info("初始化redis连接-成功！", msg)
		}
	} else {
		var newCmd = conf.NewRedisCmdable()
		CTL.Cmd = newCmd
		_, err = CTL.Ping().Result()
		if err != nil {
			log.Error("更新redis连接-失败！", msg)
			log.Error("error : ", err)
		} else {
			CTL.Config = CONFIG
			log.Info("更新redis连接-成功！", msg)
		}
	}
}

// 更新上下文配置
func (ctl *Control) SetContext(ctx context.Context) {
	ctl.Ctx = ctx
	return
}

// 更新上下文配置
func (ctl *Control) SetCancel(fn context.CancelFunc) {
	ctl.Cancel = fn
	return
}

func (ctl *Control) Pipeline() redis.Pipeliner {
	return ctl.Cmd.Pipeline()
}

func (ctl *Control) Pipelined(f func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return ctl.Cmd.Pipelined(ctl.Ctx, f)
}

func (ctl *Control) TxPipelined(f func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return ctl.Cmd.TxPipelined(ctl.Ctx, f)
}

func (ctl *Control) TxPipeline() redis.Pipeliner {
	return ctl.Cmd.TxPipeline()
}

func (ctl *Control) Command() *redis.CommandsInfoCmd {
	return ctl.Cmd.Command(ctl.Ctx)
}

func (ctl *Control) ClientGetName() *redis.StringCmd {
	return ctl.Cmd.ClientGetName(ctl.Ctx)
}

func (ctl *Control) Echo(message interface{}) *redis.StringCmd {
	return ctl.Cmd.Echo(ctl.Ctx, message)
}

func (ctl *Control) Ping() *redis.StatusCmd {
	return ctl.Cmd.Ping(ctl.Ctx)
}
func (ctl *Control) HScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return ctl.Cmd.HScan(ctl.Ctx, key, cursor, match, count)
}

func (ctl *Control) Quit() *redis.StatusCmd {
	return ctl.Cmd.Quit(ctl.Ctx)
}

func (ctl *Control) Del(keys ...string) *redis.IntCmd {
	return ctl.Cmd.Del(ctl.Ctx, keys...)
}

func (ctl *Control) Unlink(keys ...string) *redis.IntCmd {
	return ctl.Cmd.Unlink(ctl.Ctx, keys...)
}

func (ctl *Control) Dump(key string) *redis.StringCmd {
	return ctl.Cmd.Dump(ctl.Ctx, key)
}

func (ctl *Control) Exists(keys ...string) *redis.IntCmd {
	return ctl.Cmd.Exists(ctl.Ctx, keys...)
}

func (ctl *Control) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	return ctl.Cmd.Expire(ctl.Ctx, key, expiration)
}

func (ctl *Control) ExpireAt(key string, tm time.Time) *redis.BoolCmd {
	return ctl.Cmd.ExpireAt(ctl.Ctx, key, tm)
}

func (ctl *Control) Keys(pattern string) *redis.StringSliceCmd {
	return ctl.Cmd.Keys(ctl.Ctx, pattern)
}

func (ctl *Control) Migrate(host, port, key string, db int, timeout time.Duration) *redis.StatusCmd {
	return ctl.Cmd.Migrate(ctl.Ctx, host, port, key, db, timeout)
}

func (ctl *Control) Move(key string, db int) *redis.BoolCmd {
	return ctl.Cmd.Move(ctl.Ctx, key, db)
}

func (ctl *Control) ObjectRefCount(key string) *redis.IntCmd {
	return ctl.Cmd.ObjectRefCount(ctl.Ctx, key)
}

func (ctl *Control) ObjectEncoding(key string) *redis.StringCmd {
	return ctl.Cmd.ObjectEncoding(ctl.Ctx, key)
}

func (ctl *Control) ObjectIdleTime(key string) *redis.DurationCmd {
	return ctl.Cmd.ObjectIdleTime(ctl.Ctx, key)
}

func (ctl *Control) Persist(key string) *redis.BoolCmd {
	return ctl.Cmd.Persist(ctl.Ctx, key)
}

func (ctl *Control) PExpire(key string, expiration time.Duration) *redis.BoolCmd {
	return ctl.Cmd.PExpire(ctl.Ctx, key, expiration)
}

func (ctl *Control) PExpireAt(key string, tm time.Time) *redis.BoolCmd {
	return ctl.Cmd.PExpireAt(ctl.Ctx, key, tm)
}

func (ctl *Control) PTTL(key string) *redis.DurationCmd {
	return ctl.Cmd.PTTL(ctl.Ctx, key)
}

func (ctl *Control) RandomKey() *redis.StringCmd {
	return ctl.Cmd.RandomKey(ctl.Ctx)
}

func (ctl *Control) Rename(oldKey, newKey string) *redis.StatusCmd {
	return ctl.Cmd.Rename(ctl.Ctx, oldKey, newKey)
}

func (ctl *Control) RenameNX(oldKey, newKey string) *redis.BoolCmd {
	return ctl.Cmd.RenameNX(ctl.Ctx, oldKey, newKey)
}

func (ctl *Control) Restore(key string, ttl time.Duration, value string) *redis.StatusCmd {
	return ctl.Cmd.Restore(ctl.Ctx, key, ttl, value)
}

func (ctl *Control) RestoreReplace(key string, ttl time.Duration, value string) *redis.StatusCmd {
	return ctl.Cmd.RestoreReplace(ctl.Ctx, key, ttl, value)
}

func (ctl *Control) Sort(key string, sort *redis.Sort) *redis.StringSliceCmd {
	return ctl.Cmd.Sort(ctl.Ctx, key, sort)
}

func (ctl *Control) SortStore(key, store string, sort *redis.Sort) *redis.IntCmd {
	return ctl.Cmd.SortStore(ctl.Ctx, key, store, sort)
}

func (ctl *Control) SortInterfaces(key string, sort *redis.Sort) *redis.SliceCmd {
	return ctl.Cmd.SortInterfaces(ctl.Ctx, key, sort)
}
func (ctl *Control) ZScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return ctl.Cmd.ZScan(ctl.Ctx, key, cursor, match, count)
}

func (ctl *Control) Touch(keys ...string) *redis.IntCmd {
	return ctl.Cmd.Touch(ctl.Ctx, keys...)
}

func (ctl *Control) TTL(key string) *redis.DurationCmd {
	return ctl.Cmd.TTL(ctl.Ctx, key)
}

func (ctl *Control) Type(key string) *redis.StatusCmd {
	return ctl.Cmd.Type(ctl.Ctx, key)
}

func (ctl *Control) Scan(cursor uint64, match string, count int64) *redis.ScanCmd {
	return ctl.Cmd.Scan(ctl.Ctx, cursor, match, count)
}

func (ctl *Control) SScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return ctl.Cmd.SScan(ctl.Ctx, key, cursor, match, count)
}

func (ctl *Control) Append(key, value string) *redis.IntCmd {
	return ctl.Cmd.Append(ctl.Ctx, key, value)
}

func (ctl *Control) BitCount(key string, bitCount *redis.BitCount) *redis.IntCmd {
	return ctl.Cmd.BitCount(ctl.Ctx, key, bitCount)
}

func (ctl *Control) BitOpAnd(destKey string, keys ...string) *redis.IntCmd {
	return ctl.Cmd.BitOpAnd(ctl.Ctx, destKey, keys...)
}

func (ctl *Control) BitOpOr(destKey string, keys ...string) *redis.IntCmd {
	return ctl.Cmd.BitOpOr(ctl.Ctx, destKey, keys...)
}

func (ctl *Control) BitOpXor(destKey string, keys ...string) *redis.IntCmd {
	return ctl.Cmd.BitOpXor(ctl.Ctx, destKey, keys...)
}

func (ctl *Control) BitOpNot(destKey string, key string) *redis.IntCmd {
	return ctl.Cmd.BitOpNot(ctl.Ctx, destKey, key)
}

func (ctl *Control) BitPos(key string, bit int64, pos ...int64) *redis.IntCmd {
	return ctl.Cmd.BitPos(ctl.Ctx, key, bit, pos...)
}

func (ctl *Control) Decr(key string) *redis.IntCmd {
	return ctl.Cmd.Decr(ctl.Ctx, key)
}

func (ctl *Control) DecrBy(key string, decrement int64) *redis.IntCmd {
	return ctl.Cmd.DecrBy(ctl.Ctx, key, decrement)
}

func (ctl *Control) Get(key string) *redis.StringCmd {
	return ctl.Cmd.Get(ctl.Ctx, key)
}

func (ctl *Control) GetBit(key string, offset int64) *redis.IntCmd {
	return ctl.Cmd.GetBit(ctl.Ctx, key, offset)
}

func (ctl *Control) GetRange(key string, start, end int64) *redis.StringCmd {
	return ctl.Cmd.GetRange(ctl.Ctx, key, start, end)
}

func (ctl *Control) GetSet(key string, value interface{}) *redis.StringCmd {
	return ctl.Cmd.GetSet(ctl.Ctx, key, value)
}

func (ctl *Control) Incr(key string) *redis.IntCmd {
	return ctl.Cmd.Incr(ctl.Ctx, key)
}

func (ctl *Control) IncrBy(key string, value int64) *redis.IntCmd {
	return ctl.Cmd.IncrBy(ctl.Ctx, key, value)
}

func (ctl *Control) IncrByFloat(key string, value float64) *redis.FloatCmd {
	return ctl.Cmd.IncrByFloat(ctl.Ctx, key, value)
}

func (ctl *Control) MGet(keys ...string) *redis.SliceCmd {
	return ctl.Cmd.MGet(ctl.Ctx, keys...)
}

func (ctl *Control) MSet(pairs ...interface{}) *redis.StatusCmd {
	return ctl.Cmd.MSet(ctl.Ctx, pairs...)
}

func (ctl *Control) MSetNX(pairs ...interface{}) *redis.BoolCmd {
	return ctl.Cmd.MSetNX(ctl.Ctx, pairs...)
}

func (ctl *Control) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return ctl.Cmd.Set(ctl.Ctx, key, value, expiration)
}

func (ctl *Control) SetBit(key string, offset int64, value int) *redis.IntCmd {
	return ctl.Cmd.SetBit(ctl.Ctx, key, offset, value)
}

func (ctl *Control) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return ctl.Cmd.SetNX(ctl.Ctx, key, value, expiration)
}

func (ctl *Control) SetXX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return ctl.Cmd.SetXX(ctl.Ctx, key, value, expiration)
}

func (ctl *Control) SetRange(key string, offset int64, value string) *redis.IntCmd {
	return ctl.Cmd.SetRange(ctl.Ctx, key, offset, value)
}

func (ctl *Control) StrLen(key string) *redis.IntCmd {
	return ctl.Cmd.StrLen(ctl.Ctx, key)
}

func (ctl *Control) HDel(key string, fields ...string) *redis.IntCmd {
	return ctl.Cmd.HDel(ctl.Ctx, key, fields...)
}

func (ctl *Control) HExists(key, field string) *redis.BoolCmd {
	return ctl.Cmd.HExists(ctl.Ctx, key, field)
}

func (ctl *Control) HGet(key, field string) *redis.StringCmd {
	return ctl.Cmd.HGet(ctl.Ctx, key, field)
}

func (ctl *Control) HGetAll(key string) *redis.MapStringStringCmd {
	return ctl.Cmd.HGetAll(ctl.Ctx, key)
}

func (ctl *Control) HIncrBy(key, field string, incr int64) *redis.IntCmd {
	return ctl.Cmd.HIncrBy(ctl.Ctx, key, field, incr)
}

func (ctl *Control) HIncrByFloat(key, field string, incr float64) *redis.FloatCmd {
	return ctl.Cmd.HIncrByFloat(ctl.Ctx, key, field, incr)
}

func (ctl *Control) HKeys(key string) *redis.StringSliceCmd {
	return ctl.Cmd.HKeys(ctl.Ctx, key)
}

func (ctl *Control) HLen(key string) *redis.IntCmd {
	return ctl.Cmd.HLen(ctl.Ctx, key)
}

func (ctl *Control) HMGet(key string, fields ...string) *redis.SliceCmd {
	return ctl.Cmd.HMGet(ctl.Ctx, key, fields...)
}

func (ctl *Control) HMSet(key string, values ...interface{}) *redis.BoolCmd {
	return ctl.Cmd.HMSet(ctl.Ctx, key, values)
}

func (ctl *Control) HSet(key string, values ...interface{}) *redis.IntCmd {
	return ctl.Cmd.HSet(ctl.Ctx, key, values)
}

func (ctl *Control) HSetNX(key, field string, value interface{}) *redis.BoolCmd {
	return ctl.Cmd.HSetNX(ctl.Ctx, key, field, value)
}

func (ctl *Control) HVals(key string) *redis.StringSliceCmd {
	return ctl.Cmd.HVals(ctl.Ctx, key)
}

func (ctl *Control) BLPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return ctl.Cmd.BLPop(ctl.Ctx, timeout, keys...)
}

func (ctl *Control) BRPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return ctl.Cmd.BRPop(ctl.Ctx, timeout, keys...)
}

func (ctl *Control) BRPopLPush(source, destination string, timeout time.Duration) *redis.StringCmd {
	return ctl.Cmd.BRPopLPush(ctl.Ctx, source, destination, timeout)
}

func (ctl *Control) LIndex(key string, index int64) *redis.StringCmd {
	return ctl.Cmd.LIndex(ctl.Ctx, key, index)
}

func (ctl *Control) LInsert(key, op string, pivot, value interface{}) *redis.IntCmd {
	return ctl.Cmd.LInsert(ctl.Ctx, key, op, pivot, value)
}

func (ctl *Control) LInsertBefore(key string, pivot, value interface{}) *redis.IntCmd {
	return ctl.Cmd.LInsertBefore(ctl.Ctx, key, pivot, value)
}

func (ctl *Control) LInsertAfter(key string, pivot, value interface{}) *redis.IntCmd {
	return ctl.Cmd.LInsertAfter(ctl.Ctx, key, pivot, value)
}

func (ctl *Control) LLen(key string) *redis.IntCmd {
	return ctl.Cmd.LLen(ctl.Ctx, key)
}

func (ctl *Control) LPop(key string) *redis.StringCmd {
	return ctl.Cmd.LPop(ctl.Ctx, key)
}

func (ctl *Control) LPush(key string, values ...interface{}) *redis.IntCmd {
	return ctl.Cmd.LPush(ctl.Ctx, key, values...)
}

func (ctl *Control) LPushX(key string, value interface{}) *redis.IntCmd {
	return ctl.Cmd.LPushX(ctl.Ctx, key, value)
}

func (ctl *Control) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	return ctl.Cmd.LRange(ctl.Ctx, key, start, stop)
}

func (ctl *Control) LRem(key string, count int64, value interface{}) *redis.IntCmd {
	return ctl.Cmd.LRem(ctl.Ctx, key, count, value)
}

func (ctl *Control) LSet(key string, index int64, value interface{}) *redis.StatusCmd {
	return ctl.Cmd.LSet(ctl.Ctx, key, index, value)
}

func (ctl *Control) LTrim(key string, start, stop int64) *redis.StatusCmd {
	return ctl.Cmd.LTrim(ctl.Ctx, key, start, stop)
}

func (ctl *Control) RPop(key string) *redis.StringCmd {
	return ctl.Cmd.RPop(ctl.Ctx, key)
}

func (ctl *Control) RPopLPush(source, destination string) *redis.StringCmd {
	return ctl.Cmd.RPopLPush(ctl.Ctx, source, destination)
}

func (ctl *Control) RPush(key string, values ...interface{}) *redis.IntCmd {
	return ctl.Cmd.RPush(ctl.Ctx, key, values...)
}

func (ctl *Control) RPushX(key string, value interface{}) *redis.IntCmd {
	return ctl.Cmd.RPushX(ctl.Ctx, key, value)
}

func (ctl *Control) SAdd(key string, members ...interface{}) *redis.IntCmd {
	return ctl.Cmd.SAdd(ctl.Ctx, key, members...)
}

func (ctl *Control) SCard(key string) *redis.IntCmd {
	return ctl.Cmd.SCard(ctl.Ctx, key)
}

func (ctl *Control) SDiff(keys ...string) *redis.StringSliceCmd {
	return ctl.Cmd.SDiff(ctl.Ctx, keys...)
}

func (ctl *Control) SDiffStore(destination string, keys ...string) *redis.IntCmd {
	return ctl.Cmd.SDiffStore(ctl.Ctx, destination, keys...)
}

func (ctl *Control) SInter(keys ...string) *redis.StringSliceCmd {
	return ctl.Cmd.SInter(ctl.Ctx, keys...)
}

func (ctl *Control) SInterStore(destination string, keys ...string) *redis.IntCmd {
	return ctl.Cmd.SInterStore(ctl.Ctx, destination, keys...)
}

func (ctl *Control) SIsMember(key string, member interface{}) *redis.BoolCmd {
	return ctl.Cmd.SIsMember(ctl.Ctx, key, member)
}

func (ctl *Control) SMembers(key string) *redis.StringSliceCmd {
	return ctl.Cmd.SMembers(ctl.Ctx, key)
}

func (ctl *Control) SMembersMap(key string) *redis.StringStructMapCmd {
	return ctl.Cmd.SMembersMap(ctl.Ctx, key)
}

func (ctl *Control) SMove(source, destination string, member interface{}) *redis.BoolCmd {
	return ctl.Cmd.SMove(ctl.Ctx, source, destination, member)
}

func (ctl *Control) SPop(key string) *redis.StringCmd {
	return ctl.Cmd.SPop(ctl.Ctx, key)
}

func (ctl *Control) SPopN(key string, count int64) *redis.StringSliceCmd {
	return ctl.Cmd.SPopN(ctl.Ctx, key, count)
}

func (ctl *Control) SRandMember(key string) *redis.StringCmd {
	return ctl.Cmd.SRandMember(ctl.Ctx, key)
}

func (ctl *Control) SRandMemberN(key string, count int64) *redis.StringSliceCmd {
	return ctl.Cmd.SRandMemberN(ctl.Ctx, key, count)
}

func (ctl *Control) SRem(key string, members ...interface{}) *redis.IntCmd {
	return ctl.Cmd.SRem(ctl.Ctx, key, members...)
}

func (ctl *Control) SUnion(keys ...string) *redis.StringSliceCmd {
	return ctl.Cmd.SUnion(ctl.Ctx, keys...)
}

func (ctl *Control) SUnionStore(destination string, keys ...string) *redis.IntCmd {
	return ctl.Cmd.SUnionStore(ctl.Ctx, destination, keys...)
}

func (ctl *Control) XAdd(a *redis.XAddArgs) *redis.StringCmd {
	return ctl.Cmd.XAdd(ctl.Ctx, a)
}

func (ctl *Control) XDel(stream string, ids ...string) *redis.IntCmd {
	return ctl.Cmd.XDel(ctl.Ctx, stream, ids...)
}

func (ctl *Control) XLen(stream string) *redis.IntCmd {
	return ctl.Cmd.XLen(ctl.Ctx, stream)
}

func (ctl *Control) XRange(stream, start, stop string) *redis.XMessageSliceCmd {
	return ctl.Cmd.XRange(ctl.Ctx, stream, start, stop)
}

func (ctl *Control) XRangeN(stream, start, stop string, count int64) *redis.XMessageSliceCmd {
	return ctl.Cmd.XRangeN(ctl.Ctx, stream, start, stop, count)
}

func (ctl *Control) XRevRange(stream string, start, stop string) *redis.XMessageSliceCmd {
	return ctl.Cmd.XRevRange(ctl.Ctx, stream, start, stop)
}

func (ctl *Control) XRevRangeN(stream string, start, stop string, count int64) *redis.XMessageSliceCmd {
	return ctl.Cmd.XRevRangeN(ctl.Ctx, stream, start, stop, count)
}

func (ctl *Control) XRead(a *redis.XReadArgs) *redis.XStreamSliceCmd {
	return ctl.Cmd.XRead(ctl.Ctx, a)
}

func (ctl *Control) XReadStreams(streams ...string) *redis.XStreamSliceCmd {
	return ctl.Cmd.XReadStreams(ctl.Ctx, streams...)
}

func (ctl *Control) XGroupCreate(stream, group, start string) *redis.StatusCmd {
	return ctl.Cmd.XGroupCreate(ctl.Ctx, stream, group, start)
}

func (ctl *Control) XGroupCreateMkStream(stream, group, start string) *redis.StatusCmd {
	return ctl.Cmd.XGroupCreateMkStream(ctl.Ctx, stream, group, start)
}

func (ctl *Control) XGroupSetID(stream, group, start string) *redis.StatusCmd {
	return ctl.Cmd.XGroupSetID(ctl.Ctx, stream, group, start)
}

func (ctl *Control) XGroupDestroy(stream, group string) *redis.IntCmd {
	return ctl.Cmd.XGroupDestroy(ctl.Ctx, stream, group)
}

func (ctl *Control) XGroupDelConsumer(stream, group, consumer string) *redis.IntCmd {
	return ctl.Cmd.XGroupDelConsumer(ctl.Ctx, stream, group, consumer)
}

func (ctl *Control) XReadGroup(a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	return ctl.Cmd.XReadGroup(ctl.Ctx, a)
}

func (ctl *Control) XAck(stream, group string, ids ...string) *redis.IntCmd {
	return ctl.Cmd.XAck(ctl.Ctx, stream, group, ids...)
}

func (ctl *Control) XPending(stream, group string) *redis.XPendingCmd {
	return ctl.Cmd.XPending(ctl.Ctx, stream, group)
}

func (ctl *Control) XPendingExt(a *redis.XPendingExtArgs) *redis.XPendingExtCmd {
	return ctl.Cmd.XPendingExt(ctl.Ctx, a)
}

func (ctl *Control) XClaim(a *redis.XClaimArgs) *redis.XMessageSliceCmd {
	return ctl.Cmd.XClaim(ctl.Ctx, a)
}

func (ctl *Control) XClaimJustID(a *redis.XClaimArgs) *redis.StringSliceCmd {
	return ctl.Cmd.XClaimJustID(ctl.Ctx, a)
}

func (ctl *Control) BZPopMax(timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	return ctl.Cmd.BZPopMax(ctl.Ctx, timeout, keys...)
}

func (ctl *Control) BZPopMin(timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	return ctl.Cmd.BZPopMin(ctl.Ctx, timeout, keys...)
}

func (ctl *Control) ZAdd(key string, members ...redis.Z) *redis.IntCmd {
	return ctl.Cmd.ZAdd(ctl.Ctx, key, members...)
}

func (ctl *Control) ZAddNX(key string, members ...redis.Z) *redis.IntCmd {
	return ctl.Cmd.ZAddNX(ctl.Ctx, key, members...)
}

func (ctl *Control) ZAddXX(key string, members ...redis.Z) *redis.IntCmd {
	return ctl.Cmd.ZAddXX(ctl.Ctx, key, members...)
}

func (ctl *Control) ZCard(key string) *redis.IntCmd {
	return ctl.Cmd.ZCard(ctl.Ctx, key)
}

func (ctl *Control) ZCount(key, min, max string) *redis.IntCmd {
	return ctl.Cmd.ZCount(ctl.Ctx, key, min, max)
}

func (ctl *Control) ZLexCount(key, min, max string) *redis.IntCmd {
	return ctl.Cmd.ZLexCount(ctl.Ctx, key, min, max)
}

func (ctl *Control) ZIncrBy(key string, increment float64, member string) *redis.FloatCmd {
	return ctl.Cmd.ZIncrBy(ctl.Ctx, key, increment, member)
}

func (ctl *Control) ZInterStore(destination string, store *redis.ZStore) *redis.IntCmd {
	return ctl.Cmd.ZInterStore(ctl.Ctx, destination, store)
}

func (ctl *Control) ZPopMax(key string, count ...int64) *redis.ZSliceCmd {
	return ctl.Cmd.ZPopMax(ctl.Ctx, key, count...)
}

func (ctl *Control) ZPopMin(key string, count ...int64) *redis.ZSliceCmd {
	return ctl.Cmd.ZPopMin(ctl.Ctx, key, count...)
}

func (ctl *Control) ZRange(key string, start, stop int64) *redis.StringSliceCmd {
	return ctl.Cmd.ZRange(ctl.Ctx, key, start, stop)
}

func (ctl *Control) ZRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	return ctl.Cmd.ZRangeWithScores(ctl.Ctx, key, start, stop)
}

func (ctl *Control) ZRangeByScore(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return ctl.Cmd.ZRangeByScore(ctl.Ctx, key, opt)
}

func (ctl *Control) ZRangeByLex(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return ctl.Cmd.ZRangeByLex(ctl.Ctx, key, opt)
}

func (ctl *Control) ZRangeByScoreWithScores(key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return ctl.Cmd.ZRangeByScoreWithScores(ctl.Ctx, key, opt)
}

func (ctl *Control) ZRank(key, member string) *redis.IntCmd {
	return ctl.Cmd.ZRank(ctl.Ctx, key, member)
}

func (ctl *Control) ZRem(key string, members ...interface{}) *redis.IntCmd {
	return ctl.Cmd.ZRem(ctl.Ctx, key, members...)
}

func (ctl *Control) ZRemRangeByRank(key string, start, stop int64) *redis.IntCmd {
	return ctl.Cmd.ZRemRangeByRank(ctl.Ctx, key, start, stop)
}

func (ctl *Control) ZRemRangeByScore(key, min, max string) *redis.IntCmd {
	return ctl.Cmd.ZRemRangeByScore(ctl.Ctx, key, min, max)
}

func (ctl *Control) ZRemRangeByLex(key, min, max string) *redis.IntCmd {
	return ctl.Cmd.ZRemRangeByLex(ctl.Ctx, key, min, max)
}

func (ctl *Control) ZRevRange(key string, start, stop int64) *redis.StringSliceCmd {
	return ctl.Cmd.ZRevRange(ctl.Ctx, key, start, stop)
}

func (ctl *Control) ZRevRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	return ctl.Cmd.ZRevRangeWithScores(ctl.Ctx, key, start, stop)
}

func (ctl *Control) ZRevRangeByScore(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return ctl.Cmd.ZRevRangeByScore(ctl.Ctx, key, opt)
}

func (ctl *Control) ZRevRangeByLex(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return ctl.Cmd.ZRevRangeByLex(ctl.Ctx, key, opt)
}

func (ctl *Control) ZRevRangeByScoreWithScores(key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return ctl.Cmd.ZRevRangeByScoreWithScores(ctl.Ctx, key, opt)
}

func (ctl *Control) ZRevRank(key, member string) *redis.IntCmd {
	return ctl.Cmd.ZRevRank(ctl.Ctx, key, member)
}

func (ctl *Control) ZScore(key, member string) *redis.FloatCmd {
	return ctl.Cmd.ZScore(ctl.Ctx, key, member)
}

func (ctl *Control) ZUnionStore(dest string, store *redis.ZStore) *redis.IntCmd {
	return ctl.Cmd.ZUnionStore(ctl.Ctx, dest, store)
}

func (ctl *Control) PFAdd(key string, els ...interface{}) *redis.IntCmd {
	return ctl.Cmd.PFAdd(ctl.Ctx, key, els...)
}

func (ctl *Control) PFCount(keys ...string) *redis.IntCmd {
	return ctl.Cmd.PFCount(ctl.Ctx, keys...)
}

func (ctl *Control) PFMerge(dest string, keys ...string) *redis.StatusCmd {
	return ctl.Cmd.PFMerge(ctl.Ctx, dest, keys...)
}

func (ctl *Control) BgRewriteAOF() *redis.StatusCmd {
	return ctl.Cmd.BgRewriteAOF(ctl.Ctx)
}

func (ctl *Control) BgSave() *redis.StatusCmd {
	return ctl.Cmd.BgSave(ctl.Ctx)
}

func (ctl *Control) ClientKill(ipPort string) *redis.StatusCmd {
	return ctl.Cmd.ClientKill(ctl.Ctx, ipPort)
}

func (ctl *Control) ClientKillByFilter(keys ...string) *redis.IntCmd {
	return ctl.Cmd.ClientKillByFilter(ctl.Ctx, keys...)
}

func (ctl *Control) ClientList() *redis.StringCmd {
	return ctl.Cmd.ClientList(ctl.Ctx)
}

func (ctl *Control) ClientPause(dur time.Duration) *redis.BoolCmd {
	return ctl.Cmd.ClientPause(ctl.Ctx, dur)
}

func (ctl *Control) ClientID() *redis.IntCmd {
	return ctl.Cmd.ClientID(ctl.Ctx)
}

func (ctl *Control) ConfigGet(parameter string) *redis.MapStringStringCmd {
	return ctl.Cmd.ConfigGet(ctl.Ctx, parameter)
}

func (ctl *Control) ConfigResetStat() *redis.StatusCmd {
	return ctl.Cmd.ConfigResetStat(ctl.Ctx)
}

func (ctl *Control) ConfigSet(parameter, value string) *redis.StatusCmd {
	return ctl.Cmd.ConfigSet(ctl.Ctx, parameter, value)
}

func (ctl *Control) ConfigRewrite() *redis.StatusCmd {
	return ctl.Cmd.ConfigRewrite(ctl.Ctx)
}

func (ctl *Control) DBSize() *redis.IntCmd {
	return ctl.Cmd.DBSize(ctl.Ctx)
}

func (ctl *Control) FlushAll() *redis.StatusCmd {
	return ctl.Cmd.FlushAll(ctl.Ctx)
}

func (ctl *Control) FlushAllAsync() *redis.StatusCmd {
	return ctl.Cmd.FlushAllAsync(ctl.Ctx)
}

func (ctl *Control) FlushDB() *redis.StatusCmd {
	return ctl.Cmd.FlushDB(ctl.Ctx)
}

func (ctl *Control) FlushDBAsync() *redis.StatusCmd {
	return ctl.Cmd.FlushDBAsync(ctl.Ctx)
}

func (ctl *Control) Info(section ...string) *redis.StringCmd {
	return ctl.Cmd.Info(ctl.Ctx, section...)
}

func (ctl *Control) LastSave() *redis.IntCmd {
	return ctl.Cmd.LastSave(ctl.Ctx)
}

func (ctl *Control) Save() *redis.StatusCmd {
	return ctl.Cmd.Save(ctl.Ctx)
}

func (ctl *Control) Shutdown() *redis.StatusCmd {
	return ctl.Cmd.Shutdown(ctl.Ctx)
}

func (ctl *Control) ShutdownSave() *redis.StatusCmd {
	return ctl.Cmd.ShutdownSave(ctl.Ctx)
}

func (ctl *Control) ShutdownNoSave() *redis.StatusCmd {
	return ctl.Cmd.ShutdownNoSave(ctl.Ctx)
}

func (ctl *Control) SlaveOf(host, port string) *redis.StatusCmd {
	return ctl.Cmd.SlaveOf(ctl.Ctx, host, port)
}

func (ctl *Control) Time() *redis.TimeCmd {
	return ctl.Cmd.Time(ctl.Ctx)
}

func (ctl *Control) Eval(script string, keys []string, args ...interface{}) *redis.Cmd {
	return ctl.Cmd.Eval(ctl.Ctx, script, keys, args...)
}

func (ctl *Control) EvalSha(sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return ctl.Cmd.EvalSha(ctl.Ctx, sha1, keys, args...)
}

func (ctl *Control) ScriptExists(hashes ...string) *redis.BoolSliceCmd {
	return ctl.Cmd.ScriptExists(ctl.Ctx, hashes...)
}

func (ctl *Control) ScriptFlush() *redis.StatusCmd {
	return ctl.Cmd.ScriptFlush(ctl.Ctx)
}

func (ctl *Control) ScriptKill() *redis.StatusCmd {
	return ctl.Cmd.ScriptKill(ctl.Ctx)
}

func (ctl *Control) ScriptLoad(script string) *redis.StringCmd {
	return ctl.Cmd.ScriptLoad(ctl.Ctx, script)
}

func (ctl *Control) DebugObject(key string) *redis.StringCmd {
	return ctl.Cmd.DebugObject(ctl.Ctx, key)
}

func (ctl *Control) Publish(channel string, message interface{}) *redis.IntCmd {
	return ctl.Cmd.Publish(ctl.Ctx, channel, message)
}

func (ctl *Control) PubSubChannels(pattern string) *redis.StringSliceCmd {
	return ctl.Cmd.PubSubChannels(ctl.Ctx, pattern)
}

func (ctl *Control) PubSubNumSub(channels ...string) *redis.MapStringIntCmd {
	return ctl.Cmd.PubSubNumSub(ctl.Ctx, channels...)
}

func (ctl *Control) PubSubNumPat() *redis.IntCmd {
	return ctl.Cmd.PubSubNumPat(ctl.Ctx)
}

func (ctl *Control) ClusterSlots() *redis.ClusterSlotsCmd {
	return ctl.Cmd.ClusterSlots(ctl.Ctx)
}

func (ctl *Control) ClusterNodes() *redis.StringCmd {
	return ctl.Cmd.ClusterNodes(ctl.Ctx)
}

func (ctl *Control) ClusterMeet(host, port string) *redis.StatusCmd {
	return ctl.Cmd.ClusterMeet(ctl.Ctx, host, port)
}

func (ctl *Control) ClusterForget(nodeID string) *redis.StatusCmd {
	return ctl.Cmd.ClusterForget(ctl.Ctx, nodeID)
}

func (ctl *Control) ClusterReplicate(nodeID string) *redis.StatusCmd {
	return ctl.Cmd.ClusterReplicate(ctl.Ctx, nodeID)
}

func (ctl *Control) ClusterResetSoft() *redis.StatusCmd {
	return ctl.Cmd.ClusterResetSoft(ctl.Ctx)
}

func (ctl *Control) ClusterResetHard() *redis.StatusCmd {
	return ctl.Cmd.ClusterResetHard(ctl.Ctx)
}

func (ctl *Control) ClusterInfo() *redis.StringCmd {
	return ctl.Cmd.ClusterInfo(ctl.Ctx)
}

func (ctl *Control) ClusterKeySlot(key string) *redis.IntCmd {
	return ctl.Cmd.ClusterKeySlot(ctl.Ctx, key)
}

func (ctl *Control) ClusterGetKeysInSlot(slot int, count int) *redis.StringSliceCmd {
	return ctl.Cmd.ClusterGetKeysInSlot(ctl.Ctx, slot, count)
}

func (ctl *Control) ClusterCountFailureReports(nodeID string) *redis.IntCmd {
	return ctl.Cmd.ClusterCountFailureReports(ctl.Ctx, nodeID)
}

func (ctl *Control) ClusterCountKeysInSlot(slot int) *redis.IntCmd {
	return ctl.Cmd.ClusterCountKeysInSlot(ctl.Ctx, slot)
}

func (ctl *Control) ClusterDelSlots(slots ...int) *redis.StatusCmd {
	return ctl.Cmd.ClusterDelSlots(ctl.Ctx, slots...)
}

func (ctl *Control) ClusterDelSlotsRange(min, max int) *redis.StatusCmd {
	return ctl.Cmd.ClusterDelSlotsRange(ctl.Ctx, min, max)
}

func (ctl *Control) ClusterSaveConfig() *redis.StatusCmd {
	return ctl.Cmd.ClusterSaveConfig(ctl.Ctx)
}

func (ctl *Control) ClusterSlaves(nodeID string) *redis.StringSliceCmd {
	return ctl.Cmd.ClusterSlaves(ctl.Ctx, nodeID)
}

func (ctl *Control) ClusterFailover() *redis.StatusCmd {
	return ctl.Cmd.ClusterFailover(ctl.Ctx)
}

func (ctl *Control) ClusterAddSlots(slots ...int) *redis.StatusCmd {
	return ctl.Cmd.ClusterAddSlots(ctl.Ctx, slots...)
}

func (ctl *Control) ClusterAddSlotsRange(min, max int) *redis.StatusCmd {
	return ctl.Cmd.ClusterAddSlotsRange(ctl.Ctx, min, max)
}

func (ctl *Control) GeoAdd(key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	return ctl.Cmd.GeoAdd(ctl.Ctx, key, geoLocation...)
}

func (ctl *Control) GeoPos(key string, members ...string) *redis.GeoPosCmd {
	return ctl.Cmd.GeoPos(ctl.Ctx, key, members...)
}

func (ctl *Control) GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return ctl.Cmd.GeoRadius(ctl.Ctx, key, longitude, latitude, query)
}

func (ctl *Control) GeoRadiusStore(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.IntCmd {
	return ctl.Cmd.GeoRadiusStore(ctl.Ctx, key, longitude, latitude, query)
}

func (ctl *Control) GeoRadiusByMember(key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return ctl.Cmd.GeoRadiusByMember(ctl.Ctx, key, member, query)
}

func (ctl *Control) GeoRadiusByMemberRO(key, member string, query *redis.GeoRadiusQuery) *redis.IntCmd {
	return ctl.Cmd.GeoRadiusByMemberStore(ctl.Ctx, key, member, query)
}

func (ctl *Control) GeoDist(key string, member1, member2, unit string) *redis.FloatCmd {
	return ctl.Cmd.GeoDist(ctl.Ctx, key, member1, member2, unit)
}

func (ctl *Control) GeoHash(key string, members ...string) *redis.StringSliceCmd {
	return ctl.Cmd.GeoHash(ctl.Ctx, key, members...)
}

func (ctl *Control) ReadOnly() *redis.StatusCmd {
	return ctl.Cmd.ReadOnly(ctl.Ctx)
}

func (ctl *Control) ReadWrite() *redis.StatusCmd {
	return ctl.Cmd.ReadWrite(ctl.Ctx)
}

func (ctl *Control) MemoryUsage(key string, samples ...int) *redis.IntCmd {
	return ctl.Cmd.MemoryUsage(ctl.Ctx, key, samples...)
}
