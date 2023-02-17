package controller

import (
	"github.com/Scut-Corgis/tiny-tiktok/middleware/jwt"
	"strings"
	"testing"
)

func TestRelationAction(t *testing.T) {
	token := "token=" + jwt.GenerateToken("qly", 1000)
	url := "http://127.0.0.1:8080/douyin/relation/action/?to_user_id=&action_type="
	method := "POST"
	SendRequest(method, url, strings.NewReader(token))
}

func TestFollowList(t *testing.T) {
	token := jwt.GenerateToken("qly", 1000)
	url := "http://127.0.0.1:8080/douyin/relation/follow/list/?user_id=1000&token=" + token
	method := "GET"
	SendRequest(method, url, nil)
}

func TestFollowerList(t *testing.T) {
	token := jwt.GenerateToken("qly", 1000)
	url := "http://127.0.0.1:8080/douyin/relation/follower/list/?user_id=1000&token=" + token
	method := "GET"
	SendRequest(method, url, nil)
}

func TestFriendList(t *testing.T) {
	token := jwt.GenerateToken("qly", 1000)
	url := "http://127.0.0.1:8080/douyin/relation/friend/list/?user_id=1000&token=" + token
	method := "GET"
	SendRequest(method, url, nil)
}
