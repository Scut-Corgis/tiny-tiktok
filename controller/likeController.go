package controller

import (
	"net/http"
	"strconv"

	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/service"
	"github.com/gin-gonic/gin"
)

type likeResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
type GetLikeListResponse struct {
	StatusCode string  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 用户点赞视频列表
}

// FavoriteAction POST /douyin/favorite/action/ 赞操作
func FavoriteAction(c *gin.Context) {
	userId := c.GetInt64("id") //根据解析token获取赞操作用户id
	favoriteService := service.LikeServiceImpl{}
	strVideoId := c.Query("video_id") //提取url中赞操作的视频id
	videoId, _ := strconv.ParseInt(strVideoId, 10, 64)
	strActionType := c.Query("action_type") //提取url中赞操作的类型 1--点赞 2--取消点赞
	actionType, _ := strconv.ParseInt(strActionType, 10, 8)
	//布谷鸟过滤器，判断视频是否存在。防止恶意攻击，同时点赞不存在的视频导致缓存穿透
	if !redis.CuckooFilterVideoId.Contain([]byte(strconv.FormatInt(videoId, 10))) {
		c.JSON(http.StatusOK, likeResponse{StatusCode: 1, StatusMsg: "视频不存在！"})
		return
	}

	if actionType == 1 { //如果赞操作是1，进行点赞
		err := favoriteService.Like(userId, videoId)
		if err != nil {
			c.JSON(http.StatusOK, likeResponse{StatusCode: 1, StatusMsg: "点赞失败！"})
		}
		c.JSON(http.StatusOK, likeResponse{StatusCode: 0, StatusMsg: "点赞成功！"})
	} else if actionType == 2 {
		err := favoriteService.Unlike(userId, videoId)
		if err != nil { //如果赞操作是2，进行取消点赞
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
	StrUserId := c.Query("user_id") //提取url中获取喜欢列表的用户id
	UserId, err := strconv.ParseInt(StrUserId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, GetLikeListResponse{
			StatusCode: "1",
			StatusMsg:  &strfail,
			VideoList:  nil,
		})
	}

	ReturnVideos, err := favoriteService.GetLikeLists(UserId) //获取喜欢列表中每个视频的详细信息
	var videoList = make([]Video, 0, len(ReturnVideos))
	for _, videoDetail := range ReturnVideos { //将视频的详细信息格式化
		video := Video{
			Id:            videoDetail.Id,
			Author:        User(videoDetail.Author),
			PlayUrl:       videoDetail.PlayUrl,
			CoverUrl:      videoDetail.CoverUrl,
			FavoriteCount: videoDetail.FavoriteCount,
			CommentCount:  videoDetail.CommentCount,
			IsFavorite:    videoDetail.IsFavorite,
			Title:         videoDetail.Title,
		}
		videoList = append(videoList, video)
	}
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
		VideoList:  videoList,
	})
}
