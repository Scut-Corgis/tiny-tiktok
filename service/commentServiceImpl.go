package service

import (
	"log"
	"strconv"

	"github.com/Scut-Corgis/tiny-tiktok/middleware/rabbitmq"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"
)

type CommentServiceImpl struct {
	UserServiceImpl
}

// CountComments 根据视频id统计评论数量
func (CommentServiceImpl) CountComments(id int64) (int64, error) {
	// 先查redis缓存
	redisVideoKey := util.Comment_Video_Key + strconv.FormatInt(id, 10)
	cnt1, err1 := redis.RedisDb.SCard(redis.Ctx, redisVideoKey).Result()
	if err1 != nil {
		log.Println("count from redis error:", err1)
	}
	// 更新过期时间
	redis.RedisDb.Expire(redis.Ctx, redisVideoKey, redis.RandomTime())
	if cnt1 > 0 {
		return cnt1, nil
	}
	// 再查数据库
	cnt2, err2 := dao.CountComments(id)
	if err2 != nil {
		log.Println("count from db error:", err2)
		return 0, err2
	}
	return cnt2, nil
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
func (CommentServiceImpl) PostComment(comment dao.Comment) (int64, int32, string) {
	// 布谷鸟过滤器过滤
	if !redis.CuckooFilterVideoId.Contain([]byte(strconv.FormatInt(comment.VideoId, 10))) {
		return -1, 1, "Video not Found!"
	}
	commentNew, err := dao.InsertComment(comment)
	if err != nil {
		return -1, 1, "Post comment failed!"
	}
	flag := CommentInsertRedis(commentNew.VideoId, commentNew.Id) // 添加redis缓存
	if !flag {
		log.Println("Insert redis failed!")
	}
	return commentNew.Id, 0, "Post comment successfully!"
}

// DeleteComment 根据评论id删除评论
func (CommentServiceImpl) DeleteComment(id int64) (int32, string) {
	// 先查询redis缓存，若有则删除，再删除数据库记录；若无则直接删除数据库记录
	redisCommentKey := util.Comment_Comment_Key + strconv.FormatInt(id, 10)
	err0 := redis.RedisDb.Exists(redis.Ctx, redisCommentKey).Err()
	if err0 != nil {
		log.Println(err0.Error())
	} else {
		videoId, err1 := redis.RedisDb.Get(redis.Ctx, redisCommentKey).Result()
		redisVideoKey := util.Comment_Video_Key + videoId
		if err1 != nil {
			log.Println(err1.Error())
		}
		// 删除redis缓存
		if err2 := redis.RedisDb.Del(redis.Ctx, redisCommentKey).Err(); err2 != nil {
			log.Println(err2.Error())
		}
		if err3 := redis.RedisDb.SRem(redis.Ctx, redisVideoKey, id).Err(); err3 != nil {
			log.Println(err3.Error())
		}
		log.Println("Delete record in redis successfully!")
		rabbitmq.RabbitMQCommentDel.Producer(strconv.FormatInt(id, 10)) // 将id传入消息队列
		return 0, "Delete comment successfully!"
	}
	flag := dao.DeleteComment(id)
	if flag == false {
		return 1, "Delete comment failed!"
	}
	return 0, "Delete comment successfully!"
}

func CommentInsertRedis(videoId int64, commentId int64) bool {
	// 插入键值对 key:video_id value:comment_id
	redisVideoKey := util.Comment_Video_Key + strconv.FormatInt(videoId, 10)
	if _, err := redis.RedisDb.SAdd(redis.Ctx, redisVideoKey, commentId).Result(); err != nil {
		log.Println("Insert key:video_id value:comment_id into redis failed!")
		redis.RedisDb.Del(redis.Ctx, redisVideoKey) // 缓存失败就删除key
		return false
	}
	redis.RedisDb.Expire(redis.Ctx, redisVideoKey, redis.RandomTime()) // 缓存成功更新过期时间
	// 插入键值对 key:comment_id value:video_id
	redisCommentKey := util.Comment_Comment_Key + strconv.FormatInt(commentId, 10)
	if _, err := redis.RedisDb.Set(redis.Ctx, redisCommentKey, videoId, 0).Result(); err != nil {
		log.Println("Insert key:comment_id value:video_id into redis failed!")
		return false
	}
	log.Println("Insert record into redis successfully!")
	return true
}
