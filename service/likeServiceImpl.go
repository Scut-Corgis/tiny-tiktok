package service

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/rabbitmq"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"
)

type LikeServiceImpl struct {
	UserService
	VideoService
}

/*点赞*/
func (like LikeServiceImpl) Like(userId int64, videoId int64) error {
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)

	key_userId := util.Like_User_Key + strconv.FormatInt(userId, 10)
	key_videoId := util.Like_Video_key + strconv.FormatInt(videoId, 10)

	message := strings.Builder{}
	message.WriteString(strUserId)
	message.WriteString(":")
	message.WriteString(strVideoId)

	//如果点赞的用户id在redis缓存中，那就把被点赞的视频id添加到key为用户id的set中,并且把点赞数据通过rbt发送给数据库，在数据库中添加
	if n, err := redis.RedisDb.Exists(redis.Ctx, key_userId).Result(); n > 0 {
		if err != nil {
			log.Println("Redis query failed")
			return err
		}
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, strVideoId).Result(); err != nil {
			log.Println("Redis add failed")
			return err
		} else { //只有缓存中添加正确，才可以往数据库中添加
			rabbitmq.RabbitMQLikeAdd.Producer(message.String())
		}
	} else { //如果点赞的用户id不在redis缓存中，
		//在缓存中新建一个useridkey
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, util.MyDefault).Result(); err != nil {
			log.Println("Cache creation key failed!")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return err
		}
		//设置过期时间
		if _, err := redis.RedisDb.Expire(redis.Ctx, key_userId, redis.RandomTime()).Result(); err != nil {
			log.Println("Failed to set cache expiration time")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return err
		}
		//把数据库中的当前用户点赞的videoId全部添加到缓存中
		videoIdList, err1 := dao.GetLikeVideoIdList(userId)
		if err1 != nil {
			log.Println("Failed to get the likes video id list")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return err1
		}
		for _, videoId := range videoIdList {
			strvideoId := strconv.FormatInt(videoId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, strvideoId).Result(); err != nil {
				log.Println("Failed to add cache for videoId")
				redis.RedisDb.Del(redis.Ctx, key_userId)
				return err
			}
		}
		//把该次点赞的videoId添加到缓存中
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, strVideoId).Result(); err != nil {
			log.Println("Failed to add cache for videoId")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return err
		} else { //只有缓存中添加正确，才可以往数据库中添加
			rabbitmq.RabbitMQLikeAdd.Producer(message.String())
		}
	}
	//查看该次点赞的strVideoId是否在缓存中
	if n, err := redis.RedisDb.Exists(redis.Ctx, key_videoId).Result(); n > 0 {
		if err != nil {
			log.Println("Redis query failed")
			return err
		}
		//如果在缓存中，则把当前点赞的strUserId添加到key为strVideoId的set中
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, strUserId).Result(); err != nil {
			log.Println("Redis add failed")
			return err
		}
	} else { //如果在缓存中
		//先添加一个默认值
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, util.MyDefault).Result(); err != nil {
			log.Println("Cache creation key failed!")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err
		}
		//设置过期时间
		if _, err := redis.RedisDb.Expire(redis.Ctx, key_videoId, redis.RandomTime()).Result(); err != nil {
			log.Println("Failed to set cache expiration time")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err
		}
		//把数据库中给当前视频的点赞的userId全部添加到缓存中
		userIdList, err1 := dao.GetLikeUserIdList(videoId)
		if err1 != nil {
			log.Println("Failed to get video id like user list")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err1
		}
		for _, userId := range userIdList {
			struserid := strconv.FormatInt(userId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, struserid).Result(); err != nil {
				log.Println("Failed to add cache for videoId")
				redis.RedisDb.Del(redis.Ctx, key_videoId)
				return err
			}
		}
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, strUserId).Result(); err != nil {
			log.Println("Failed to add cache for videoId")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err
		}
	}
	return nil
}

