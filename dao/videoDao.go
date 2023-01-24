package dao

import (
	"log"
	"time"
)

type Video struct {
	Id          int64
	AuthorId    int64
	PlayUrl     string
	CoverUrl    string
	PublishTime time.Time
	Title       string
}

func InsertVideosTable(video *Video) error {
	if err := Db.Create(&video).Error; err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
