package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

/*
关注用户
*/
func Follow(userId int64, followId int64) bool {
	return dao.InsertFollow(userId, followId)
}

/*
取关用户
*/
func UnFollow(userId int64, followId int64) bool {
	return dao.DeleteFollow(userId, followId)
}

/*
获取用户关注列表
#优化：这里遍历了整个user
*/
func FollowList(userId int64) ([]dao.UserTable, error) {
	return dao.QueryFollowsByUserId(userId)

}