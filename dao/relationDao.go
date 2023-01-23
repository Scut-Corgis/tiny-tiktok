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

func DeleteFollow(userId int64, followId int64) bool {
	follow := Follow{}
	if err := Db.Table("follows").Where("user_id = ? and follower_id = ?", userId, followId).Delete(&follow).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func QueryFollowsByUserId(userId int64) ([]UserTable, error) {
	followList := make([]UserTable, 1)
	if err := Db.Raw("select id, name, "+
		"\ncount(if(tag = 'follower', 1, null)) follower_count, "+
		"\ncount(if(tag = 'follow', 1, null)) follow_count, 'true' is_follow "+
		"\nfrom (select f1.follower_id fid, u.id, name, 'follower' tag "+
		"\nfrom follows f1 join users u on f1.user_id = u.id left join follows f2 on u.id = f2.user_id union all "+
		"\nselect f1.follower_id fid, u.id, name, 'follow' tag from follows f1 join users u on f1.user_id = u.id "+
		"\nleft join follows f2 on u.id = f2.follower_id) T where fid = ? group by fid, id, name", userId).Scan(&followList).Error; nil != err {
		return nil, err
	}

	return followList, nil
}
