package service

import (
	"fmt"
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/rabbitmq"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"testing"
)

func LikeServiceImplInit() {
	dao.Init()
	redis.InitRedis()
	rabbitmq.Init()
	rabbitmq.InitLikeRabbitMQ()
}

func TestLikeServiceImpl_Like(t *testing.T) {
	LikeServiceImplInit()
	lsi := LikeServiceImpl{}
	err := lsi.Like(1000, 1000)
	fmt.Println(err)
}

func TestLikeServiceImpl_Unlike(t *testing.T) {
	LikeServiceImplInit()
	lsi := LikeServiceImpl{}
	err := lsi.Unlike(1000, 1000)
	fmt.Println(err)
}

func TestLikeServiceImpl_GetLikeLists(t *testing.T) {
	LikeServiceImplInit()
	lsi := LikeServiceImpl{}
	likes, err := lsi.GetLikeLists(1000)
	fmt.Println(likes)
	fmt.Println(err)
}

func TestLikeServiceImpl_IsLike(t *testing.T) {
	LikeServiceImplInit()
	lsi := LikeServiceImpl{}
	flag, err := lsi.IsLike(1000, 1000)
	fmt.Println(flag, err)
}

func TestLikeServiceImpl_LikeCount(t *testing.T) {
	LikeServiceImplInit()
	lsi := LikeServiceImpl{}
	cnt, err := lsi.LikeCount(1000)
	fmt.Println(cnt, err)
}
