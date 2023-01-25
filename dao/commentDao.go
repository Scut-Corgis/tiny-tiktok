package dao

import (
	"log"
)

func QueryCommentByVideoId(id int64) ([]Comment, error) {
	var comment []Comment
	if err := Db.Where("video_id = ?", id).Find(&comment).Error; err != nil {
		log.Println(err.Error())
		return comment, err
	}
	return comment, nil
}

func InsertComment(comment *Comment) bool {
	if err := Db.Create(&comment).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
