package redis

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/config"
	"github.com/go-redis/redis/v8"
)

var mutex sync.Mutex

// RedisDb 定义一个全局变量
var RedisDb *redis.Client
var RedisDbCommentIdVideoId *redis.Client // key:comment_id value:video_id relation 1:1
var RedisDbVideoIdCommentId *redis.Client // key:video_id value:comment_id ralation 1:n
var Ctx = context.Background()

func InitRedis() {
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     config.Redis_addr_port,
		Password: config.Redis_password,
		DB:       0, // redis一共16个库，指定其中一个库即可
	})
	// 将key:comment_id value:video_id存入DB1
	RedisDbCommentIdVideoId = redis.NewClient(&redis.Options{
		Addr:     config.Redis_addr_port,
		Password: config.Redis_password,
		DB:       1,
	})
	// 将key:video_id value:comment_id存入DB2
	RedisDbVideoIdCommentId = redis.NewClient(&redis.Options{
		Addr:     config.Redis_addr_port,
		Password: config.Redis_password,
		DB:       2,
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
	ret, err := RedisDb.SetNX(Ctx, key, value, time.Second*10).Result()
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
