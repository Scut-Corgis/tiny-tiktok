package controller

import (
	"log"
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

type FriendListResponse struct {
	Response
	FriendList []dao.FriendResp `json:"user_list"`
}

// RelationAction POST /douyin/relation/action/ 关注操作
func RelationAction(c *gin.Context) {
	// Step1. 取出用户id
	userId := c.GetInt64("id")
	usi := service.UserServiceImpl{}
	// Step2. 判断to_user_id解析是否有误
	followId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if nil != err {
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "关注用户id错误"})
		return
	}
	followUser := usi.QueryUserById(followId)
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

// FollowList GET /douyin/relation/follow/list/ 关注列表
func FollowList(c *gin.Context) {
	// Step1. 判断user_id解析是否有误
	realUserId := c.GetInt64("id")
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
	// Step1. 判断user_id解析是否有误
	realUserId := c.GetInt64("id")
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
	// Step1. 判断user_id解析是否有误
	realUserId := c.GetInt64("id")
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || realUserId != userId {
		c.JSON(http.StatusOK, FriendListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "用户id有误",
			},
			FriendList: nil,
		})
		return
	}
	// Step2. 判断获取好友列表是否有误
	rsi := service.RelationServiceImpl{}
	friendList, err := rsi.GetFriendList(userId)
	if err != nil {
		c.JSON(http.StatusOK, FriendListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取好友列表失败",
			},
			FriendList: nil,
		})
		return
	}
	log.Println(friendList)
	c.JSON(http.StatusOK, FriendListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "获取好友列表成功！",
		},
		FriendList: friendList,
	})
}
