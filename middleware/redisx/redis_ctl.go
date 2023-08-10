package redisx

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var CTL *Controller

// redis控制器
type Controller struct {
	Config *Config       // redis配置
	Cmd    redis.Cmdable // redis接口方法
}

// 初始化redis控制器
func Init(conf *Config) {
	var cmd = conf.NewRedisCmdable()
	if ok, err := Ping(cmd); ok && err == nil {
		CTL = &Controller{Config: conf, Cmd: cmd}
		log.Error("Redis连接成功！", conf.Format())
	} else {
		log.Error("Redis连接失败！", conf.Format())
		log.Error("error : ", err)
	}
}

func Ping(cmd redis.Cmdable) (bool, error) {
	_, err := cmd.Ping(context.Background()).Result()
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (ctl *Controller) Pipeline() redis.Pipeliner {
	return ctl.Cmd.Pipeline()
}

func (ctl *Controller) Command(ctx context.Context) *redis.CommandsInfoCmd {
	return ctl.Cmd.Command(ctx)
}

func (ctl *Controller) ClientGetName(ctx context.Context) *redis.StringCmd {
	return ctl.Cmd.ClientGetName(ctx)
}

func (ctl *Controller) Echo(ctx context.Context, message interface{}) *redis.StringCmd {
	return ctl.Cmd.Echo(ctx, message)
}

func (ctl *Controller) Quit(ctx context.Context) *redis.StatusCmd {
	return ctl.Cmd.Quit(ctx)
}

func (ctl *Controller) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return ctl.Cmd.Del(ctx, keys...)
}

func (ctl *Controller) Dump(ctx context.Context, key string) *redis.StringCmd {
	return ctl.Cmd.Dump(ctx, key)
}

func (ctl *Controller) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return ctl.Cmd.Exists(ctx, keys...)
}

func (ctl *Controller) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return ctl.Cmd.Expire(ctx, key, expiration)
}

func (ctl *Controller) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	return ctl.Cmd.ExpireAt(ctx, key, tm)
}

func (ctl *Controller) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return ctl.Cmd.Keys(ctx, pattern)
}

func (ctl *Controller) Move(ctx context.Context, key string, db int) *redis.BoolCmd {
	return ctl.Cmd.Move(ctx, key, db)
}

func (ctl *Controller) Sort(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	return ctl.Cmd.Sort(ctx, key, sort)
}

func (ctl *Controller) SortStore(ctx context.Context, key, store string, sort *redis.Sort) *redis.IntCmd {
	return ctl.Cmd.SortStore(ctx, key, store, sort)
}

func (ctl *Controller) TTL(ctx context.Context, key string) *redis.DurationCmd {
	return ctl.Cmd.TTL(ctx, key)
}

func (ctl *Controller) Decr(ctx context.Context, key string) *redis.IntCmd {
	return ctl.Cmd.Decr(ctx, key)
}

func (ctl *Controller) DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd {
	return ctl.Cmd.DecrBy(ctx, key, decrement)
}

func (ctl *Controller) Incr(ctx context.Context, key string) *redis.IntCmd {
	return ctl.Cmd.Incr(ctx, key)
}

func (ctl *Controller) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	return ctl.Cmd.IncrBy(ctx, key, value)
}

func (ctl *Controller) Get(ctx context.Context, key string) *redis.StringCmd {
	return ctl.Cmd.Get(ctx, key)
}

func (ctl *Controller) GetRange(ctx context.Context, key string, start, end int64) *redis.StringCmd {
	return ctl.Cmd.GetRange(ctx, key, start, end)
}

func (ctl *Controller) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	return ctl.Cmd.MGet(ctx, keys...)
}

func (ctl *Controller) MSet(ctx context.Context, pairs ...interface{}) *redis.StatusCmd {
	return ctl.Cmd.MSet(ctx, pairs...)
}

func (ctl *Controller) MSetNX(ctx context.Context, pairs ...interface{}) *redis.BoolCmd {
	return ctl.Cmd.MSetNX(ctx, pairs...)
}

func (ctl *Controller) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return ctl.Cmd.Set(ctx, key, value, expiration)
}

func (ctl *Controller) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return ctl.Cmd.SetNX(ctx, key, value, expiration)
}

func (ctl *Controller) SetRange(ctx context.Context, key string, offset int64, value string) *redis.IntCmd {
	return ctl.Cmd.SetRange(ctx, key, offset, value)
}

func (ctl *Controller) StrLen(ctx context.Context, key string) *redis.IntCmd {
	return ctl.Cmd.StrLen(ctx, key)
}

func (ctl *Controller) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	return ctl.Cmd.HDel(ctx, key, fields...)
}

func (ctl *Controller) HExists(ctx context.Context, key, field string) *redis.BoolCmd {
	return ctl.Cmd.HExists(ctx, key, field)
}

func (ctl *Controller) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return ctl.Cmd.HGet(ctx, key, field)
}

func (ctl *Controller) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	return ctl.Cmd.HGetAll(ctx, key)
}

func (ctl *Controller) HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd {
	return ctl.Cmd.HIncrBy(ctx, key, field, incr)
}

func (ctl *Controller) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	return ctl.Cmd.HKeys(ctx, key)
}

func (ctl *Controller) HLen(ctx context.Context, key string) *redis.IntCmd {
	return ctl.Cmd.HLen(ctx, key)
}

func (ctl *Controller) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	return ctl.Cmd.HMGet(ctx, key, fields...)
}

func (ctl *Controller) HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {
	return ctl.Cmd.HMSet(ctx, key, values)
}

func (ctl *Controller) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return ctl.Cmd.HSet(ctx, key, values)
}

func (ctl *Controller) HSetNX(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {
	return ctl.Cmd.HSetNX(ctx, key, field, value)
}

func (ctl *Controller) ClientList(ctx context.Context) *redis.StringCmd {
	return ctl.Cmd.ClientList(ctx)
}

func (ctl *Controller) ClientID(ctx context.Context) *redis.IntCmd {
	return ctl.Cmd.ClientID(ctx)
}
