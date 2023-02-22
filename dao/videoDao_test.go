package dao

import (
	"fmt"
	"testing"
	"time"
)

func TestQueryVideoById(t *testing.T) {
	Init()
	video, err := QueryVideoById(1000)
	fmt.Println(video)
	fmt.Println(err)
}

func TestInsertVideosTable(t *testing.T) {
	Init()
	video := Video{
		AuthorId:    1000,
		PlayUrl:     "http://xxx.com",
		CoverUrl:    "http://xxx.com",
		PublishTime: time.Now(),
		Title:       "test",
	}
	_, err := InsertVideosTable(video)
	fmt.Println(err)
}

func TestGetMost30videosIdList(t *testing.T) {
	Init()
	videos := GetMost30videosIdList(time.Now())
	fmt.Println(videos)
}

func TestGetVideoIdListByUserId(t *testing.T) {
	Init()
	videos := GetVideoIdListByUserId(1000)
	fmt.Println(videos)
}

func TestQueryAllVideoIds(t *testing.T) {
	Init()
	videos := QueryAllVideoIds()
	fmt.Println(videos)
}
