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

// QueryUserByName 根据name查询User
func QueryUserByName(name string) (User, error) {
	user := User{}
	if err := Db.Where("name = ?", name).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// QueryUserById 根据id查询User
func QueryUserById(id int64) (User, error) {
	user := User{}
	if err := Db.Where("id = ?", id).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// QueryUserRespById 根据id查询UserResp
func QueryUserRespById(id int64) (UserResp, error) {
	userInfo := UserResp{}
	user, err := QueryUserById(id)
	if err != nil {
		log.Println(err.Error())
		return userInfo, err
	} else {
		userInfo.FollowerCount = CountFollowers(id) // 统计粉丝数量
		userInfo.FollowCount = CountFollowings(id)  // 统计关注博主的数量
		userInfo.Id = user.Id
		userInfo.Name = user.Name
		return userInfo, err
	}
}

// InsertUser 将User插入user表内
func InsertUser(user *User) bool {
	if err := Db.Create(&user).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
