package controller

import (
	"net/http"
	"strconv"

	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/service"
	"github.com/gin-gonic/gin"
)

type likeResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
type GetLikeListResponse struct {
	StatusCode string            `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string           `json:"status_msg"`  // 返回状态描述
	VideoList  []dao.VideoDetail `json:"video_list"`  // 用户点赞视频列表
}

// FavoriteAction POST /douyin/favorite/action/ 赞操作
func FavoriteAction(c *gin.Context) {
	userId := c.GetInt64("id")
	favoriteService := service.LikeServiceImpl{}
	//user := favoriteService.QueryUserByName(username) //这种调用会出错！
	//user := service.UserServiceImpl{}.QueryUserByName(username)
	//user, err := dao.QueryUserByName(username)
	// if err != nil {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "查询用户出错"})
	// 	return
	// }

	//userId, _ := strconv.ParseInt(StrUserId, 10, 64)
	strVideoId := c.Query("video_id")
	videoId, _ := strconv.ParseInt(strVideoId, 10, 64)
	strActionType := c.Query("action_type")
	actionType, _ := strconv.ParseInt(strActionType, 10, 8)
	if !redis.CuckooFilterVideoId.Contain([]byte(strconv.FormatInt(videoId, 10))) {
		c.JSON(http.StatusOK, likeResponse{StatusCode: 1, StatusMsg: "视频不存在！"})
		return
	}

	if actionType == 1 {
		err := favoriteService.Like(userId, videoId)
		if err != nil {
			c.JSON(http.StatusOK, likeResponse{StatusCode: 1, StatusMsg: "点赞失败！"})
		}
		c.JSON(http.StatusOK, likeResponse{StatusCode: 0, StatusMsg: "点赞成功！"})
	} else if actionType == 2 {
		err := favoriteService.Unlike(userId, videoId)
		if err != nil {
			c.JSON(http.StatusOK, likeResponse{StatusCode: 1, StatusMsg: "取消点赞失败！"})
		}
		c.JSON(http.StatusOK, likeResponse{StatusCode: 0, StatusMsg: "取消点赞成功！"})
	}
}

// FavoriteList GET /douyin/favorite/list/ 喜欢列表
func FavoriteList(c *gin.Context) {
	favoriteService := service.LikeServiceImpl{}
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

	ReturnVideos, err := favoriteService.GetLikeLists(UserId)
	if err != nil {
		c.JSON(http.StatusOK, GetLikeListResponse{
			StatusCode: "1",
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
