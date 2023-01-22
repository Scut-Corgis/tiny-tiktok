package controller

import (
	"net/http"
	"strconv"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/jwt"
	"github.com/Scut-Corgis/tiny-tiktok/service"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []dao.UserTable `json:"user_list"`
}

/*
处理关注和取关接口
*/
func RelationAction(c *gin.Context) {
	// 若token含userid，获取用户可以省去查数据库操作，或使用redis减少对数据库的访问
	jwt.AuthPost()
	username := c.GetString("username")
	user, _ := dao.QueryUserByUsername(username)
	if username != user.Username {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  "token错误",
		})
		return
	}
	userId := user.Id

	followId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if nil != err {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  "关注用户id错误",
		})
		return
	}
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if nil != err || (actionType < 1 && actionType > 2) {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  "关注/取关操作错误",
		})
		return
	}

	switch {
	case actionType == 1:
		flag := service.Follow(userId, followId)
		if !flag {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				StatusMsg:  "关注失败",
			})
			return
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "关注成功!",
		})
		return
	case actionType == 2:
		flag := service.UnFollow(userId, followId)
		if !flag {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				StatusMsg:  "取关失败",
			})
			return
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "取关成功!",
		})
		return
	}

}

/*
处理获取当前用户的关注列表
*/
func FollowList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "用户id有误",
			},
			UserList: nil,
		})
		return
	}
	followList, err := service.FollowList(userId)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取关注列表失败",
			},
			UserList: nil,
		})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: followList,
	})
}

/*
处理获取当前用户的粉丝列表
*/
func FollowerList(c *gin.Context) {
	// userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusOK, UserListResponse{
	// 		Response: Response{
	// 			StatusCode: -1,
	// 			StatusMsg:  "用户id有误",
	// 		},
	// 		UserList: nil,
	// 	})
	// 	return
	// }
	// followList, err := service.FollowList(userId)
	// if err != nil {
	// 	c.JSON(http.StatusOK, UserListResponse{
	// 		Response: Response{
	// 			StatusCode: -1,
	// 			StatusMsg:  "获取关注列表失败",
	// 		},
	// 		UserList: nil,
	// 	})
	// 	return
	// }

	// c.JSON(http.StatusOK, UserListResponse{
	// 	Response: Response{
	// 		StatusCode: 0,
	// 	},
	// 	UserList: followList,
	// })
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	// c.JSON(http.StatusOK, UserListResponse{
	// 	Response: Response{
	// 		StatusCode: 0,
	// 	},
	// 	UserList: []User{DemoUser},
	// })
}
