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
				log.Println("Token right : ", claims.Name)
				c.Set("username", claims.Name)
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

func AuthGetWithoutLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		if len(tokenString) == 0 {
			c.Set("username", "")
			c.Next()
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
		log.Println("token error!")
		c.Abort()
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: -1,
			StatusMsg:  "Token Error",
		})

	}
}

func AuthPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.PostForm("token")
		// log.Println(tokenString)
		// log.Println(len(tokenString))
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
