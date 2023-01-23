package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/jwt"
	"github.com/gin-gonic/gin"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

// var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	token := jwt.GenerateToken(username)
	// token := username + password
	user, _ := dao.QueryUserByUsername(username)
	if username == user.Name {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
			UserId:   user.Id,
		})
	} else {
		newUser := dao.User{
			Name:     username,
			Password: password,
		}
		if dao.InsertUser(&newUser) == false {
			log.Println("Insert Data Failed")
		}
		u, _ := dao.QueryUserByUsername(username)
		log.Println("注册返回的 id", u.Id)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Register success"},
			UserId:   u.Id,
			Token:    token,
		})
	}
}

// Login 登录功能，需要补充密码编码解码
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user, err := dao.QueryUserByUsername(username)

	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else {
		if user.Password == password {
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
}

func UserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	id, _ := strconv.ParseInt(user_id, 10, 64)
	if user, err := dao.QueryUserById(id); err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User: User{
				user.Id,
				user.Name,
				0,
				0,
				false,
			},
		})
	}
}
