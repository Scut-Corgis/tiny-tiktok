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

	// redis缓存操作
	redisFollowerCntKey := util.Relation_FollowerCnt_Key + strconv.FormatInt(followId, 10)
	redisFollowingCntKey := util.Relation_FollowingCnt_Key + strconv.FormatInt(userId, 10)
	// 不加锁版：redis读写本身就很快，直接失效即可保证一致性。
	// 保证数据一致性：主动使count缓存失效
	redis.DelRedisCatchBatch(redisFollowerCntKey, redisFollowingCntKey)

	// // 加锁版：加锁是否反而会影响性能
	// rand.Seed(time.Now().UnixNano())
	// value := strconv.Itoa(rand.Int())
	// lockFollower := redis.Lock(redisFollowerCntKey, value)
	// lockFollowing := redis.Lock(redisFollowingCntKey, value)
	// if lockFollower && lockFollowing {
	// 	_, err1 := redis.RedisDb.Incr(redis.Ctx, redisFollowerCntKey).Result()
	// 	_, err2 := redis.RedisDb.Incr(redis.Ctx, redisFollowingCntKey).Result()
	// 	if err1 != nil || err2 != nil {
	// 		// 保证数据一致性：主动使count缓存失效
	// 		redis.DelRedisCatchBatch(redisFollowerCntKey, redisFollowingCntKey)
	// 	}
	// } else {
	// 	// 保证数据一致性：主动使count缓存失效
	// 	redis.DelRedisCatchBatch(redisFollowerCntKey, redisFollowingCntKey)
	// }
	// if lockFollower {
	// 	unlock := redis.Unlock(redisFollowerCntKey)
	// 	if !unlock {
	// 		return false, nil
	// 	}
	// }
	// if lockFollowing {
	// 	unlock := redis.Unlock(redisFollowerCntKey)
	// 	if !unlock {
	// 		return false, nil
	// 	}
	// }
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

	// 保证数据一致性：主动使count缓存失效
	redisFollowerCntKey := util.Relation_FollowerCnt_Key + strconv.FormatInt(followId, 10)
	redisFollowingCntKey := util.Relation_FollowingCnt_Key + strconv.FormatInt(userId, 10)
	redis.DelRedisCatchBatch(redisFollowerCntKey, redisFollowingCntKey)
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

			friendResp.Avatar = "http://" + config.Url_addr + config.Url_Image_prefix + "male.png"
			friendResp.FavoriteCount = tmpFriendInfo.FavoriteCount
			friendResp.WorkCount = tmpFriendInfo.WorkCount
			friendResp.TotalFavorited = tmpFriendInfo.TotalFavorited
			msi := MessageServiceImpl{}
			latestMsg, err := msi.GetLatestMessage(userId, tmpFriendInfo.Id)
			// 如果没有最新消息,消息列表填空,以防出bug
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
