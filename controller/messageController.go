package controller

import (
	"net/http"
	"strconv"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/service"
	"github.com/gin-gonic/gin"
)

type ChatResponse struct {
	Response
	MessageList []dao.MessageResp `json:"message_list"`
}

/*
处理发送消息请求
*/
func MessageAction(c *gin.Context) {
	username := c.GetString("username")
	user, _ := dao.QueryUserByName(username)
	if username != user.Name {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "token错误"})
		return
	}
	userId := user.Id
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	toUser, _ := dao.QueryUserById(toUserId)
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

	if actionType == 3 {
		flag, err := service.SendMessage(userId, toUserId, content)
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

/*
处理消息列表请求
*/
func ChatRecord(c *gin.Context) {
	username := c.GetString("username")
	user, _ := dao.QueryUserByName(username)
	if username != user.Name {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "token错误"})
		return
	}
	userId := user.Id
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	toUser, _ := dao.QueryUserById(toUserId)
	if nil != err || toUserId == userId {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "接收消息用户id错误"})
		return
	}
	if toUser.Name == "" {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "接收消息用户不存在"})
		return
	}

	chatRecords, err := service.GetChatRecord(userId, toUserId)
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
