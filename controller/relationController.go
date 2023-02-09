package controller

import (
	"net/http"
	"strconv"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/service"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []dao.UserResp `json:"user_list"`
}

/*
处理关注和取关接口
*/
func RelationAction(c *gin.Context) {
	//#优化：若token含userid，获取用户可以省去查数据库操作，或使用redis减少对数据库的访问
	// Step1. 判断token解析是否有误
	username := c.GetString("username")
	usi := service.UserServiceImpl{}
	user := usi.QueryUserByName(username)
	if username != user.Name {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "token错误"})
		return
	}
	userId := user.Id
	// Step2. 判断to_user_id解析是否有误
	followId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	followUser := usi.QueryUserById(followId)
	if nil != err {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "关注用户id错误"})
		return
	}
	if followUser.Name == "" {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "该用户不存在"})
		return
	}
	if followId == userId {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "不可关注/取关自己"})
		return
	}
	// Step3. 判断action_type解析是否有误
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if nil != err || actionType < 1 || actionType > 2 {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "关注/取关失败"})
		return
	}
	// Step4. 关注或取关
	rsi := service.RelationServiceImpl{}
	switch {
	case actionType == 1:
		// 关注
		flag, err := rsi.Follow(userId, followId)
		if nil != err {
			c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "关注失败"})
			return
		}
		if !flag {
			c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "该用户已关注"})
			return
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "关注成功!",
		})
		return
	case actionType == 2:
		// 取关
		flag, err := rsi.UnFollow(userId, followId)
		if nil != err {
			c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "取关失败"})
			return
		}
		if !flag {
			c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "该用户未关注,取关失败"})
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
	// Step1. 判断token、user_id解析是否有误
	username := c.GetString("username")
	usi := service.UserServiceImpl{}
	user := usi.QueryUserByName(username)
	if username != user.Name {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "token错误"})
		return
	}
	realUserId := user.Id
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || realUserId != userId {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "用户id有误",
			},
			UserList: nil,
		})
		return
	}
	// Step2. 判断获取关注列表是否有误
	rsi := service.RelationServiceImpl{}
	followList, err := rsi.GetFollowList(userId)
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
			StatusMsg:  "获取关注列表成功！",
		},
		UserList: followList,
	})
}

/*
处理获取当前用户的粉丝列表
*/
func FollowerList(c *gin.Context) {
	// Step1. 判断token、user_id解析是否有误
	username := c.GetString("username")
	usi := service.UserServiceImpl{}
	user := usi.QueryUserByName(username)
	if username != user.Name {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  "token错误",
		})
		return
	}
	realUserId := user.Id
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || realUserId != userId {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "用户id有误",
			},
			UserList: nil,
		})
		return
	}
	// Step2. 判断获取粉丝列表是否有误
	rsi := service.RelationServiceImpl{}
	followerList, err := rsi.GetFollowerList(userId)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取粉丝列表失败",
			},
			UserList: nil,
		})
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "获取粉丝列表成功！",
		},
		UserList: followerList,
	})
}

/*
处理获取好友/互关列表接口
*/
func FriendList(c *gin.Context) {
	// ------------------------------------------------------------
	// 方案说明中 好友列表描述与粉丝列表相同，直接调用FollowerList即可
	// FollowerList(c)
	// 这里将好友列表请求识别为互关列表，执行以下逻辑
	// ------------------------------------------------------------
	// Step1. 判断token、user_id解析是否有误
	username := c.GetString("username")
	usi := service.UserServiceImpl{}
	user := usi.QueryUserByName(username)
	if username != user.Name {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  "token错误",
		})
		return
	}
	realUserId := user.Id
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || realUserId != userId {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "用户id有误",
			},
			UserList: nil,
		})
		return
	}
	// Step2. 判断获取好友列表是否有误
	rsi := service.RelationServiceImpl{}
	friendList, err := rsi.GetFriendList(userId)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取好友列表失败",
			},
			UserList: nil,
		})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "获取好友列表成功！",
		},
		UserList: friendList,
	})
}
