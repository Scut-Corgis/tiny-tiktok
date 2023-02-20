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
			log.Println("redis 查询失败")
			return err
		}
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, strVideoId).Result(); err != nil {
			log.Println("redis 添加失败")
			return err
		} else { //只有缓存中添加正确，才可以往数据库中添加
			rabbitmq.RabbitMQLikeAdd.Producer(message.String())
		}
		//list, _ := redis.RedisDb.SMembers(redis.Ctx, key_userId).Result()
		//log.Println(key_userId, list)
	} else { //如果点赞的用户id不在redis缓存中，
		//在缓存中新建一个useridkey
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, util.MyDefault).Result(); err != nil {
			log.Println("缓存创建key失败！")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return err
		}
		//设置过期时间
		if _, err := redis.RedisDb.Expire(redis.Ctx, key_userId, util.Day).Result(); err != nil { //这个过期时间是随便给的，下面需要想想具体给多少
			log.Println("缓存过期时间设置失败")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return err
		}
		//把数据库中的当前用户点赞的videoId全部添加到缓存中
		videoIdList, err1 := dao.GetLikeVideoIdList(userId)
		if err1 != nil {
			log.Println("获取点赞视频id列表")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return err1
		}
		for _, videoId := range videoIdList {
			strvideoId := strconv.FormatInt(videoId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, strvideoId).Result(); err != nil {
				log.Println("videoId添加缓存失败")
				redis.RedisDb.Del(redis.Ctx, key_userId)
				return err
			}
		}
		//把该次点赞的videoId添加到缓存中
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, strVideoId).Result(); err != nil {
			log.Println("videoId添加缓存失败")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return err
		} else { //只有缓存中添加正确，才可以往数据库中添加
			rabbitmq.RabbitMQLikeAdd.Producer(message.String())
		}
		//list, _ := redis.RedisDb.SMembers(redis.Ctx, key_userId).Result()
		//log.Println(key_userId, list)
	}
	//查看该次点赞的strVideoId是否在缓存中
	if n, err := redis.RedisDb.Exists(redis.Ctx, key_videoId).Result(); n > 0 {
		if err != nil {
			log.Println("redis 查询失败")
			return err
		}
		//如果在缓存中，则把当前点赞的strUserId添加到key为strVideoId的set中
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, strUserId).Result(); err != nil {
			log.Println("redis 添加失败")
			return err
		}
		//list, _ := redis.RedisDb.SMembers(redis.Ctx, key_videoId).Result()
		//log.Println(key_videoId, list)
	} else { //如果在缓存中
		//先添加一个默认值
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, util.MyDefault).Result(); err != nil {
			log.Println("缓存创建key失败！")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err
		}
		//设置过期时间
		if _, err := redis.RedisDb.Expire(redis.Ctx, key_videoId, util.Day).Result(); err != nil {
			log.Println("缓存过期时间设置失败")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err
		}
		//把数据库中给当前视频的点赞的userId全部添加到缓存中
		userIdList, err1 := dao.GetLikeUserIdList(videoId)
		if err1 != nil {
			log.Println("获取视频id点赞用户列表失败")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err1
		}
		for _, userId := range userIdList {
			struserid := strconv.FormatInt(userId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, struserid).Result(); err != nil {
				log.Println("videoId添加缓存失败")
				redis.RedisDb.Del(redis.Ctx, key_videoId)
				return err
			}
		}
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, strUserId).Result(); err != nil {
			log.Println("videoId添加缓存失败")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err
		}
		//list, _ := redis.RedisDb.SMembers(redis.Ctx, key_videoId).Result()
		//log.Println(key_videoId, list)
	}
	return nil
}

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
			log.Println("redis 查询失败")
			return err
		}
		if _, err := redis.RedisDb.SRem(redis.Ctx, key_userId, strVideoId).Result(); err != nil { //把该videoid从缓存中删除
			log.Println("redis 取消点赞删除缓存失败")
		} else { //只有缓存操作成功，数据库才可以操作
			rabbitmq.RabbitMQLikeDel.Producer(message.String())
		}
		//list, _ := redis.RedisDb.SMembers(redis.Ctx, key_userId).Result()
		//log.Println(key_userId, list)
	} else { //如果取消点赞的用户id不在redis缓存中 过期了
		//在缓存中新建一个useridkey
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, util.MyDefault).Result(); err != nil {
			log.Println("缓存创建key失败！")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return err
		}
		//设置过期时间
		if _, err := redis.RedisDb.Expire(redis.Ctx, key_userId, util.Day).Result(); err != nil { //这个过期时间是随便给的，下面需要想想具体给多少
			log.Println("缓存过期时间设置失败")
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
				log.Println("videoId添加缓存失败")
				redis.RedisDb.Del(redis.Ctx, key_userId)
				return err
			}
		}
		//把该次取消点赞的videoId移除缓存
		if _, err := redis.RedisDb.SRem(redis.Ctx, key_userId, strVideoId).Result(); err != nil {
			log.Println("videoId添加缓存失败")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return err
		} else {
			rabbitmq.RabbitMQLikeDel.Producer(message.String())
		}
		//list, _ := redis.RedisDb.SMembers(redis.Ctx, key_userId).Result()
		//log.Println(key_userId, list)
	}

	//查看该次取消点赞的strVideoId是否在缓存中
	if n, err := redis.RedisDb.Exists(redis.Ctx, key_videoId).Result(); n > 0 {
		if err != nil {
			log.Println("redis 查询失败")
			return err
		}
		//如果在缓存中，则把当前点赞的strUserId移除key为strVideoId的set中
		if _, err := redis.RedisDb.SRem(redis.Ctx, key_videoId, strUserId).Result(); err != nil {
			log.Println("redis 移除失败")
			return err
		}
		//list, _ := redis.RedisDb.SMembers(redis.Ctx, key_videoId).Result()
		//log.Println(key_videoId, list)
	} else { //如果不在缓存中
		//先添加一个默认值
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, util.MyDefault).Result(); err != nil {
			log.Println("缓存创建key失败！")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err
		}
		//设置过期时间
		if _, err := redis.RedisDb.Expire(redis.Ctx, key_videoId, util.Day).Result(); err != nil {
			log.Println("缓存过期时间设置失败")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err
		}
		//把数据库中给当前视频的点赞的userId全部添加到缓存中
		userIdList, err1 := dao.GetLikeUserIdList(videoId)
		if err1 != nil {
			log.Println("获取视频id点赞用户列表失败")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err1
		}
		for _, userId := range userIdList {
			struserid := strconv.FormatInt(userId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, struserid).Result(); err != nil {
				log.Println("videoId添加缓存失败")
				redis.RedisDb.Del(redis.Ctx, key_videoId)
				return err
			}
		}
		if _, err := redis.RedisDb.SRem(redis.Ctx, key_videoId, strUserId).Result(); err != nil { //Srem 命令用于移除集合中的一个或多个成员元素，不存在的成员元素会被忽略
			log.Println("videoId移除缓存失败")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return err
		}
		//list, _ := redis.RedisDb.SMembers(redis.Ctx, key_videoId).Result()
		//log.Println(key_videoId, list)
	}
	return nil
}

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
	//如果没有获取这个video_id的视频，视频可能被删除了,打印异常,并且不加入此视频  ???

	//将视频信息类型对象添加到集合中去
	*result = append(*result, video)
}

