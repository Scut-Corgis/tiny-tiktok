package service

import (
	"log"
	"strconv"
	"strings"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/rabbitmq"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"
)

type LikeServiceImpl struct {
	UserService //出错
}

func (LikeServiceImpl) Like(userId int64, videoId int64) error {
	//如果点赞关系已经存在，则直接返回
	//IsFavorite := dao.JudgeIsFavorite(userId, videoId)
	//if IsFavorite {
	//	return nil
	//}
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)

	message := strings.Builder{}
	message.WriteString(strUserId)
	message.WriteString(":")
	message.WriteString(strVideoId)

	//如果点赞的用户id在redis缓存中，那就把被点赞的视频id添加到key为用户id的set中,并且把点赞数据通过rbt发送给数据库，在数据库中添加
	if n, err := redis.RedisDbLikeUserIdVideoId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
		if err != nil {
			log.Println("redis 查询失败")
			return err
		}
		if _, err := redis.RedisDbLikeUserIdVideoId.SAdd(redis.Ctx, strUserId, strVideoId).Result(); err != nil {
			log.Println("redis 添加失败")
			return err
		} else { //只有缓存中添加正确，才可以往数据库中添加
			rabbitmq.RabbitMQLikeAdd.Producer(message.String())
		}
	} else { //如果点赞的用户id不在redis缓存中，
		//在缓存中新建一个useridkey
		if _, err := redis.RedisDbLikeUserIdVideoId.SAdd(redis.Ctx, strUserId, -1).Result(); err != nil {
			log.Println("缓存创建key失败！")
			redis.RedisDbLikeUserIdVideoId.Del(redis.Ctx, strUserId)
			return err
		}
		//设置过期时间
		if _, err := redis.RedisDbLikeUserIdVideoId.Expire(redis.Ctx, strUserId, util.Day).Result(); err != nil { //这个过期时间是随便给的，下面需要想想具体给多少
			log.Println("缓存过期时间设置失败")
			redis.RedisDbLikeUserIdVideoId.Del(redis.Ctx, strUserId)
		}
		//把数据库中的当前用户点赞的videoId全部添加到缓存中
		videoIdList, err1 := dao.GetLikeVideoIdList(userId)
		if err1 != nil {
			return err1
		}
		for _, videoId := range videoIdList {
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RedisDbLikeUserIdVideoId.SAdd(redis.Ctx, strUserId, videoId).Result(); err != nil {
				log.Println("videoId添加缓存失败")
				redis.RedisDbLikeUserIdVideoId.Del(redis.Ctx, strUserId)
				return err
			}
		}
		//把该次点赞的videoId添加到缓存中
		if _, err := redis.RedisDbLikeUserIdVideoId.SAdd(redis.Ctx, strUserId, strVideoId).Result(); err != nil {
			log.Println("videoId添加缓存失败")
			redis.RedisDbLikeUserIdVideoId.Del(redis.Ctx, strUserId)
			return err
		} else { //只有缓存中添加正确，才可以往数据库中添加
			rabbitmq.RabbitMQLikeAdd.Producer(message.String())
		}
	}
	return nil
}

func (LikeServiceImpl) Unlike(userId int64, videoId int64) error {
	return dao.DeleteLike(userId, videoId)
}

func (LikeServiceImpl) GetVideo(videoId int64, userId int64) dao.VideoDetail {
	vsi := VideoServiceImpl{}
	VideoDetail, _ := vsi.QueryVideoDetailByVideoId(videoId, userId)
	return VideoDetail
}

func (LikeServiceImpl) GetLikeLists(userId int64) ([]dao.VideoDetail, error) {
	var result []dao.VideoDetail

	videoIdList, err := dao.GetLikeVideoIdList(userId)
	if err != nil {
		log.Println("获取用户点赞视频Id出错！")
	}
	//fmt.Println(videoIdList)
	for _, videoId := range videoIdList {
		video := GetVideo(videoId, userId)
		result = append(result, video)
	}
	return result, nil
}

func (LikeServiceImpl) IsFavourite(videoId int64, userId int64) (bool, error) {
	likeinfon, err := dao.GetLikInfo(videoId, userId)
	if err != nil {
		log.Println("get likeInfo failed")
		return false, nil
	}
	if likeinfon.UserId != userId || likeinfon.VideoId != videoId {
		return false, nil
	}
	return true, nil
}

func (LikeServiceImpl) FavouriteCount(videoId int64) (int64, error) {
	return dao.GetLikeCountByVideoId(videoId)
}
