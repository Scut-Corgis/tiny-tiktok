package service

import (
	"fmt"
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"
	"testing"
)

func MessageServiceImplInit() {
	dao.Init()
	redis.InitRedis()
	util.InitWordsFilter()
}

func TestMessageServiceImpl_SendMessage(t *testing.T) {
	MessageServiceImplInit()
	msi := MessageServiceImpl{}
	flag, err := msi.SendMessage(1000, 1000, "text")
	fmt.Println(flag, err)
}

func TestMessageServiceImpl_GetChatRecord(t *testing.T) {
	MessageServiceImplInit()
	msi := MessageServiceImpl{}
	messages, err := msi.GetChatRecord(1000, 1000)
	fmt.Println(messages)
	fmt.Println(err)
}

func TestMessageServiceImpl_GetLatestMessage(t *testing.T) {
	MessageServiceImplInit()
	msi := MessageServiceImpl{}
	message, err := msi.GetLatestMessage(1000, 1000)
	fmt.Println(message)
	fmt.Println(err)
}