/*取消点赞*/
func (like LikeServiceImpl) Unlike(userId int64, videoId int64) error {
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)

	key_userId := util.Like_User_Key + strconv.FormatInt(userId, 10)
	key_videoId := util.Like_Video_key + strconv.FormatInt(videoId, 10)

	message := strings.Builder{}
	message.WriteString(strUserId)
	message.WriteString(":")
	message.WriteString(strVideoId)

	if n, err := redis.RedisDb.Exists(redis.Ctx, key_userId).Result(); n > 0 { //如果取消点赞的用户id在redis缓存中
		if err != nil {
			log.Println("Redis query failed")
			return err
		}
		if _, err := redis.RedisDb.SRem(redis.Ctx, key_userId, strVideoId).Result(); err != nil { //把该videoid从缓存中删除
			log.Println("Redis cancel likes delete cache failed")
		} else { //只有缓存操作成功，数据库才可以操作
			rabbitmq.RabbitMQLikeDel.Producer(message.String())
		}
	} else { //如果取消点赞的用户id不在redis缓存中 过期了
		//在缓存中新建一个useridkey
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, util.MyDefault).Result(); err != nil {
			log.Println("Cache creation key failed!")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return err
		}
		//设置过期时间
		if _, err := redis.RedisDb.Expire(redis.Ctx, key_userId, redis.RandomTime()).Result(); err != nil {
			log.Println("Failed to set cache expiration time")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return err
		}
		//把数据库中的当前用户点赞的videoId全部添加到缓存中
		videoIdList, err1 := dao.GetLikeVideoIdList(userId)
		if err1 != nil {
			return err1
		}
		for _, videoId := range videoIdList {
			strvideoId := strconv.FormatInt(videoId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, strvideoId).Result(); err != nil {
				log.Println("Failed to add cache for videoId")
				redis.RedisDb.Del(redis.Ctx, key_userId)
				return err
			}
		}
		//把该次取消点赞的videoId移除缓存
		if _, err := redis.RedisDb.SRem(redis.Ctx, key_userId, strVideoId).Result(); err != nil {
			log.Println("Failed to add cache for videoId")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return err
		} else {
			rabbitmq.RabbitMQLikeDel.Producer(message.String())
		}
	}

	//查看该次取消点赞的strVideoId是否在缓存中
	if n, err := redis.RedisDb.Exists(redis.Ctx, key_videoId).Result(); n > 0 {
		if err != nil {
			log.Println("Redis query failed")
			return err
		}
		//如果在缓存中，则把当前点赞的strUserId移除key为strVideoId的set中
		if _, err := redis.RedisDb.SRem(redis.Ctx, key_videoId, strUserId).Result(); err != nil {
			log.Println("Redis removal failed")
			return err
		}
	} else { //如果不在缓存中
		//先添加一个默认值
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, util.MyDefault).Result(); err != nil {
			log.Println("Cache creation key failed!")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err
		}
		//设置过期时间
		if _, err := redis.RedisDb.Expire(redis.Ctx, key_videoId, redis.RandomTime()).Result(); err != nil {
			log.Println("Failed to set cache expiration time")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err
		}
		//把数据库中给当前视频的点赞的userId全部添加到缓存中
		userIdList, err1 := dao.GetLikeUserIdList(videoId)
		if err1 != nil {
			log.Println("Failed to get video id like user list")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err1
		}
		for _, userId := range userIdList {
			struserid := strconv.FormatInt(userId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, struserid).Result(); err != nil {
				log.Println("Failed to add cache for videoId")
				redis.RedisDb.Del(redis.Ctx, key_videoId)
				return err
			}
		}
		if _, err := redis.RedisDb.SRem(redis.Ctx, key_videoId, strUserId).Result(); err != nil { //Srem 命令用于移除集合中的一个或多个成员元素，不存在的成员元素会被忽略
			log.Println("VideoId removal cache failed")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err
		}
	}
	return nil
}

/*根据videoId userId获取视频的详细信息，并添加到result中*/
func (like LikeServiceImpl) GetVideo(videoId int64, userId int64, result *[]dao.VideoDetail, wg *sync.WaitGroup) {
	defer wg.Done()

	//根据userId和videoId查询数据库中视频信息
	videoDetail, _ := VideoServiceImpl{}.QueryVideoDetailByVideoId(videoId, userId)
	video := dao.VideoDetail{
		Id:            videoDetail.Id,
		Author:        dao.UserResp(videoDetail.Author),
		PlayUrl:       videoDetail.PlayUrl,
		CoverUrl:      videoDetail.CoverUrl,
		FavoriteCount: videoDetail.FavoriteCount,
		CommentCount:  videoDetail.CommentCount,
		IsFavorite:    videoDetail.IsFavorite,
		Title:         videoDetail.Title,
	}
	//如果没有获取这个video_id的视频，视频可能被删除了,打印异常,并且不加入此视频
	//将视频信息类型对象添加到集合中去
	*result = append(*result, video)
}

