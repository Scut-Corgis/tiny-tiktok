package service

import "github.com/Scut-Corgis/tiny-tiktok/dao"

type FollowServiceImpl struct{}

// JudgeIsFollowById 判断用户id1是否关注id2用户
func (FollowServiceImpl) JudgeIsFollowById(id1 int64, id2 int64) bool {
	return dao.JudgeIsFollowById(id1, id2)
}

func (FollowServiceImpl) CountFollowers(id int64) int64 {
	return dao.CountFollowers(id)
}

func (FollowServiceImpl) CountFollowings(id int64) int64 {
	return dao.CountFollowings(id)
}
