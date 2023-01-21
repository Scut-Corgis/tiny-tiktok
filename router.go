package main

import (
	"github.com/Scut-Corgis/tiny-tiktok/controller"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/jwt"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	// r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", jwt.AuthGet(), controller.Feed)
	apiRouter.GET("/user/", jwt.AuthGet(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", jwt.AuthPost(), controller.Publish)
	apiRouter.GET("/publish/list/", jwt.AuthGet(), controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", jwt.AuthPost(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", jwt.AuthGet(), controller.FavoriteList)
	apiRouter.POST("/comment/action/", jwt.AuthPost(), controller.CommentAction)
	apiRouter.GET("/comment/list/", jwt.AuthGet(), controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", jwt.AuthGet(), controller.FollowList)
	apiRouter.GET("/relation/follower/list/", jwt.AuthGet(), controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", jwt.AuthGet(), controller.FriendList)
	apiRouter.GET("/message/chat/", jwt.AuthGet(), controller.MessageChat)
	apiRouter.POST("/message/action/", jwt.AuthPost(), controller.MessageAction)
}
