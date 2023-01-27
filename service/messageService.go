package service

import (
	"fmt"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

/*
发送消息
*/
func SendMessage(userId int64, toUserId int64, content string) (bool, error) {
	msgKey := genMsgKey(userId, toUserId)
	createTime := time.Now().Format("2006-01-02 15:04:05")
	return dao.InsertMessage(msgKey, content, createTime)
}

/*
读取聊天记录
*/
func GetChatRecord(userId int64, toUserId int64) ([]dao.MessageResp, error) {
	msgKey := genMsgKey(userId, toUserId)

	return dao.QueryMessagesByMsgKey(msgKey)
}

/*
生成消息key
参数：sendMsgUserId int64 发送消息用户id，recvMsgUserId int64 接收消息用户id
返回：msgKey string 消息key
*/
func genMsgKey(sendMsgUserId int64, recvMsgUserId int64) string {
	return fmt.Sprintf("%d_%d", sendMsgUserId, recvMsgUserId)
}
