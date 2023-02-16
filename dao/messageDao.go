package dao

import (
	"log"
)

type Message struct {
	Id         int64  `gorm:"column:id"`
	MessageKey string `gorm:"column:message_key"`
	Content    string `gorm:"column:content"`
	CreateTime string `gorm:"column:create_time"`
}

type MessageResp struct {
	Id         int64  `json:"id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	FromUserId int64  `json:"from_user_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

func InsertMessage(msgKey string, content string, createTime string) (bool, error) {
	message := Message{
		MessageKey: msgKey,
		Content:    content,
		CreateTime: createTime,
	}
	if err := Db.Table("messages").Create(&message).Error; err != nil {
		log.Println(err.Error())
		return false, err
	}
	return true, nil
}

func QueryMessagesByMsgKey(msgKey string) ([]Message, error) {
	message := make([]Message, 1)
	//if err := Db.Raw("select id, content, DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s') create_time from messages where message_key = ?", msgKey).Scan(&messageList).Error; nil != err {
	//	return nil, err
	//}
	if err := Db.Select([]string{"id", "content", "create_time"}).Where("message_key = ?", msgKey).Find(&message).Error; err != nil {
		return nil, err
	}
	return message, nil
}
