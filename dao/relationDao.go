package dao

import (
	"log"
)

type Follow struct {
	UserId     int64 `gorm:"column:user_id"`
	FollowerId int64 `gorm:"column:follower_id"`
}

type UserResp struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	FollowCount    int64  `json:"follow_count"`
	FollowerCount  int64  `json:"follower_count"`
	IsFollow       bool   `json:"is_follow"`
	TotalFavorited int64  `json:"total_favorited"`
	WorkCount      int64  `json:"work_count"`
	FavoriteCount  int64  `json:"favorite_count"`
}

type FriendResp struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	FollowCount    int64  `json:"follow_count"`
	FollowerCount  int64  `json:"follower_count"`
	IsFollow       bool   `json:"is_follow"`
	Avatar         string `json:"avatar"`
	TotalFavorited int64  `json:"total_favorited"`
	WorkCount      int64  `json:"work_count"`
	FavoriteCount  int64  `json:"favorite_count"`
	Message        string `json:"message"`
	MsgType        int64  `json:"msgType"`
}

// InsertFollow 增加follow关系 userId 关注 followId
func InsertFollow(userId int64, followId int64) error {
	follow := Follow{
		UserId:     followId,
		FollowerId: userId, // 登录的用户是粉丝
	}
	if err := Db.Table("follows").Create(&follow).Error; err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// DeleteFollow 删除follow关系 userId 关注 followId
func DeleteFollow(userId int64, followId int64) error {
	follow := Follow{}
	if err := Db.Table("follows").Where("user_id = ? and follower_id = ?", followId, userId).Delete(&follow).Error; err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// JudgeIsFollowById 查询是否已关注 用户id1是否关注id2用户
func JudgeIsFollowById(id1 int64, id2 int64) bool { // 判断用户id1是否关注id2用户
	var count int64
	Db.Model(&Follow{}).Where("user_id = ? and follower_id = ?", id2, id1).Count(&count)
	return count > 0
}

// QueryFollowsIdByUserId 查询用户关注id列表
func QueryFollowsIdByUserId(userId int64) ([]int64, error) {
	followIds := make([]int64, 0)
	if err := Db.Table("follows").Select("user_id").Where("follower_id = ?", userId).Find(&followIds).Error; nil != err {
		return nil, err
	}
	return followIds, nil
}

// QueryFollowersIdByUserId 查询用户粉丝id列表
func QueryFollowersIdByUserId(userId int64) ([]int64, error) {
	followerIds := make([]int64, 0)
	if err := Db.Table("follows").Select("follower_id").Where("user_id = ?", userId).Find(&followerIds).Error; nil != err {
		return nil, err
	}
	return followerIds, nil
}

// CountFollowers 统计粉丝数
func CountFollowers(id int64) int64 {
	var count int64
	Db.Model(&Follow{}).Where("user_id = ?", id).Count(&count)
	return count
}

// CountFollowings 统计关注数
func CountFollowings(id int64) int64 {
	var count int64
	Db.Model(&Follow{}).Where("follower_id = ?", id).Count(&count)
	return count
}
