/*
userController测试文件
go test -v userController_test.go common_test.go
*/
package controller

import (
	"github.com/Scut-Corgis/tiny-tiktok/middleware/jwt"
	"testing"
)

func TestRegister(t *testing.T) {
	// 用户注册——成功
	url1 := "http://127.0.0.1:8080/douyin/user/register/?username=qly&password=123"
	method1 := "POST"
	SendRequest(method1, url1, nil)
	// 用户注册——用户名重复
	url2 := "http://127.0.0.1:8080/douyin/user/register/?username=qly&password=123"
	method2 := "POST"
	SendRequest(method2, url2, nil)
}

func TestLogin(t *testing.T) {
	// 用户登录——成功
	url1 := "http://127.0.0.1:8080/douyin/user/login/?username=qly&password=123"
	method1 := "POST"
	SendRequest(method1, url1, nil)
	// 用户登录——密码错误
	url2 := "http://127.0.0.1:8080/douyin/user/login/?username=qly&password=122"
	method2 := "POST"
	SendRequest(method2, url2, nil)
}

func TestUserInfo(t *testing.T) {
	token := jwt.GenerateToken("wpy51", 1089)

	// 用户信息——用户不存在
	url1 := "http://127.0.0.1:8080/douyin/user/?user_id=9999&token=" + token
	method1 := "GET"
	SendRequest(method1, url1, nil)

	// 用户信息——鉴权失败
	url2 := "http://127.0.0.1:8080/douyin/user/?user_id=1000&token=" + "fuck"
	method2 := "GET"
	SendRequest(method2, url2, nil)

	// 用户信息——成功
	url3 := "http://127.0.0.1:8080/douyin/user/?user_id=1000&token=" + token
	method3 := "GET"
	SendRequest(method3, url3, nil)
}
