package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/rabbitmq"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"
)

func VideoServiceImplInit() {
	dao.Init()
	rabbitmq.Init()
	rabbitmq.InitCommentRabbitMQ()
	redis.InitRedis()
	util.InitWordsFilter()
	redis.InitCuckooFilter()
}

func TestVideoServiceImpl_QueryVideoById(t *testing.T) {
	VideoServiceImplInit()
	vsi := VideoServiceImpl{}
	video, _ := vsi.QueryVideoById(1000)
	fmt.Println(video)
}

func TestVideoServiceImpl_QueryVideoDetailByVideoId(t *testing.T) {
	VideoServiceImplInit()
	vsi := VideoServiceImpl{}
	videoDetail, publishTime := vsi.QueryVideoDetailByVideoId(1000, 1000)
	fmt.Println(videoDetail)
	fmt.Println(publishTime)
}

func TestVideoServiceImpl_GetMost30videosIdList(t *testing.T) {
	VideoServiceImplInit()
	vsi := VideoServiceImpl{}
	videos := vsi.GetMost30videosIdList(time.Now())
	fmt.Println(videos)
}

func TestVideoServiceImpl_GetVideoIdListByUserId(t *testing.T) {
	VideoServiceImplInit()
	vsi := VideoServiceImpl{}
	videos := vsi.GetVideoIdListByUserId(1000)
	fmt.Println(videos)
}

func TestVideoServiceImpl_InsertVideosTable(t *testing.T) {
	VideoServiceImplInit()
	vsi := VideoServiceImpl{}
	video := dao.Video{
		AuthorId:    1000,
		PlayUrl:     "http://xxx.com",
		CoverUrl:    "http://xxx.com",
		PublishTime: time.Now(),
		Title:       "test",
	}
	flag := vsi.InsertVideosTable(video)
	fmt.Println(flag)
}
