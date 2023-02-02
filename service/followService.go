package service

type FollowService interface {
	// JudgeIsFollowById 判断用户id1是否关注了用户id2
	JudgeIsFollowById(id1 int64, id2 int64) bool

	CountFollowers(id int64) int64
	CountFollowings(id int64) int64
}
