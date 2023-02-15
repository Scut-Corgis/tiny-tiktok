package controller

import (
	"github.com/Scut-Corgis/tiny-tiktok/service"
	"github.com/Scut-Corgis/tiny-tiktok/util"
	"net/http"
	"strconv"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []CommentInfo `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	CommentInfo CommentInfo `json:"comment,omitempty"`
}

// CommentAction POST /douyin/comment/action/ 评论操作
func CommentAction(c *gin.Context) {
	csi := service.CommentServiceImpl{}
	usi := service.UserServiceImpl{}

	// 获取当前用户
	currentUserName := c.GetString("username")
	user := usi.QueryUserByName(currentUserName)

	// 获取当前视频
	id := c.Query("video_id")
	videoId, _ := strconv.ParseInt(id, 10, 64)
	video := csi.QueryVideoById(videoId)

	actionType := c.Query("action_type")
	if actionType == "1" {
		text := c.Query("comment_text")
		text = util.Filter.Replace(text, '#') // 评论敏感词过滤
		t := time.Now()
		comment := dao.Comment{
			UserId:      user.Id,
			VideoId:     videoId,
			CommentText: text,
			CreateDate:  t,
		}
		commentId, code, message := csi.PostComment(comment)
		if code != 0 {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: code, StatusMsg: message},
			})
		}
		userInfo, _ := usi.QueryUserRespById(user.Id)
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
			CommentInfo: CommentInfo{
				Id: commentId,
				User: User{
					Id:            userInfo.Id,
					Name:          userInfo.Name,
					FollowCount:   userInfo.FollowCount,
					FollowerCount: userInfo.FollowerCount,
					IsFollow:      usi.JudgeIsFollowById(userInfo.Id, video.AuthorId),
				},
				CommentText: text,
				CreateDate:  t,
			},
		})
	} else {
		id := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(id, 10, 64)
		code, message := csi.DeleteComment(commentId)
		c.JSON(http.StatusOK, likeResponse{StatusCode: code, StatusMsg: message})
	}
}

// CommentList GET /douyin/comment/list/ 评论列表
func CommentList(c *gin.Context) {
	usi := service.UserServiceImpl{}
	csi := service.CommentServiceImpl{}

	id := c.Query("video_id")
	videoId, _ := strconv.ParseInt(id, 10, 64)
	video := csi.QueryVideoById(videoId)

	//if !dao.JudgeVideoIsExist(id) {
	//	c.JSON(http.StatusOK, likeResponse{StatusCode: 1, StatusMsg: "Video doesn't exist"})
	//}
	comments := csi.QueryCommentsByVideoId(videoId)
	var commonList []CommentInfo
	for _, comment := range comments {
		user, err := usi.QueryUserRespById(comment.UserId)
		if err != nil {
			continue
		}
		commonList = append(commonList, CommentInfo{
			comment.Id,
			User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      usi.JudgeIsFollowById(user.Id, video.AuthorId),
			},
			comment.CommentText,
			comment.CreateDate,
		})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "success"},
		CommentList: commonList,
	})
}
