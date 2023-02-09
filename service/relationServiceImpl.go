package service

import (
	"log"
	"strconv"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
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
	err := dao.InsertFollow(userId, followId)
	if err != nil {
		return false, err
	}

	// 将查询到的关注关系注入Redis
	redisFollowKey := util.Relation_Follow_Key + strconv.FormatInt(userId, 10)
	redis.RedisDb.SAdd(redis.Ctx, redisFollowKey, followId)
	// 更新过期时间
	redis.RedisDb.Expire(redis.Ctx, redisFollowKey, util.Relation_Follow_TTL)
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
	err := dao.DeleteFollow(userId, followId)
	if err != nil {
		return false, err
	}

	// 删除Redis中 redisFollowKey set集合中的followId元素
	redisFollowKey := util.Relation_Follow_Key + strconv.FormatInt(userId, 10)
	redis.RedisDb.SRem(redis.Ctx, redisFollowKey, followId)
	// 更新过期时间
	redis.RedisDb.Expire(redis.Ctx, redisFollowKey, util.Relation_Follow_TTL)

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
	return dao.JudgeIsFollowById(userId, followId)
}

/*
获取用户关注列表
*/
func (RelationServiceImpl) GetFollowList(userId int64) ([]dao.UserResp, error) {
	followList := make([]dao.UserResp, 0)

	followIds, err := dao.QueryFollowsIdByUserId(userId)
	if nil != err {
		return followList, err
	}
	for _, followId := range followIds {
		followInfo, err := dao.QueryUserRespById(followId)
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
	followerList := make([]dao.UserResp, 0)
	followerIds, err := dao.QueryFollowersIdByUserId(userId)
	if nil != err {
		return followerList, err
	}
	// 注：range获取数组项不能修改数组中结构体的值
	for _, followerId := range followerIds {
		followerInfo, err := dao.QueryUserRespById(followerId)
		isFollow := dao.JudgeIsFollowById(userId, followerId)
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
func (RelationServiceImpl) GetFriendList(userId int64) ([]dao.UserResp, error) {
	friendList := make([]dao.UserResp, 0)
	// 查出关注列表
	followIds, err := dao.QueryFollowsIdByUserId(userId)
	if nil != err {
		return friendList, err
	}
	for _, followId := range followIds {
		tmpFriendInfo, err := dao.QueryUserRespById(followId)
		// 判断是否回关，回关了即为好友
		isFollow := dao.JudgeIsFollowById(followId, userId)
		if nil != err {
			return friendList, err
		}
		if isFollow {
			tmpFriendInfo.IsFollow = true
			friendList = append(friendList, tmpFriendInfo)
		}
	}
	return friendList, nil
}

// 统计id用户粉丝数
func (RelationServiceImpl) CountFollowers(id int64) int64 {
	return dao.CountFollowers(id)
}

// 统计id用户关注数
func (RelationServiceImpl) CountFollowings(id int64) int64 {
	return dao.CountFollowings(id)
}
