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

// InsertUser 将user插入表内
func InsertUser(user *User) bool {
	if err := Db.Create(&user).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
