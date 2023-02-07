package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

type RelationServiceImpl struct{}

/*
关注用户
userId 关注 followId
*/
func (RelationServiceImpl) Follow(userId int64, followId int64) (bool, error) {
	return dao.InsertFollow(userId, followId)
}

/*
取关用户
userId 关注 followId
*/
func (RelationServiceImpl) UnFollow(userId int64, followId int64) (bool, error) {
	return dao.DeleteFollow(userId, followId)
}

/*
查询是否已关注
userId 关注 followId
*/
func (RelationServiceImpl) IsFollowed(userId int64, followId int64) (bool, error) {
	isFollow, err := dao.QueryIsFollowByUserId(userId, followId)
	// isFollow为0，表示未关注
	if nil == err && isFollow == 0 {
		return false, err
	}
	return true, err
}

/*
获取用户关注列表
*/
func (RelationServiceImpl) GetFollowList(userId int64) ([]dao.UserResp, error) {
	return dao.QueryFollowsByUserId(userId)
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
		followerInfo, err1 := dao.QueryUserRespById(followerId)
		isFollow, err2 := dao.QueryIsFollowByUserId(userId, followerId)
		if nil != err1 || nil != err2 {
			return followerList, err
		}
		if isFollow == 0 {
			followerInfo.IsFollow = false
		} else {
			followerInfo.IsFollow = true
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
	// 查出好友的id
	friendIds, err := dao.QueryFriendsIdByUserId(userId)
	if nil != err {
		return friendList, err
	}
	// 查每个好友的信息
	for _, friendId := range friendIds {
		friendInfo, err1 := dao.QueryUserRespById(friendId)
		isFollow, err2 := dao.QueryIsFollowByUserId(userId, friendId)
		if nil != err1 || nil != err2 {
			return friendList, err
		}
		if isFollow == 0 {
			friendInfo.IsFollow = false
		} else {
			friendInfo.IsFollow = true
		}
		friendList = append(friendList, friendInfo)
	}
	return friendList, nil
}
