package redis

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/util"

	"github.com/Scut-Corgis/tiny-tiktok/config"
	"github.com/go-redis/redis/v8"
)

var mutex sync.Mutex

// RedisDb 定义一个全局变量
var RedisDb *redis.Client

var Ctx = context.Background()

func InitRedis() {
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     config.Redis_addr_port,
		Password: config.Redis_password,
		DB:       0,
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
	ret, err := RedisDb.SetNX(Ctx, key, value, time.Second*5).Result()
	if err != nil {
		log.Println("Lock error:", err.Error())
		return ret
	}
	return ret
}

func Unlock(key string) bool {
	err := RedisDb.Del(Ctx, key).Err()
	if err != nil {
		log.Println("Unlock error:", err.Error())
		return false
	}
	return true
}

func RandomTime() time.Duration {
	rand.Seed(time.Now().Unix())
	return time.Duration(rand.Int63n(25))*time.Hour + util.Day // 设置1-2天的随机过期时间
}
