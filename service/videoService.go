package service

import "github.com/Scut-Corgis/tiny-tiktok/dao"

type VideoService interface {
	// QueryVideoById 根据视频id获取视频
	QueryVideoById(id int64) dao.Video
}
