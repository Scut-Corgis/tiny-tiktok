package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"
	"log"
	"strconv"
)

type CommentServiceImpl struct {
	UserService
}

func (CommentServiceImpl) QueryCommentsByVideoId(id int64) []dao.Comment {
	commentsList, err := dao.QueryCommentsByVideoId(id)
	if err != nil {
		log.Println("error:", err.Error())
		log.Println("Video not found!")
		return commentsList
	}
	log.Println("Query comments successfully!")
	return commentsList
}

func (CommentServiceImpl) InsertComment(comment *dao.Comment) bool {
	flag := dao.InsertComment(comment)
	if flag == false {
		log.Println("Insert comment failed!")
		return flag
	}

	// 注入redis
	redisCommentKey := util.Relation_Comment_Key + strconv.FormatInt(comment.Id, 10)
	redis.RedisDb.SAdd(redis.Ctx, redisCommentKey, comment.VideoId)
	// 更新过期时间
	redis.RedisDb.Expire(redis.Ctx, redisCommentKey, util.Relation_Follow_TTL)
	log.Println("Insert comment successfully!")
	return flag
}

func (CommentServiceImpl) DeleteComment(id int64) bool {
	flag := dao.DeleteComment(id)
	if flag == false {
		log.Println("Delete comment failed!")
		return flag
	}
	log.Println("Delete comment successfully!")
	return flag
}
