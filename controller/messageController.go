package controller

import (
	"net/http"
	"strconv"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/service"
	"github.com/Scut-Corgis/tiny-tiktok/util"
	"github.com/gin-gonic/gin"
)

type ChatResponse struct {
	Response
	MessageList []dao.MessageResp `json:"message_list"`
}

// MessageAction POST /douyin/message/action/ 发送信息
func MessageAction(c *gin.Context) {
	userId := c.GetInt64("id")
	usi := service.UserServiceImpl{}
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	toUser := usi.QueryUserById(toUserId)
	if nil != err {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "接收消息用户id错误"})
		return
	}
	if toUser.Name == "" {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "该用户不存在"})
		return
	}
	if toUserId == userId {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "不可给自己发消息"})
		return
	}
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if nil != err {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "发送消息操作错误"})
		return
	}
	content := c.Query("content")
	// 评论敏感词过滤
	content = util.Filter.Replace(content, '*')

	msi := service.MessageServiceImpl{}
	if actionType == 1 {
		flag, err := msi.SendMessage(userId, toUserId, content)
		if nil != err || !flag {
			c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "发送消息失败"})
			return
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "发送成功",
		})
		return
	}

}

// ChatRecord GET /douyin/message/chat/ 聊天记录
func ChatRecord(c *gin.Context) {
	userId := c.GetInt64("id")
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	usi := service.UserServiceImpl{}
	toUser := usi.QueryUserById(toUserId)
	if nil != err || toUserId == userId {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "接收消息用户id错误"})
		return
	}
	if toUser.Name == "" {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "接收消息用户不存在"})
		return
	}

	msi := service.MessageServiceImpl{}
	chatRecords, err := msi.GetChatRecord(userId, toUserId)
	if nil != err {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "聊天记录请求失败"})
		return
	}
	c.JSON(http.StatusOK, ChatResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "聊天记录请求成功",
		},
		MessageList: chatRecords,
	})
}
