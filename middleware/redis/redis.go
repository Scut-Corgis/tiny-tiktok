package redis

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/config"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var mutex sync.Mutex

// RedisDb 定义一个全局变量
var RedisDb *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     config.Redis_addr_port,
		Password: config.Redis_password,
		DB:       0, // redis一共16个库，指定其中一个库即可
	})
	_, err := RedisDb.Ping(Ctx).Result()
	if err != nil {
		log.Panicln("err:", err.Error())
	}
	log.Println("redis has connected!")
}

func Lock(key string, value string) bool {
	mutex.Lock() // 保证程序不存在并发冲突问题
	defer mutex.Unlock()
	ret, err := RedisDb.SetNX(ctx, key, value, time.Second*10).Result()
	if err != nil {
		log.Println("Lock error:", err.Error())
		return ret
	}
	return ret
}

func Unlock(key string) bool {
	err := RedisDb.Del(ctx, key).Err()
	if err != nil {
		log.Println("Unlock error:", err.Error())
		return false
	}
	return true
}
