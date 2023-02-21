package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"time"
)

type VideoService interface {
	// QueryVideoById 根据视频id获取视频
	QueryVideoById(id int64) dao.Video

	// QueryVideoDetailByVideoId 根据视频id和查询用户id查询视频的详细信息
	QueryVideoDetailByVideoId(videoId int64, queryUserId int64) (dao.VideoDetail, time.Time)

	// GetVideoIdListByUserId 根据作者id查询视频id列表
	GetVideoIdListByUserId(authorId int64) []int64

	// GetMost30videosIdList 根据时间获取视频id列表
	GetMost30videosIdList(latestTime time.Time) []int64

	// InsertVideosTable 将video插入videos表内
	InsertVideosTable(video dao.Video) bool

	// CountWorks 统计用户id的作品数
	CountWorks(id int64) int
}
