package dao

import (
	"log"
)

// QueryCommentsByVideoId 根据视频id查询评论列表
func QueryCommentsByVideoId(id int64) ([]Comment, error) {
	var comment []Comment
	if err := Db.Where("video_id = ?", id).Find(&comment).Error; err != nil {
		log.Println(err.Error())
		return comment, err
	}
	return comment, nil
}

// InsertComment 将comment插入到comments表内
func InsertComment(comment *Comment) bool {
	if err := Db.Create(&comment).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

// DeleteComment 根据评论id将评论删除
func DeleteComment(id int64) bool {
	var comment Comment
	if err := Db.Where("id = ?", id).First(&comment).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	Db.Delete(&comment)
	return true
}
