package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func QueryUserByName(name string) dao.User {
	user, err := dao.QueryUserByName(name)
	if err != nil {
		log.Println("error:", err.Error())
		log.Println("User not found!")
		return user
	}
	log.Println("Query user successfully!")
	return user
}

func InsertUser(user *dao.User) bool {
	flag := dao.InsertUser(user)
	if flag == false {
		log.Println("Insert user failed!")
		return flag
	}
	log.Println("Insert user successfully!")
	return flag
}

// HashEncode hash加密密码
func HashEncode(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePasswords(password1 string, password2 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
	if err != nil {
		return false
	}
	return true
}
