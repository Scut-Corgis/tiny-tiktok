package controller

import (
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

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	currentUsername := c.GetString("username")
	user, err := dao.QueryUserByName(currentUsername)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	video_id := c.Query("video_id")
	id, _ := strconv.ParseInt(video_id, 10, 64)
	video, _ := dao.QueryVideoById(id)
	actionType := c.Query("action_type")
	if !dao.JudgeVideoIsExist(id) {
		c.JSON(http.StatusOK, likeResponse{StatusCode: 1, StatusMsg: "Video doesn't exist"})
	}
	if actionType == "1" {
		text := c.Query("comment_text")
		user_info, _ := dao.QueryUserRespById(user.Id)
		t := time.Now()
		comment := dao.Comment{
			UserId:      user.Id,
			VideoId:     id,
			CommentText: text,
			CreateDate:  t,
		}
		dao.InsertComment(&comment)
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0, StatusMsg: "success"},
			CommentInfo: CommentInfo{
				comment.Id,
				User{
					user_info.Id,
					user_info.Name,
					user_info.FollowCount,
					user_info.FollowerCount,
					dao.JudgeIsFollowById(user.Id, video.AuthorId),
				},
				text,
				time.Now(),
			},
		})
	} else {
		comment_id := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(comment_id, 10, 64)
		if !dao.DeleteComment(commentId) {
			c.JSON(http.StatusOK, likeResponse{StatusCode: 1, StatusMsg: "Comment doesn't exist"})
		} else {
			c.JSON(http.StatusOK, likeResponse{StatusCode: 0, StatusMsg: "Delete successfully!"})
		}
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	video_id := c.Query("video_id")
	currentName := c.GetString("username")
	id, _ := strconv.ParseInt(video_id, 10, 64)
	if comments, err := dao.QueryCommentsByVideoId(id); err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Comment doesn't exist"},
		})
	} else {
		var commonList []CommentInfo
		for _, comment := range comments {
			user, _ := dao.QueryUserRespById(comment.UserId)
			commonList = append(commonList, CommentInfo{
				comment.Id,
				User{
					user.Id,
					user.Name,
					user.FollowCount,
					user.FollowerCount,
					dao.JudgeIsFollow(id, currentName),
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
}
