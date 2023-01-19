package jwt

import (
	"fmt"
	"net/http"

	"github.com/Scut-Corgis/tiny-tiktok/controller"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		if len(tokenString) == 0 {
			c.Abort()
			c.JSON(http.StatusUnauthorized, controller.Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
		}
		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		})

		if err == nil {
			if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
				fmt.Println("token right!")
				c.Set("username", claims.Name)
				c.Next()
			}
		}

		c.Abort()
		c.JSON(http.StatusUnauthorized, controller.Response{
			StatusCode: -1,
			StatusMsg:  "Token Error",
		})

	}
}

func AuthPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.PostForm("token")
		if len(tokenString) == 0 {
			c.Abort()
			c.JSON(http.StatusUnauthorized, controller.Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
		}
		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		})

		if err == nil {
			if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
				fmt.Println("token right!")
				c.Set("username", claims.Name)
				c.Next()
			}
		}

		c.Abort()
		c.JSON(http.StatusUnauthorized, controller.Response{
			StatusCode: -1,
			StatusMsg:  "Token Error",
		})

	}
}
