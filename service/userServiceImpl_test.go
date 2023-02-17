package service

import (
	"fmt"
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"testing"
)

func UserServiceImplInit() {
	dao.Init()
	redis.InitRedis()
	redis.InitCuckooFilter()
}

func TestUserServiceImpl_QueryUserByName(t *testing.T) {
	UserServiceImplInit()
	usi := UserServiceImpl{}
	user := usi.QueryUserByName("qly")
	fmt.Println(user)
}

func TestUserServiceImpl_QueryUserById(t *testing.T) {
	UserServiceImplInit()
	usi := UserServiceImpl{}
	user := usi.QueryUserById(1000)
	fmt.Println(user)
}

func TestUserServiceImpl_QueryUserRespById(t *testing.T) {
	UserServiceImplInit()
	usi := UserServiceImpl{}
	userResp, err := usi.QueryUserRespById(1000)
	fmt.Println(userResp)
	fmt.Println(err)
}

func TestUserServiceImpl_Register(t *testing.T) {
	UserServiceImplInit()
	usi := UserServiceImpl{}
	id, code, message := usi.Register("qly", "1000")
	fmt.Println(id, code, message)
}

func TestUserServiceImpl_Login(t *testing.T) {
	UserServiceImplInit()
	usi := UserServiceImpl{}
	code, message := usi.Login("qly", "1000")
	fmt.Println(code, message)
}
