package jwt

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// AuthGet 鉴权中间件，若用户token正确，则解析后将username和id注册到上下文里，否则返回错误信息。
func AuthGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		if len(tokenString) == 0 {
			c.Abort()
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		})

		if err == nil {
			if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
				log.Println("Token right:", claims.Name, claims.Id)
				c.Set("username", claims.Name)
				c.Set("id", claims.Id)
				c.Next()
				return
			}
		}
		log.Println("token error!")
		c.Abort()
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: -1,
			StatusMsg:  "Token Error",
		})

	}
}

// AuthGetWithoutLogin 未登录情况下，若携带token则解析出用户username和id注册到上下文里;若未携带,则username为空，id默认值-1
func AuthGetWithoutLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		if len(tokenString) == 0 {
			c.Set("username", "")
			c.Set("id", -1)
			c.Next()
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		})

		if err == nil {
			if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
				log.Println("Token right:", claims.Name, claims.Id)
				c.Set("username", claims.Name)
				c.Set("id", claims.Id)
				c.Next()
				return
			}
		}
		log.Println("token error!")
		c.Abort()
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: -1,
			StatusMsg:  "Token Error",
		})

	}
}

// AuthPost post请求鉴权
func AuthPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.PostForm("token")
		if tokenString == "" {
			tokenString = c.Query("token")
		}
		if len(tokenString) == 0 {
			c.Abort()
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		})
		if err == nil {
			if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
				log.Println("Token right : ", claims.Name)
				c.Set("username", claims.Name)
				c.Next()
				return
			}
		}
		log.Println("token error!!!")
		c.Abort()
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: -1,
			StatusMsg:  "Token Error",
		})
	}
}
