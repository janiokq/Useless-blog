package cache

import (
	"context"
	"github.com/janiokq/Useless-blog/cinit"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	"github.com/prometheus/common/log"
	"math/rand"
	"strconv"
	"time"
)

const (
	KeyMaxExpire     = 500 // 秒
	AgainGetStopTime = 100 * time.Millisecond
)

//func CacheSetbydefaultexpiration(ctx context.Context, prefix string, id int64, data interface{}) {
//	CacheSetHM(ctx, prefix, id, &data, KeyMaxExpire)
//}

func CacheGetHM(ctx context.Context, prefix string, id int64) (map[string]string, error) {
	k := GetIdKey(prefix, id)
	r, err := cinit.RedisCli.HGetAll(k).Result()
	if err != nil {
		log.Info(err.Error(), ctx)
	}
	return r, err
}

func CacheDel(ctx context.Context, prefix string, id int64) {
	k := GetIdKey(prefix, id)
	err := cinit.RedisCli.Del(k).Err()
	if err != nil {
		log.Info(err.Error(), ctx)
	}
}

func CacheGet(ctx context.Context, prefix string, id int64) (string, error) {
	_k := GetIdKey(prefix, id)
	return CacheGetBuyKey(ctx, _k)
}

func CacheGetBuyKey(ctx context.Context, key string) (string, error) {
	results := cinit.RedisCli.Get(key)
	if results.Err() != nil {
		logx.Error(results.Err().Error(), ctx)
		return "", results.Err()
	}
	return results.Val(), nil
}

func CacheSet(ctx context.Context, prefix string, id int64, data interface{}, maxExpire int) {
	_k := GetIdKey(prefix, id)
	err := cinit.RedisCli.Set(_k, data, time.Second*time.Duration(rand.Intn(maxExpire))).Err()
	if err != nil {
		logx.Error(err.Error(), ctx)
		return
	}
}

func CacheSetHM(ctx context.Context, prefix string, id int64, data map[string]interface{}, maxExpire int) {
	_k := GetIdKey(prefix, id)
	err := cinit.RedisCli.HMSet(_k, data).Err()
	if err != nil {
		logx.Error(err.Error(), ctx)
		return
	}
	setKeyExpire(ctx, maxExpire)
}

func GetIdKey(prefix string, ids ...int64) string {
	var s = prefix
	for _, id := range ids {
		s += "_" + strconv.FormatInt(id, 10)
	}
	return s
}

func setKeyExpire(ctx context.Context, maxExpire int, ks ...string) {
	rand.Seed(time.Now().UnixNano())
	t := time.Second * time.Duration(rand.Intn(maxExpire))
	for _, key := range ks {
		err := cinit.RedisCli.Expire(key, t).Err()
		if err != nil {
			logx.Error(err.Error(), ctx)
		}
	}
}
