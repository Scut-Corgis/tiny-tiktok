package service

import (
	"log"
	"strconv"
	"strings"

	"github.com/Scut-Corgis/tiny-tiktok/config"
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/rabbitmq"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"
)

type RelationServiceImpl struct{}

/*
关注用户
userId 关注 followId
*/
func (RelationServiceImpl) Follow(userId int64, followId int64) (bool, error) {
	// #优化 实例化了过多的对象
	rsi := RelationServiceImpl{}
	isFollowed := rsi.JudgeIsFollowById(userId, followId)
	if isFollowed {
		return false, nil
	}
	// 插入数据库操作放入消息队列
	sb := strings.Builder{}
	sb.WriteString(strconv.FormatInt(userId, 10))
	sb.WriteString(" ")
	sb.WriteString(strconv.FormatInt(followId, 10))
	rabbitmq.RabbitMQRelationAdd.Producer(sb.String())
	log.Println("****relationMQ ADD success!")
	// #优化: 下边注释勿删, 下为只含redis版本，以备测试用
	// err := dao.InsertFollow(userId, followId)
	// if err != nil {
	// 	return false, err
	// }

	// 将查询到的关注关系注入Redis
	redisFollowKey := util.Relation_Follow_Key + strconv.FormatInt(userId, 10)
	redis.RedisDb.SAdd(redis.Ctx, redisFollowKey, followId)
	// 更新过期时间
	redis.RedisDb.Expire(redis.Ctx, redisFollowKey, util.Relation_Follow_TTL)

	// 保证数据一致性：主动使count缓存失效
	rsi.ExpireFollowerCnt(followId)
	rsi.ExpireFollowingCnt(userId)
	return true, nil
}

/*
取关用户
userId 关注 followId
*/
func (RelationServiceImpl) UnFollow(userId int64, followId int64) (bool, error) {
	rsi := RelationServiceImpl{}
	isFollowed := rsi.JudgeIsFollowById(userId, followId)
	// 未关注 isFollowed为false, 返回false，表示userId未关注followId
	if !isFollowed {
		return false, nil
	}
	// 数据库删除操作放入消息队列
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(userId)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(followId)))
	rabbitmq.RabbitMQRelationDel.Producer(sb.String())
	log.Println("****relationMQ Del success!")
	// err := dao.DeleteFollow(userId, followId)
	// if err != nil {
	// 	return false, err
	// }

	// 删除Redis中 redisFollowKey set集合中的followId元素
	redisFollowKey := util.Relation_Follow_Key + strconv.FormatInt(userId, 10)
	redis.RedisDb.SRem(redis.Ctx, redisFollowKey, followId)
	// 更新过期时间
	redis.RedisDb.Expire(redis.Ctx, redisFollowKey, util.Relation_Follow_TTL)
	// 保证数据一致性：主动使count缓存失效
	rsi.ExpireFollowerCnt(followId)
	rsi.ExpireFollowingCnt(userId)
	return true, nil
}

/*
查询是否已关注
userId 关注 followId
*/
func (RelationServiceImpl) JudgeIsFollowById(userId int64, followId int64) bool {
	// 查redis是否已有记录
	redisFollowKey := util.Relation_Follow_Key + strconv.FormatInt(userId, 10)
	flag, err := redis.RedisDb.SIsMember(redis.Ctx, redisFollowKey, followId).Result()
	if err != nil {
		log.Println("redis query error!")
		return false
	}
	if flag {
		// 重现设置过期时间
		redis.RedisDb.Expire(redis.Ctx, redisFollowKey, util.Relation_Follow_TTL)
		return true
	}
	// 对于出现redis过期的情况：
	// redis中没记录，则在mysql中查
	flag = dao.JudgeIsFollowById(userId, followId)
	if flag {
		// 更新redis
		redis.RedisDb.SAdd(redis.Ctx, redisFollowKey, followId)
		redis.RedisDb.Expire(redis.Ctx, redisFollowKey, util.Relation_Follow_TTL)
	}
	return flag
}

/*
获取用户关注列表
*/
func (RelationServiceImpl) GetFollowList(userId int64) ([]dao.UserResp, error) {
	//#优化：关注列表由于要(大V/经常登录用户)返回多个用户的信息，为确保数据一致性，会带来频繁的缓存删除和增加操作，暂不做redis缓存
	usi := UserServiceImpl{}
	followList := make([]dao.UserResp, 0)
	followIds, err := dao.QueryFollowsIdByUserId(userId)
	if nil != err {
		return followList, err
	}
	for _, followId := range followIds {
		followInfo, err := usi.QueryUserRespById(followId)
		if nil != err {
			return followList, err
		}
		//关注列表一定已关注
		followInfo.IsFollow = true
		followList = append(followList, followInfo)
	}
	return followList, nil
}

