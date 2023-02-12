package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"
	"log"
	"strconv"
)

type CommentServiceImpl struct {
	UserServiceImpl
	VideoServiceImpl
}

// QueryCommentsByVideoId 根据视频id获取comment列表
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

// PostComment 发表评论
func (CommentServiceImpl) PostComment(comment dao.Comment) (int32, string) {
	flag1 := dao.InsertComment(&comment)
	if flag1 == false {
		return 1, "Post comment failed!"
	}
	flag2 := InsertRedis(comment.VideoId, comment.Id) // 添加redis缓存
	if !flag2 {
		log.Println("Insert redis failed!")
	}
	return 0, "Post comment successfully!"
}

// DeleteComment 根据评论id删除评论
func (CommentServiceImpl) DeleteComment(id int64) (int32, string) {
	// 先查询redis缓存，若有则删除，再删除数据库记录；若无则直接删除数据库记录
	redisCommentKey := util.Relation_Comment_Key + strconv.FormatInt(id, 10)
	if err0 := redis.RedisDbCommentIdVideoId.Exists(redis.Ctx, redisCommentKey).Err(); err0 != nil {
		log.Println(err0.Error())
	}
	videoId, err1 := redis.RedisDbCommentIdVideoId.Get(redis.Ctx, redisCommentKey).Result()
	redisVideoKey := util.Relation_Video_Key + videoId
	if err1 != nil {
		log.Println(err1.Error())
	}
	// 删除redis缓存
	if err2 := redis.RedisDbCommentIdVideoId.Del(redis.Ctx, redisCommentKey).Err(); err2 != nil {
		log.Println(err2.Error())
	}
	if err3 := redis.RedisDbVideoIdCommentId.Del(redis.Ctx, redisVideoKey).Err(); err3 != nil {
		log.Println(err3.Error())
	}
	log.Println("Delete record in redis successfully!")
	flag := dao.DeleteComment(id)
	if flag == false {
		return 1, "Delete comment failed!"
	}
	return 0, "Delete comment successfully!"
}

func InsertRedis(video_id int64, comment_id int64) bool {
	// 更新RedisDbVideoIdCommentId
	redisVideoKey := util.Relation_Video_Key + strconv.FormatInt(video_id, 10)
	if err := redis.RedisDbVideoIdCommentId.SAdd(redis.Ctx, redisVideoKey, comment_id).Err(); err != nil {
		log.Println("Insert RedisDbCommentIdVideoId failed!")
		redis.RedisDbVideoIdCommentId.Del(redis.Ctx, redisVideoKey) // 缓存失败就删除key
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
