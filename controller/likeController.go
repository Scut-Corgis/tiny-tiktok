package controller

import (
	"net/http"
	"strconv"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/service"
	"github.com/gin-gonic/gin"
)

type likeResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
type GetLikeListResponse struct {
	StatusCode string          `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string         `json:"status_msg"`  // 返回状态描述
	VideoList  []service.Video `json:"video_list"`  // 用户点赞视频列表
}

// 点赞和取消点赞功能
func FavoriteAction(c *gin.Context) {
	username := c.GetString("username")
	user, err := dao.QueryUserByName(username)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "请登录后点赞！"})
		return
	}

	userId := user.Id
	strVideoId := c.Query("video_id")
	videoId, _ := strconv.ParseInt(strVideoId, 10, 64)
	strActionType := c.Query("action_type")
	actionType, _ := strconv.ParseInt(strActionType, 10, 8)

	if actionType == 1 {
		err := service.Like(userId, videoId)
		if err != nil {
			c.JSON(http.StatusOK, likeResponse{StatusCode: 1, StatusMsg: "点赞失败！"})
		}
		c.JSON(http.StatusOK, likeResponse{StatusCode: 0, StatusMsg: "点赞成功！"})
	} else if actionType == 2 {
		err := service.Unlike(userId, videoId)
		if err != nil {
			c.JSON(http.StatusOK, likeResponse{StatusCode: 1, StatusMsg: "取消点赞失败！"})
		}
		c.JSON(http.StatusOK, likeResponse{StatusCode: 0, StatusMsg: "取消点赞成功！"})
	}
}

func FavoriteList(c *gin.Context) {
	strsuccess := "获取点赞列表成功"
	strfail := "获取点赞列表失败"
	StrUserId := c.Query("user_id")
	UserId, err := strconv.ParseInt(StrUserId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, GetLikeListResponse{
			StatusCode: "1",
			StatusMsg:  &strfail,
			VideoList:  nil,
		})
	}

	ReturnVideos, err := service.GetLikeLists(UserId)
	if err != nil {
		c.JSON(http.StatusOK, GetLikeListResponse{
			StatusCode: "0",
			StatusMsg:  &strfail,
			VideoList:  nil,
		})
	}
	c.JSON(http.StatusOK, GetLikeListResponse{
		StatusCode: "0",
		StatusMsg:  &strsuccess,
		VideoList:  ReturnVideos,
	})
}
