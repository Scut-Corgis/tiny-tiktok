package controller

import (
	"net/http"
	"strconv"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
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
	user := usi.QueryUserByName(username)
	if username == user.Name {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist!"},
		})
	} else {
		encoderPassword, err := service.HashEncode(password)
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Incorrect password format!"},
			})
		}
		newUser := dao.User{
			Name:     username,
			Password: encoderPassword,
		}
		if !usi.InsertUser(&newUser) {
			println("Insert user failed！")
		}
		user := usi.QueryUserByName(username)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Register successfully!"},
			UserId:   user.Id,
			Token:    jwt.GenerateToken(username),
		})
	}
}

// Login POST /douyin/user/login/ 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	usi := service.UserServiceImpl{}
	user := usi.QueryUserByName(username)
	if service.ComparePasswords(user.Password, password) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Login success"},
			UserId:   user.Id,
			Token:    jwt.GenerateToken(user.Name),
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Username or Password error"},
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
		currentName := c.GetString("username")
		user := usi.QueryUserByName(currentName)
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User: User{
				userResp.Id,
				userResp.Name,
				userResp.FollowCount,
				userResp.FollowerCount,
				usi.JudgeIsFollowById(id, user.Id),
			},
		})
	}
}
