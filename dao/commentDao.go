package dao

import (
	"log"
)

// CommentCount 根据视频id统计评论数量
func CommentCount(id int64) (int64, error) {
	var count int64
	err := Db.Model(&Comment{}).Where("video_id = ?", id).Count(&count).Error
	if err != nil {
		return -1, err
	}
	return count, nil
}

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
func InsertComment(comment Comment) (Comment, error) {
	err := Db.Create(&comment).Error
	if err != nil {
		log.Println(err.Error())
		return Comment{}, err
	}
	return comment, err
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
