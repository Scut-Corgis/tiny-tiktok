package controller

import (
	"github.com/Scut-Corgis/tiny-tiktok/service"
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
	currentName := c.GetString("username")
	user := usi.QueryUserByName(currentName)

	// 获取当前视频
	video_id := c.Query("video_id")
	id, _ := strconv.ParseInt(video_id, 10, 64)
	video, _ := dao.QueryVideoById(id)

	actionType := c.Query("action_type")
	if !dao.JudgeVideoIsExist(id) {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Video doesn't exist"},
		})
	}
	if actionType == "1" {
		text := c.Query("comment_text")
		t := time.Now()
		comment := dao.Comment{
			UserId:      user.Id,
			VideoId:     id,
			CommentText: text,
			CreateDate:  t,
		}
		if !csi.InsertComment(&comment) {
			println("Insert comment failed!")
		}
		userInfo, _ := usi.QueryUserRespById(user.Id)
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0, StatusMsg: "success"},
			CommentInfo: CommentInfo{
				Id: comment.Id,
				User: User{
					Id:            userInfo.Id,
					Name:          userInfo.Name,
					FollowCount:   userInfo.FollowCount,
					FollowerCount: userInfo.FollowerCount,
					IsFollow:      usi.JudgeIsFollowById(userInfo.Id, video.AuthorId),
				},
				CommentText: text,
				CreateDate:  time.Now(),
			},
		})
	} else {
		comment_id := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(comment_id, 10, 64)
		if !csi.DeleteComment(commentId) {
			c.JSON(http.StatusOK, likeResponse{StatusCode: 1, StatusMsg: "Comment not found!"})
		} else {
			c.JSON(http.StatusOK, likeResponse{StatusCode: 0, StatusMsg: "Comment delete successfully!"})
		}
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	usi := service.UserServiceImpl{}
	csi := service.CommentServiceImpl{}

	video_id := c.Query("video_id")
	id, _ := strconv.ParseInt(video_id, 10, 64)
	video, _ := dao.QueryVideoById(id)

	if !dao.JudgeVideoIsExist(id) {
		c.JSON(http.StatusOK, likeResponse{StatusCode: 1, StatusMsg: "Video doesn't exist"})
	}
	comments := csi.QueryCommentsByVideoId(id)
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