/*
获取用户粉丝列表
*/
func (RelationServiceImpl) GetFollowerList(userId int64) ([]dao.UserResp, error) {
	//#优化：关注列表由于要(大V/经常登录用户)返回多个用户的信息，为确保数据一致性，会带来频繁的缓存删除和增加操作，暂不做redis缓存
	rsi := RelationServiceImpl{}
	usi := UserServiceImpl{}
	followerList := make([]dao.UserResp, 0)
	followerIds, err := dao.QueryFollowersIdByUserId(userId)
	if nil != err {
		return followerList, err
	}

	// 注：range获取数组项不能修改数组中结构体的值
	for _, followerId := range followerIds {
		followerInfo, err := usi.QueryUserRespById(followerId)
		isFollow := rsi.JudgeIsFollowById(userId, followerId)
		if nil != err {
			return followerList, err
		}
		if isFollow {
			followerInfo.IsFollow = true
		} else {
			followerInfo.IsFollow = false
		}
		followerList = append(followerList, followerInfo)
	}
	return followerList, nil
}

/*
获取用户好友列表
*/
func (RelationServiceImpl) GetFriendList(userId int64) ([]dao.FriendResp, error) {
	//#优化：关注列表由于要(大V/经常登录用户)返回多个用户的信息，为确保数据一致性，会带来频繁的缓存删除和增加操作，暂不做redis缓存
	//进入好友页面，先将useId的所有msgid的redis缓存给删掉，实现进入聊天框后重新访问一次全部聊天记录
	redisMessageIdKey := util.Message_MessageId_Key + strconv.FormatInt(userId, 10) + "_"
	redis.DelRedisCatchBatch(redisMessageIdKey)
	friendList := make([]dao.FriendResp, 0)
	rsi := RelationServiceImpl{}

	// 查出关注列表
	usi := UserServiceImpl{}
	followIds, err := dao.QueryFollowsIdByUserId(userId)
	if nil != err {
		return friendList, err
	}
	for _, followId := range followIds {
		tmpFriendInfo, err := usi.QueryUserRespById(followId)
		friendResp := dao.FriendResp{}
		// 判断是否回关，回关了即为好友
		isFollow := rsi.JudgeIsFollowById(followId, userId)
		if nil != err {
			return friendList, err
		}
		if isFollow {
			friendResp.Id = tmpFriendInfo.Id
			friendResp.Name = tmpFriendInfo.Name
			friendResp.FollowCount = tmpFriendInfo.FollowCount
			friendResp.FollowerCount = tmpFriendInfo.FollowerCount
			friendResp.IsFollow = true
			//friendResp.Avatar = "E:/personal/Go_workplace/tiny-tiktok/data/female.png"
			friendResp.Avatar = config.Ftp_image_path + "female.png"
			friendResp.FavoriteCount = tmpFriendInfo.FavoriteCount
			friendResp.WorkCount = tmpFriendInfo.WorkCount
			friendResp.TotalFavorited = tmpFriendInfo.TotalFavorited
			msi := MessageServiceImpl{}
			latestMsg, err := msi.GetLatestMessage(userId, tmpFriendInfo.Id)
			if err != nil {
				friendResp.Message = ""
				friendResp.MsgType = 0
			}
			friendResp.Message = latestMsg.Content
			friendResp.MsgType = latestMsg.MsgType

			friendList = append(friendList, friendResp)
		}
	}

	return friendList, nil
}

// 统计id用户粉丝数
func (RelationServiceImpl) CountFollowers(id int64) int64 {
	redisFollowerCntKey := util.Relation_FollowerCnt_Key + strconv.FormatInt(id, 10)
	// redis是否存在该键值对
	if cnt, err := redis.RedisDb.SCard(redis.Ctx, redisFollowerCntKey).Result(); cnt > 0 {
		if err != nil {
			log.Println("redis query error!")
			return -1
		}
		// 更新过期时间
		redis.RedisDb.Expire(redis.Ctx, redisFollowerCntKey, util.Relation_FollowerCnt_TTL)
		return cnt
	}
	cnt := dao.CountFollowers(id)
	redis.RedisDb.Set(redis.Ctx, redisFollowerCntKey, cnt, util.Relation_FollowerCnt_TTL)
	return cnt
}

// 统计id用户关注数
func (RelationServiceImpl) CountFollowings(id int64) int64 {
	redisFollowingCntKey := util.Relation_FollowingCnt_Key + strconv.FormatInt(id, 10)
	// redis是否存在该键值对
	if cnt, err := redis.RedisDb.SCard(redis.Ctx, redisFollowingCntKey).Result(); cnt > 0 {
		if err != nil {
			log.Println("redis query error!")
			return -1
		}
		// 更新过期时间
		redis.RedisDb.Expire(redis.Ctx, redisFollowingCntKey, util.Relation_FollowingCnt_TTL)
		return cnt
	}
	cnt := dao.CountFollowings(id)

	redis.RedisDb.Set(redis.Ctx, redisFollowingCntKey, cnt, util.Relation_FollowingCnt_TTL)
	return cnt
}

// 主动使cnt的redis缓存失效
func (RelationServiceImpl) ExpireFollowerCnt(id int64) {
	// 由于关注或取关操作导致cnt缓存不一致
	redisFollowerCntKey := util.Relation_FollowerCnt_Key + strconv.FormatInt(id, 10)
	redis.RedisDb.Del(redis.Ctx, redisFollowerCntKey)
}

// 主动使cnt的redis缓存失效
func (RelationServiceImpl) ExpireFollowingCnt(id int64) {
	// 由于关注或取关操作导致cnt缓存不一致
	redisFollowingCntKey := util.Relation_FollowingCnt_Key + strconv.FormatInt(id, 10)
	redis.RedisDb.Del(redis.Ctx, redisFollowingCntKey)
}
