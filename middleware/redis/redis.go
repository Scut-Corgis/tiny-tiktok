package redis

import (
	"context"
	"log"

	"github.com/Scut-Corgis/tiny-tiktok/config"
	"github.com/go-redis/redis/v8"
)

// 定义一个全局变量
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
