package redis

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/util"

	"github.com/Scut-Corgis/tiny-tiktok/config"
	"github.com/go-redis/redis/v8"
)

var mutex sync.Mutex

// RedisDb 定义一个全局变量
var RedisDb *redis.Client
var RedisDbBackUp *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     config.Redis_addr_port,
		Password: config.Redis_password,
		DB:       0,
	})
	RedisDbBackUp = redis.NewClient(&redis.Options{
		Addr:     config.Redis_addr_port,
		Password: config.Redis_password,
		DB:       1,
	})
	_, err1 := RedisDb.Ping(Ctx).Result()
	if err1 != nil {
		log.Panicln("err:", err1.Error())
	}
	_, err2 := RedisDbBackUp.Ping(Ctx).Result()
	if err2 != nil {
		log.Panicln("err:", err2.Error())
	}
	log.Println("redis has connected!")
	BackUp()
	log.Println("redis backup successfully!")
}

func BackUp() {
	keys, _ := RedisDb.Keys(Ctx, "*").Result()
	for _, key := range keys {
		RedisDb.Copy(Ctx, key, key, 1, true)
	}
}

func Lock(key string, value string) bool {
	mutex.Lock() // 保证程序不存在并发冲突问题
	defer mutex.Unlock()
	ret, err := RedisDb.SetNX(Ctx, key, value, time.Second).Result()
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

// 匹配指定前缀删除redis缓存
func DelRedisCatchBatch(keys ...string) {
	for _, redisKey := range keys {
		keysMatch, err := RedisDb.Do(Ctx, "keys", redisKey+"*").Result()
		if err != nil {
			log.Println(err)
		}
		if reflect.TypeOf(keysMatch).Kind() == reflect.Slice {
			val := reflect.ValueOf(keysMatch)
			if val.Len() == 0 {
				continue
			}
			for i := 0; i < val.Len(); i++ {
				RedisDb.Del(Ctx, val.Index(i).Interface().(string))
				fmt.Printf("删除了rediskey::%s \n", val.Index(i).Interface().(string))
			}
		}
	}
}
