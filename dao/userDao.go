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

// InsertUser 将User插入users表内
func InsertUser(user *User) bool {
	if err := Db.Create(&user).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

// QueryAllNames 查询所有用户名
func QueryAllNames() []string {
	names := make([]string, 0)
	if err := Db.Table("users").Select("name").Find(&names).Error; err != nil {
		log.Println(err.Error())
		return names
	}
	return names
}
