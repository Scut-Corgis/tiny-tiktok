package dao

import (
	"log"
)

type Follow struct {
	UserId     int64 `gorm:"column:user_id"`
	FollowerId int64 `gorm:"column:follower_id"`
}

type UserResp struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

/*
增加follow关系
userId 关注 followId
*/
func InsertFollow(userId int64, followId int64) (bool, error) {
	follow := Follow{
		UserId:     followId,
		FollowerId: userId, // 登录的用户是粉丝
	}
	if err := Db.Table("follows").Create(&follow).Error; err != nil {
		log.Println(err.Error())
		return false, err
	}
	return true, nil
}

/*
删除follow关系
userId 关注 followId
*/
func DeleteFollow(userId int64, followId int64) (bool, error) {
	follow := Follow{}
	if err := Db.Table("follows").Where("user_id = ? and follower_id = ?", followId, userId).Delete(&follow).Error; err != nil {
		log.Println(err.Error())
		return false, err
	}
	return true, nil
}

/*
查询是否已关注
用户id1是否关注id2用户
*/
func JudgeIsFollowById(id1 int64, id2 int64) bool { // 判断用户id1是否关注id2用户
	var count int64
	Db.Model(&Follow{}).Where("user_id = ? and follower_id = ?", id2, id1).Count(&count)
	return count > 0
}

/*
查询用户关注id列表
*/
func QueryFollowsIdByUserId(userId int64) ([]int64, error) {
	followIds := make([]int64, 0)
	if err := Db.Table("follows").Select("user_id").Where("follower_id = ?", userId).Find(&followIds).Error; nil != err {
		return nil, err
	}
	return followIds, nil
}

/*
查询用户粉丝id列表
*/
func QueryFollowersIdByUserId(userId int64) ([]int64, error) {
	followerIds := make([]int64, 0)
	if err := Db.Table("follows").Select("follower_id").Where("user_id = ?", userId).Find(&followerIds).Error; nil != err {
		return nil, err
	}
	return followerIds, nil
}

func CountFollowers(id int64) int64 {
	var count int64
	Db.Model(&Follow{}).Where("user_id = ?", id).Count(&count)
	return count
}

func CountFollowings(id int64) int64 {
	var count int64
	Db.Model(&Follow{}).Where("follower_id = ?", id).Count(&count)
	return count
}

// 原生sql写法不启用
// /*
// 查询用户关注列表
// */
// func QueryFollowsByUserId(userId int64) ([]UserResp, error) {
// 	followList := make([]UserResp, 0)
// 	//子查询： 第一个select：先内连接查出 users.id = follows.user_id 的一个表，作为左基础表，再左外连接查询 users.id = follows.user_id，得到用户关注表
// 	//		  第二个select：先内连接查出 users.id = follows.user_id 的一个表，作为左基础表，再左外连接查询 users.id = follows.follower_id，得到用户粉丝表,只能体现粉丝个数，没有粉丝具体信息
// 	//外查询： 两个表联合在一起，按照follower_id、id、name进行group by分组。查出 follower.id = userId 的结果，即当前用户作为粉丝，他关注的人的一个列表。
// 	// 		  tag字段用来统计有多少个粉丝和关注。
// 	if err := Db.Raw("select id, name, "+
// 		"\ncount(if(tag = 'follower', 1, null)) follower_count, "+
// 		"\ncount(if(tag = 'follow', 1, null)) follow_count, 'true' is_follow "+
// 		"\nfrom (select f1.follower_id fid, u.id, name, 'follower' tag "+
// 		"\nfrom follows f1 join users u on f1.user_id = u.id left join follows f2 on u.id = f2.user_id union all "+
// 		"\nselect f1.follower_id fid, u.id, name, 'follow' tag "+
// 		"\nfrom follows f1 join users u on f1.user_id = u.id "+
// 		"\nleft join follows f2 on u.id = f2.follower_id) T where fid = ? group by fid, id, name", userId).Scan(&followList).Error; nil != err {
// 		return nil, err
// 	}

// 	return followList, nil
// }

// /*
// 查询用户粉丝列表
// */
// func QueryFollowersByUserId(userId int64) ([]UserResp, error) {
// 	followerList := make([]UserResp, 1)
// 	if err := Db.Raw("select id, name, "+
// 		"\ncount(if(tag = 'follower', 1, null)) follower_count, "+
// 		"\ncount(if(tag = 'follow', 1, null)) follow_count, 'true' is_follow "+
// 		"\nfrom (select f1.user_id uid, u.id, name, 'follow' tag from follows f1 "+
// 		"\njoin users u on f1.follower_id = u.id left join follows f2 on u.id = f2.follower_id "+
// 		"\nunion all select f1.user_id uid, u.id, name, 'follower' tag from follows f1 "+
// 		"\njoin users u on f1.follower_id = u.id left join follows f2 on u.id = f2.user_id) T "+
// 		"\nwhere uid = ? group by uid, id, name", userId).Scan(&followerList).Error; nil != err {
// 		return nil, err
// 	}

// 	return followerList, nil
// }

// /*
// 查询用户好友id列表
// */
// func QueryFriendsIdByUserId(userId int64) ([]int64, error) {
// 	friendIds := make([]int64, 0)
// 	// 不用join,实现高效的互相关注查询：https://blog.csdn.net/a_void/article/details/103273954
// 	// 按照字典顺序做一次排序，那么排序后的结果都是(A, B), (A, B)。
// 	// 思路：把特征相同的数据分到一组，计算组里面的数据条数，为1则是单向关注，为2则是双向关注。这样，利用窗口函数，不用join就能得到答案
// 	if err := Db.Raw("select b.follower_id from "+
// 		"\n(select a.follower_id, a.user_id, if(sum(1) over (partition by feature) > 1, 1, 0) as is_friend from "+
// 		"\n(select follower_id, user_id, "+
// 		"\nif(follower_id > user_id, concat(user_id, follower_id), concat(follower_id, user_id)) as feature "+
// 		"\nfrom follows)a "+
// 		"\n)b"+
// 		"\nwhere b.is_friend > 0 and b.user_id = ?", userId).Scan(&friendIds).Error; nil != err {
// 		return nil, err
// 	}
// 	return friendIds, nil
// }
