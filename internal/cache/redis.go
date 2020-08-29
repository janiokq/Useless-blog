package cache

import (
	"context"
	"github.com/janiokq/Useless-blog/cinit"
	"github.com/janiokq/Useless-blog/internal/utils"
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

func CacheSetbydefaultexpiration(ctx context.Context, prefix string, id int64, data interface{}) {
	CacheSet(ctx, prefix, id, data, KeyMaxExpire)
}

func CacheGet(ctx context.Context, prefix string, id int64) (map[string]string, error) {
	k := getIdKey(prefix, id)
	r, err := cinit.RedisCli.HGetAll(k).Result()
	if err != nil {
		log.Info(err.Error(), ctx)
	}
	return r, err
}

func CacheDel(ctx context.Context, prefix string, id int64) {
	k := getIdKey(prefix, id)
	err := cinit.RedisCli.Del(k).Err()
	if err != nil {
		log.Info(err.Error(), ctx)
	}
}

func CacheSet(ctx context.Context, prefix string, id int64, data interface{}, maxExpire int) {
	_d := utils.Strcut2Map(data)
	_k := getIdKey(prefix, id)
	err := cinit.RedisCli.HMSet(_k, _d).Err()
	if err != nil {
		logx.Error(err.Error(), ctx)
		return
	}
	setKeyExpire(ctx, maxExpire)
}

func getIdKey(prefix string, ids ...int64) string {
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
