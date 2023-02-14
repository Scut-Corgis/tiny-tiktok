package main

import (
	"log"

	"github.com/Scut-Corgis/tiny-tiktok/util"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/ffmpeg"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/ftp"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/rabbitmq"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode("release") //关闭gin debug日志信息，没什么用
	r := gin.Default()

	//	go service.RunMessageServer() 消息聊天处理协程

	initRouter(r)
	initDependencies()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// 加载项目依赖
func initDependencies() {
	log.SetFlags(log.Lshortfile)
	dao.Init()
	rabbitmq.Init()
	rabbitmq.InitLikeRabbitMQ()
	redis.InitRedis()
	ffmpeg.Init()
	ftp.Init()
	util.InitWordsFilter()

}
