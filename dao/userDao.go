package dao

import (
	"log"
)

// User 定义用户表结构
type User struct {
	Id       int64
	Name     string
	Password string
}

// QueryUserByName 根据用户名查询用户
func QueryUserByName(name string) (User, error) {
	user := User{}
	if err := Db.Where("name = ?", name).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// QueryUserById 根据用户id查询用户
func QueryUserById(id int64) (User, error) {
	user := User{}
	if err := Db.Where("id = ?", id).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

func JudgeIsFollow(id int64, name string) bool { // 判断name用户是否关注id用户
	user1, err := QueryUserById(id)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	user2, _ := QueryUserByName(name)
	var count int64
	Db.Model(&FollowTable{}).Where("user_id = ? and follower_id = ?", user1.Id, user2.Id).Count(&count)
	return count > 0
}

func QueryUserTableById(id int64) (UserTable, error) {
	userInfo := UserTable{}
	user, err := QueryUserById(id)
	if err != nil {
		log.Println(err.Error())
		return userInfo, err
	} else {
		Db.Model(&FollowTable{}).Where("user_id = ?", id).Count(&userInfo.FollowerCount)   // 统计粉丝数量
		Db.Model(&FollowTable{}).Where("follower_id = ?", id).Count(&userInfo.FollowCount) // 统计关注博主的数量
		userInfo.Id = user.Id
		userInfo.Name = user.Name
		return userInfo, err
	}
}

// InsertUser 将user插入表内
func InsertUser(user *User) bool {
	if err := Db.Create(&user).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
