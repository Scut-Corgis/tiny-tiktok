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

func (CommentServiceImpl) PostComment(comment *dao.Comment) (int32, string) {
	flag1 := dao.InsertComment(comment)
	if flag1 == false {
		return 1, "Insert comment failed!"
	}
	flag2 := InsertRedis(comment.VideoId, comment.Id) // 添加redis缓存
	if !flag2 {
		log.Println("Insert redis failed!")
	}
	return 0, "Post comment successfully!"
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

func InsertRedis(video_id int64, comment_id int64) bool {
	// 更新RedisDbVideoIdCommentId
	redisVideoKey := util.Relation_Video_Key + strconv.FormatInt(video_id, 10)
	if err := redis.RedisDbVideoIdCommentId.SAdd(redis.Ctx, redisVideoKey, video_id).Err(); err != nil {
		log.Println("Insert RedisDbCommentIdVideoId failed!")
		redis.RedisDbVideoIdCommentId.Del(redis.Ctx) // 缓存失败就删除key
		return false
	}
	redis.RedisDbVideoIdCommentId.Expire(redis.Ctx, redisVideoKey, util.Relation_Follow_TTL) // 缓存成功更新过期时间
	// 更新RedisDbCommentIdVideoId
	redisCommentKey := util.Relation_Comment_Key + strconv.FormatInt(comment_id, 10)
	if err := redis.RedisDbCommentIdVideoId.Set(redis.Ctx, redisCommentKey, video_id, 0).Err(); err != nil {
		log.Println("Insert RedisDbVideoIdCommentId failed!")
		return false
	}
	return true
}
