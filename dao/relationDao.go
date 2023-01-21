package dao

import "log"

type FollowTable struct {
	UserId     int64 `gorm:"column:user_id"`
	FollowerId int64 `gorm:"column:follower_id"`
}

func InsertFollow(userId int64, followId int64) bool {
	follow := FollowTable{
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
	follow := FollowTable{}
	if err := Db.Table("follows").Where("user_id = ? and follower_id = ?", userId, followId).Delete(&follow).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
