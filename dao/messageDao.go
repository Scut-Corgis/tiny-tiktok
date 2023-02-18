package dao

import (
	"log"
)

type Message struct {
	Id         int64  `gorm:"column:id"`
	ToUserId   int64  `gorm:"column:to_user_id"`
	FromUserId int64  `gorm:"column:from_user_id"`
	Content    string `gorm:"column:content"`
	CreateTime string `gorm:"column:create_time"`
}

type MessageResp struct {
	Id         int64  `json:"id"`
	ToUserId   int64  `json:"to_user_id"`
	FromUserId int64  `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

// InsertMessage 插入message
func InsertMessage(userId int64, toUserId int64, content string, createTime string) (int64, error) {
	message := Message{
		ToUserId:   toUserId,
		FromUserId: userId,
		Content:    content,
		CreateTime: createTime,
	}
	if err := Db.Table("messages").Create(&message).Error; err != nil {
		log.Println(err.Error())
		return -1, err
	}
	return message.Id, nil
}

// QueryMessagesByMsgKey 根据userId和toUserId获取所有message记录
func QueryMessagesByMsgKey(userId int64, toUserId int64) ([]Message, error) {
	message := make([]Message, 1)
	if err := Db.Table("messages").Where("to_user_id = ? AND from_user_id = ?", toUserId, userId).Or("to_user_id = ? AND from_user_id = ?", userId, toUserId).Find(&message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

// QueryLatestMessageByUserId 根据userId和toUserId获取最新的message记录
func QueryLatestMessageByUserId(userId int64, toUserId int64) (Message, error) {
	message := Message{}
	if err := Db.Table("messages").Where("to_user_id = ? AND from_user_id = ?", toUserId, userId).Or("to_user_id = ? AND from_user_id = ?", userId, toUserId).Last(&message).Error; err != nil {
		return message, err
	}
	return message, nil
}
