package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"log"
)

type VideoServiceImpl struct{}

func (VideoServiceImpl) QueryVideoById(id int64) dao.Video {
	video, err := dao.QueryVideoById(id)
	if err != nil {
		log.Println("error:", err.Error())
		log.Println("Video not found!")
		return video
	}
	log.Println("Query video successfully!")
	return video
}
