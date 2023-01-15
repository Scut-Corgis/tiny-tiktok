package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//	go service.RunMessageServer() 消息聊天处理协程

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
