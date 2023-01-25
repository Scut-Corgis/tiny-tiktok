package dao

import (
	"log"
)

func QueryCommentsByVideoId(id int64) ([]Comment, error) {
	var comment []Comment
	if err := Db.Where("video_id = ?", id).Find(&comment).Error; err != nil {
		log.Println(err.Error())
		return comment, err
	}
	return comment, nil
}

//func JudgeCommentIsExist(user_id int64, video_id int64) bool {
//	if err := Db.Where("user_id = ? and video_id = ?", user_id, video_id).Error; err != nil {
//		log.Println(err.Error())
//		return false
//	}
//	return true
//}

func InsertComment(comment *Comment) bool {
	if err := Db.Create(&comment).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func DeleteComment(id int64) bool {
	var comment Comment
	if err := Db.Where("id = ?", id).First(&comment).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	Db.Delete(&comment)
	return true
}
