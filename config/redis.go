package config

import (
	"github.com/go-redis/redis"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
)

var RedisCli *redis.Client

func redisInit() {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Addr,
		Password: Config.Redis.Password,
		DB:       Config.Redis.Db,
	})

	_, err := RedisCli.Ping().Result()
	if err != nil {
		logx.Fatal(err.Error())
	}
}

func redisClose() {
	err := RedisCli.Close()
	if err != nil {
		logx.Error(err.Error())
	}
}
