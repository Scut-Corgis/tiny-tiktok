package main

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//	go service.RunMessageServer() 消息聊天处理协程

	initRouter(r)
	initDependencies()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// 加载项目依赖
func initDependencies() {
	dao.Init()
}
