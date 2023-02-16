package controller

import (
	"net/http"
	"strconv"

	"github.com/Scut-Corgis/tiny-tiktok/middleware/jwt"
	"github.com/Scut-Corgis/tiny-tiktok/service"
	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

// Register POST /douyin/user/register/ 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	usi := service.UserServiceImpl{}
	userId, code, message := usi.Register(username, password)

	if code != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
			UserId:   userId,
			Token:    jwt.GenerateToken(username, userId),
		})
	}
}

// Login POST /douyin/user/login/ 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	usi := service.UserServiceImpl{}
	code, message := usi.Login(username, password)
	if code != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
		})
	} else {
		user := usi.QueryUserByName(username)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
			UserId:   user.Id,
			Token:    jwt.GenerateToken(user.Name, user.Id),
		})
	}
}

// UserInfo GET /douyin/user/ 用户信息
func UserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	id, _ := strconv.ParseInt(user_id, 10, 64)
	usi := service.UserServiceImpl{}
	if userResp, err := usi.QueryUserRespById(id); err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else {
		currentUserName := c.GetString("username")
		user := usi.QueryUserByName(currentUserName)
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User: User{
				userResp.Id,
				userResp.Name,
				userResp.FollowCount,
				userResp.FollowerCount,
				usi.JudgeIsFollowById(user.Id, id),
			},
		})
	}
}
