package dao

import (
	"log"
	"time"
)

type Message struct {
	Id         int64     `gorm:"column:id"`
	ToUserId   int64     `gorm:"column:to_user_id"`
	FromUserId int64     `gorm:"column:from_user_id"`
	Content    string    `gorm:"column:content"`
	CreateTime time.Time `gorm:"column:create_time"`
}

type MessageResp struct {
	Id         int64  `json:"id"`
	ToUserId   int64  `json:"to_user_id"`
	FromUserId int64  `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

type LatestMessage struct {
	Id         int64  `json:"id"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
	MsgType    int64  `json:"msgType"` // 1为发送方，0为接收方
}

// InsertMessage 插入message
func InsertMessage(userId int64, toUserId int64, content string, createTime time.Time) (int64, error) {
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
	message := make([]Message, 0)
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

// QueryChartsAfterLatestMessage 根据userId、toUserId、最新消息的时间获取剩余聊天记录
func QueryChartsAfterLatestMsg(userId int64, toUserId int64, latestTime time.Time) ([]Message, error) {
	message := make([]Message, 0)
	if err := Db.Table("messages").Where("to_user_id = ? AND from_user_id = ? AND create_time > ?", toUserId, userId, latestTime).Or("to_user_id = ? AND from_user_id = ? AND create_time > ?", userId, toUserId, latestTime).Find(&message).Error; err != nil {
		return message, err
	}
	return message, nil
}