/*获取点赞列表*/
func (like LikeServiceImpl) GetLikeLists(userId int64) ([]dao.VideoDetail, error) {
	key_userId := util.Like_User_Key + strconv.FormatInt(userId, 10)
	var result []dao.VideoDetail
	if n, err := redis.RedisDb.Exists(redis.Ctx, key_userId).Result(); n > 0 {
		if err != nil {
			log.Println("redis Query failed")
			return nil, err
		}
		videoIdList, err1 := redis.RedisDb.SMembers(redis.Ctx, key_userId).Result()
		log.Println(key_userId, videoIdList)
		if err1 != nil {
			log.Println("redis Query failed")
			return nil, err
		}

		videoIdListLen := len(videoIdList)
		log.Println(videoIdListLen)
		if videoIdListLen == 0 {
			return result, nil
		}

		var wg sync.WaitGroup
		wg.Add(videoIdListLen - 1)
		for i := 0; i < videoIdListLen; i++ {
			videoId, _ := strconv.ParseInt(videoIdList[i], 10, 64)
			log.Println(videoId)
			if videoId == util.MyDefault {
				continue
			}
			go like.GetVideo(videoId, userId, &result, &wg)
		}
		wg.Wait()
		return result, nil
	} else { //如果key_userId不存在缓存中，需要把数据库中的信息添加到缓存中
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, util.MyDefault).Result(); err != nil {
			log.Println("Failed to add cache")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return nil, err
		}
		//设置有效期
		if _, err := redis.RedisDb.Expire(redis.Ctx, key_userId, redis.RandomTime()).Result(); err != nil {
			log.Println("Failed to set cache expiration time")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return nil, err
		}
		//把数据库中的视频id添加到缓存中
		//把数据库中的当前用户点赞的videoId全部添加到缓存中
		videoIdList, err1 := dao.GetLikeVideoIdList(userId)
		if err1 != nil {
			log.Println("Failed to get the likes video id list")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return nil, err1
		}
		for _, videoId := range videoIdList {
			strvideoId := strconv.FormatInt(videoId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, strvideoId).Result(); err != nil {
				log.Println("Failed to add cache for videoId")
				redis.RedisDb.Del(redis.Ctx, key_userId)
				return nil, err
			}
		}

		videoIdListLen := len(videoIdList)
		if videoIdListLen == 0 {
			return result, nil
		}
		//从缓存中把点赞列表获取
		var wg sync.WaitGroup
		wg.Add(videoIdListLen)
		for i := 0; i < videoIdListLen; i++ {
			go like.GetVideo(videoIdList[i], userId, &result, &wg)
		}
		wg.Wait()
		return result, nil
	}

}

/*判断用户userId是否点赞视频videoId*/
func (like LikeServiceImpl) IsLike(videoId int64, userId int64) (bool, error) {
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)
	key_userId := util.Like_User_Key + strconv.FormatInt(userId, 10)
	key_videoId := util.Like_Video_key + strconv.FormatInt(videoId, 10)

	if n, err := redis.RedisDb.Exists(redis.Ctx, key_userId).Result(); n > 0 { //如果key_userId存在缓存中
		if err != nil {
			log.Println("Redis query failed")
			return false, err
		}
		isLike, err := redis.RedisDb.SIsMember(redis.Ctx, key_userId, strVideoId).Result()
		if err != nil {
			log.Println("Redis query failed")
			return false, err
		}
		return isLike, nil
	} else { //如果key_userId不存在缓存中，查询key_videoId是否在缓存中
		if n, err := redis.RedisDb.Exists(redis.Ctx, key_videoId).Result(); n > 0 { //如果key_userId存在缓存中
			if err != nil {
				log.Println("Redis query failed")
				return false, err
			}
			isLike, err := redis.RedisDb.SIsMember(redis.Ctx, key_videoId, strUserId).Result()
			if err != nil {
				log.Println("Redis query failed")
				return false, err
			}
			return isLike, nil
		} else { //如果key_userId不存在缓存中 那么key_userId key_videoId都不存在缓存当中，把数据库中的用户userId点赞的视频id添加到key_userId中
			if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, util.MyDefault).Result(); err != nil {
				log.Println("Redis Failed to add cache")
				redis.RedisDb.Del(redis.Ctx, key_userId)
				return false, err
			}
			//设置有效期
			if _, err := redis.RedisDb.Expire(redis.Ctx, key_userId, redis.RandomTime()).Result(); err != nil {
				log.Println("Failed to set cache expiration time")
				redis.RedisDb.Del(redis.Ctx, key_userId)
				return false, err
			}
			//把数据库中的视频id添加到缓存中
			//把数据库中的当前用户点赞的videoId全部添加到缓存中
			videoIdList, err1 := dao.GetLikeVideoIdList(userId)
			if err1 != nil {
				log.Println("Failed to get the likes video id list")
				redis.RedisDb.Del(redis.Ctx, key_userId)
				return false, err1
			}
			for _, videoId := range videoIdList {
				strvideoId := strconv.FormatInt(videoId, 10)
				//如果出现一次不对的就把这个键值删除
				if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, strvideoId).Result(); err != nil {
					log.Println("Failed to add cache for videoId")
					redis.RedisDb.Del(redis.Ctx, key_userId)
					return false, err
				}
			}
			isLike, err := redis.RedisDb.SIsMember(redis.Ctx, key_userId, strVideoId).Result()
			if err != nil {
				log.Println("Redis query failed")
				return false, err
			}

			return isLike, nil
		}
	}
}

