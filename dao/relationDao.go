package dao

import (
	"log"
)

type Follow struct {
	UserId     int64 `gorm:"column:user_id"`
	FollowerId int64 `gorm:"column:follower_id"`
}

type UserTable struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

/*
增加一条follows表数据
*/
func InsertFollow(userId int64, followId int64) bool {
	follow := Follow{
		UserId:     userId,
		FollowerId: followId,
	}
	if err := Db.Table("follows").Create(&follow).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

/*
删除一条follows表数据
*/
func DeleteFollow(userId int64, followId int64) bool {
	follow := Follow{}
	if err := Db.Table("follows").Where("user_id = ? and follower_id = ?", userId, followId).Delete(&follow).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

/*
查询用户关注列表
*/
func QueryFollowsByUserId(userId int64) ([]UserTable, error) {
	followList := make([]UserTable, 1)

	//子查询： 第一个select：先内连接查出 users.id = follows.user_id 的一个表，作为左基础表，再左外连接查询 users.id = follows.user_id，得到用户关注表
	//		  第二个select：先内连接查出 users.id = follows.user_id 的一个表，作为左基础表，再左外连接查询 users.id = follows.follower_id，得到用户粉丝表,只能体现粉丝个数，没有粉丝具体信息
	//外查询： 两个表联合在一起，按照follower_id、id、name进行group by分组。查出 follower.id = userId 的结果，即当前用户作为粉丝，他关注的人的一个列表。
	// 		  tag字段用来统计有多少个粉丝和关注。
	if err := Db.Raw("select id, name, "+
		"\ncount(if(tag = 'follower', 1, null)) follower_count, "+
		"\ncount(if(tag = 'follow', 1, null)) follow_count, 'true' is_follow "+
		"\nfrom (select f1.follower_id fid, u.id, name, 'follower' tag "+
		"\nfrom follows f1 join users u on f1.user_id = u.id left join follows f2 on u.id = f2.user_id union all "+
		"\nselect f1.follower_id fid, u.id, name, 'follow' tag "+
		"\nfrom follows f1 join users u on f1.user_id = u.id "+
		"\nleft join follows f2 on u.id = f2.follower_id) T where fid = ? group by fid, id, name", userId).Scan(&followList).Error; nil != err {
		return nil, err
	}

	return followList, nil
}

/*
查询用户粉丝列表
*/
func QueryFollowersByUserId(userId int64) ([]UserTable, error) {
	followerList := make([]UserTable, 1)
	if err := Db.Raw("select id, name, "+
		"\ncount(if(tag = 'follower', 1, null)) follower_count, "+
		"\ncount(if(tag = 'follow', 1, null)) follow_count, "+
		"\nif((select count(*) from follows where follower_id = ? and user_id = id) > 0, 'true', 'false') is_follow "+
		"\nfrom (select f1.user_id uid, u.id, name, 'follower' tag from follows f1 "+
		"\njoin users u on f1.follower_id = u.id left join follows f2 on u.id = f2.follower_id "+
		"\nunion all select f1.user_id uid, u.id, name, 'follow' tag from follows f1 "+
		"\njoin users u on f1.follower_id = u.id left join follows f2 on u.id = f2.user_id) T "+
		"\nwhere uid = ? group by uid, id, name", userId, userId).Scan(&followerList).Error; nil != err {
		return nil, err
	}

	return followerList, nil
}

/*
通过UserId查询当前用户的 UserTable 信息
*/
func QueryUserByUserId(userId int64) (UserTable, error) {
	userTable := UserTable{}
	if err := Db.Where("id = ?", userId).First(&userTable).Error; err != nil {
		return userTable, err
	}
	return userTable, nil
}