// func (LikeServiceImpl) GetLikeLists(userId int64) ([]dao.VideoDetail, error) {
// 	var result []dao.VideoDetail

// 	videoIdList, err := dao.GetLikeVideoIdList(userId)
// 	if err != nil {
// 		log.Println("获取用户点赞视频Id出错！")
// 	}

// 	for _, videoId := range videoIdList {
// 		video := GetVideo(videoId, userId)
// 		result = append(result, video)
// 	}
// 	return result, nil
// }

func (like LikeServiceImpl) GetLikeLists(userId int64) ([]dao.VideoDetail, error) {
	key_userId := util.Like_User_Key + strconv.FormatInt(userId, 10)
	var result []dao.VideoDetail
	if n, err := redis.RedisDb.Exists(redis.Ctx, key_userId).Result(); n > 0 {
		if err != nil {
			log.Println("查询失败")
			return nil, err
		}
		videoIdList, err1 := redis.RedisDb.SMembers(redis.Ctx, key_userId).Result()
		log.Println(key_userId, videoIdList)
		if err1 != nil {
			log.Println("查询失败")
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
			log.Println("添加缓存失败")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return nil, err
		}
		//设置有效期
		if _, err := redis.RedisDb.Expire(redis.Ctx, key_userId, util.Day).Result(); err != nil {
			log.Println("缓存过期时间设置失败")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return nil, err
		}
		//把数据库中的视频id添加到缓存中
		//把数据库中的当前用户点赞的videoId全部添加到缓存中
		videoIdList, err1 := dao.GetLikeVideoIdList(userId)
		if err1 != nil {
			log.Println("获取点赞视频id列表")
			redis.RedisDb.Del(redis.Ctx, key_userId)
			return nil, err1
		}
		for _, videoId := range videoIdList {
			strvideoId := strconv.FormatInt(videoId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, strvideoId).Result(); err != nil {
				log.Println("videoId添加缓存失败")
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

//	func (like LikeServiceImpl) IsFavourite(videoId int64, userId int64) (bool, error) {
//		//key_userId := util.Like_User_Key + strconv.FormatInt(userId, 10)
//		//key_videoId := util.Like_Video_key + strconv.FormatInt(videoId, 10)
//		likeinfon, err := dao.GetLikInfo(videoId, userId)
//		if err != nil {
//			log.Println("get likeInfo failed")
//			return false, nil
//		}
//		if likeinfon.UserId != userId || likeinfon.VideoId != videoId {
//			return false, nil
//		}
//		return true, nil
//	}
func (like LikeServiceImpl) IsLike(videoId int64, userId int64) (bool, error) {
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)
	key_userId := util.Like_User_Key + strconv.FormatInt(userId, 10)
	key_videoId := util.Like_Video_key + strconv.FormatInt(videoId, 10)

	if n, err := redis.RedisDb.Exists(redis.Ctx, key_userId).Result(); n > 0 { //如果key_userId存在缓存中
		if err != nil {
			log.Println("redis 查询失败")
			return false, err
		}
		isLike, err := redis.RedisDb.SIsMember(redis.Ctx, key_userId, strVideoId).Result()
		if err != nil {
			log.Println("redis 查询失败")
			return false, err
		}
		return isLike, nil
	} else { //如果key_userId不存在缓存中，查询key_videoId是否在缓存中
		if n, err := redis.RedisDb.Exists(redis.Ctx, key_videoId).Result(); n > 0 { //如果key_userId存在缓存中
			if err != nil {
				log.Println("redis 查询失败")
				return false, err
			}
			isLike, err := redis.RedisDb.SIsMember(redis.Ctx, key_videoId, strUserId).Result()
			if err != nil {
				log.Println("redis 查询失败")
				return false, err
			}
			return isLike, nil
		} else { //如果key_userId不存在缓存中 那么key_userId key_videoId都不存在缓存当中，把数据库中的用户userId点赞的视频id添加到key_userId中
			if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, util.MyDefault).Result(); err != nil {
				log.Println("添加缓存失败")
				redis.RedisDb.Del(redis.Ctx, key_userId)
				return false, err
			}
			//设置有效期
			if _, err := redis.RedisDb.Expire(redis.Ctx, key_userId, util.Day).Result(); err != nil {
				log.Println("缓存过期时间设置失败")
				redis.RedisDb.Del(redis.Ctx, key_userId)
				return false, err
			}
			//把数据库中的视频id添加到缓存中
			//把数据库中的当前用户点赞的videoId全部添加到缓存中
			videoIdList, err1 := dao.GetLikeVideoIdList(userId)
			if err1 != nil {
				log.Println("获取点赞视频id列表")
				redis.RedisDb.Del(redis.Ctx, key_userId)
				return false, err1
			}
			for _, videoId := range videoIdList {
				strvideoId := strconv.FormatInt(videoId, 10)
				//如果出现一次不对的就把这个键值删除
				if _, err := redis.RedisDb.SAdd(redis.Ctx, key_userId, strvideoId).Result(); err != nil {
					log.Println("videoId添加缓存失败")
					redis.RedisDb.Del(redis.Ctx, key_userId)
					return false, err
				}
			}
			isLike, err := redis.RedisDb.SIsMember(redis.Ctx, key_userId, strVideoId).Result()
			if err != nil {
				log.Println("redis 查询失败")
				return false, err
			}

			return isLike, nil
		}
	}
}

func (like LikeServiceImpl) LikeCount(videoId int64) (int64, error) {
	key_videoId := util.Like_Video_key + strconv.FormatInt(videoId, 10)
	var result int64
	result = 1
	//如果键值key_videoId在缓存中
	if n, err := redis.RedisDb.Exists(redis.Ctx, key_videoId).Result(); n > 0 {
		if err != nil {
			log.Println("redis 查询失败")
			return -1, err
		}
		result, err1 := redis.RedisDb.SCard(redis.Ctx, key_videoId).Result()
		if err1 != nil {
			log.Println("redis 查询失败")
			return -1, err
		}

		return result - 1, nil
	} else { //如果键值key_videoId不在缓存中
		if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, util.MyDefault).Result(); err != nil {
			log.Println("添加缓存失败")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return -1, err
		}
		//设置有效期
		if _, err := redis.RedisDb.Expire(redis.Ctx, key_videoId, util.Day).Result(); err != nil {
			log.Println("缓存过期时间设置失败")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return -1, err
		}
		//把数据库中的用户id添加到缓存中
		//把数据库中的点赞该视频的userid全部添加到缓存中
		userIdList, err1 := dao.GetLikeUserIdList(videoId)
		if err1 != nil {
			log.Println("获取点赞视频id列表")
			redis.RedisDb.Del(redis.Ctx, key_videoId)
			return -1, err1
		}
		for _, videoId := range userIdList {
			struserId := strconv.FormatInt(videoId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RedisDb.SAdd(redis.Ctx, key_videoId, struserId).Result(); err != nil {
				log.Println("videoId添加缓存失败")
				redis.RedisDb.Del(redis.Ctx, key_videoId)
				return -1, err
			}
		}
		result = int64(len(userIdList))

		return result, nil
	}
}
