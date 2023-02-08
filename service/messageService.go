package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

type MessageService interface {
	// SendMessage 发送消息
	SendMessage(userId int64, toUserId int64, content string) (bool, error)

	// GetChatRecord 读取聊天记录
	GetChatRecord(userId int64, toUserId int64) ([]dao.MessageResp, error)
}
