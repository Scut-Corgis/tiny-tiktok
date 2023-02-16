package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

type MessageService interface {
	// SendMessage 发送消息
	SendMessage(userId int64, toUserId int64, content string) (bool, error)

	// GetChatRecord 读取聊天记录
	GetChatRecord(userId int64, toUserId int64) ([]dao.MessageResp, error)

	// GetLatestMessage 获取最新的聊天记录
	GetLatestMessage(userId1 int64, userId2 int64) (LatestMessage, error)
}
