package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

type RelationService interface {
	// Follow 关注用户followId
	Follow(userId int64, followId int64) (bool, error)

	// UnFollow 取关用户followId
	UnFollow(userId int64, followId int64) (bool, error)

	// IsFollowed 查询是否已关注followId
	IsFollowed(userId int64, followId int64) (bool, error)

	// GetFollowList 获取用户关注列表
	GetFollowList(userId int64) ([]dao.UserResp, error)

	// GetFollowerList 获取用户粉丝列表
	GetFollowerList(userId int64) ([]dao.UserResp, error)

	// GetFriendList 获取用户好友列表
	GetFriendList(userId int64) ([]dao.UserResp, error)
}
