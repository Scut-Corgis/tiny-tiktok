package service

import "github.com/Scut-Corgis/tiny-tiktok/dao"

type UserService interface {
	// QueryUserByName 根据name获取User对象
	QueryUserByName(name string) dao.User

	// QueryUserById 根据id获取User对象
	QueryUserById(id int64) dao.User

	// QueryUserRespById 根据id获取UserResp对象
	QueryUserRespById(id int64) (dao.UserResp, error)

	// Register 用户注册
	Register(username string, password string) (dao.User, int32, string)

	// Login 用户登录
	Login(username string, password string) (int32, string)
}
