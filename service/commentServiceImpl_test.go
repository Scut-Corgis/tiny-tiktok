package service

import (
	"fmt"
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/rabbitmq"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"
	"testing"
	"time"
)

func CommentServiceImplInit() {
	dao.Init()
	rabbitmq.Init()
	rabbitmq.InitCommentRabbitMQ()
	redis.InitRedis()
	util.InitWordsFilter()
	redis.InitCuckooFilter()
}

func TestCommentServiceImpl_CountComments(t *testing.T) {
	CommentServiceImplInit()
	csi := CommentServiceImpl{}
	cnt, err := csi.CountComments(1000)
	fmt.Println(cnt)
	fmt.Println(err)
}

func TestCommentServiceImpl_QueryCommentsByVideoId(t *testing.T) {
	CommentServiceImplInit()
	csi := CommentServiceImpl{}
	comments := csi.QueryCommentsByVideoId(1000)
	fmt.Println(comments)
}

func TestCommentServiceImpl_PostComment(t *testing.T) {
	CommentServiceImplInit()
	csi := CommentServiceImpl{}
	comment := dao.Comment{
		UserId:      1000,
		VideoId:     1000,
		CommentText: "test",
		CreateDate:  time.Now(),
	}
	id, code, message := csi.PostComment(comment)
	fmt.Println(id, code, message)
}

func TestCommentServiceImpl_DeleteComment(t *testing.T) {
	CommentServiceImplInit()
	csi := CommentServiceImpl{}
	code, message := csi.DeleteComment(1000)
	fmt.Println(code, message)
}
