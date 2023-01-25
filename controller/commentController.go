package controller

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentListResponse struct {
	Response
	CommentList []CommentInfo `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")

	if user, exist := usersLoginInfo[token]; exist {
		if actionType == "1" {
			text := c.Query("comment_text")
			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
				Comment: Comment{
					Id:         1,
					User:       user,
					Content:    text,
					CreateDate: "05-01",
				}})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	video_id := c.Query("video_id")
	currentName := c.GetString("username")
	id, _ := strconv.ParseInt(video_id, 10, 64)
	if comments, err := dao.QueryCommentByVideoId(id); err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Comment doesn't exist"},
		})
	} else {
		var commonList []CommentInfo
		for _, comment := range comments {
			user, _ := dao.QueryUserTableById(comment.UserId)
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