/*获取视频videoId的点赞数*/
func (like LikeServiceImpl) LikeCount(videoId int64) (int64, error) {
	key_videoId := util.Like_Video_key + strconv.FormatInt(videoId, 10)
	var result int64
	result = 1
	//如果键值key_videoId在缓存中
	if n, err := redis.RedisDb.Exists(redis.Ctx, key_videoId).Result(); n > 0 {
		if err != nil {
			log.Println("Redis query failed")
			return -1, err
		}
		result, err1 := redis.RedisDb.SCard(redis.Ctx, key_videoId).Result()
		if err1 != nil {
			log.Println("Redis query failed")
			return -1, err
		}

		return result - 1, nil
	} else { //如果键值key_videoId不在缓存中
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, util.MyDefault).Result(); err != nil {
			log.Println("Failed to add cache")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return -1, err
		}
		//设置有效期
		if _, err := redis.RedisDb.Expire(redis.Ctx, key_videoId, redis.RandomTime()).Result(); err != nil {
			log.Println("Failed to set cache expiration time")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return -1, err
		}
		//把数据库中的用户id添加到缓存中
		//把数据库中的点赞该视频的userid全部添加到缓存中
		userIdList, err1 := dao.GetLikeUserIdList(videoId)
		if err1 != nil {
			log.Println("Get likes video id list")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return -1, err1
		}
		for _, videoId := range userIdList {
			struserId := strconv.FormatInt(videoId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, struserId).Result(); err != nil {
				log.Println("Failed to add cache for videoId")
				redis.RedisDb.Del(redis.Ctx, key_videoId)
				return -1, err
			}
		}
		result = int64(len(userIdList))

		return result, nil
	}
}
func (like LikeServiceImpl) addVideoLikeCount(videoId int64, sum *int64, wg *sync.WaitGroup) {
	defer wg.Done()

	count, err := like.LikeCount(videoId)
	if err != nil {
		//如果有错误，输出错误信息，并不加入该视频点赞数
		log.Println("video query likes failed")
		return
	}
	*sum += count
}

/*获取用户userId的获取的点赞总数*/
func (like LikeServiceImpl) TotalLiked(userId int64) int64 {
	//获取用户userId发布视频的videoId列表
	videoIdList := VideoServiceImpl{}.GetVideoIdListByUserId(userId)
	listlLen := len(videoIdList)
	//videoLikecountList := make([]int64,listlLen)
	var result int64 = 0

	var wg sync.WaitGroup
	wg.Add(listlLen)
	for i := 0; i < listlLen; i++ {
		go like.addVideoLikeCount(videoIdList[i], &result, &wg)
	}
	wg.Wait()

	return result
}

/*获取用户userId喜欢的视频数量*/
func (like LikeServiceImpl) LikeVideoCount(userId int64) (int64, error) {
	key_userId := util.Like_User_Key + strconv.FormatInt(userId, 10)
	//先判断key_userId键值是否在缓存中
	if n, err := redis.RedisDb.Exists(redis.Ctx, key_userId).Result(); n > 0 { //key_userId键值在缓存中
		if err != nil {
			log.Println("redis Query failed")
			return 0, err
		} else { //查询成功
			result, err := redis.RedisDb.SCard(redis.Ctx, key_userId).Result() //获取key_userId键值有几个val
			if err != nil {
				log.Println("redis Query failed")
				return 0, err
			}
			//减去添加的默认值
			return result - 1, nil
		}
	} else { //key_userId键值不在缓存中，需要把MySQL中的数据添加到缓存中
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, util.MyDefault).Result(); err != nil {
			log.Println("Failed to add cache")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return 0, err
		}
		//设置有效期
		if _, err := redis.RedisDb.Expire(redis.Ctx, key_userId, redis.RandomTime()).Result(); err != nil {
			log.Println("Failed to set cache expiration time")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return 0, err
		}
		//把数据库中的当前用户点赞的videoId全部添加到缓存中
		likevideoIdList, err1 := dao.GetLikeVideoIdList(userId)
		if err1 != nil {
			log.Println("Failed to get the likes video id list")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return 0, err1
		}
		for _, videoId := range likevideoIdList {
			strvideoId := strconv.FormatInt(videoId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, strvideoId).Result(); err != nil {
				log.Println("Failed to add cache for videoId")
				redis.RedisDb.Del(redis.Ctx, key_userId)
				return 0, err
			}
		}

		return int64(len(likevideoIdList)), nil
	}
}
