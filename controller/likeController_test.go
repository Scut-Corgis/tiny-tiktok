package controller

import (
	"strings"
	"testing"

	"github.com/Scut-Corgis/tiny-tiktok/middleware/jwt"
)

func TestFavoriteAction(t *testing.T) {
	token := "token=" + jwt.GenerateToken("Corgis", 1000)
	url := "http://127.0.0.1:8080/douyin/favorite/action/?video_id=1000&action_type=1"
	method := "POST"
	SendRequest(method, url, strings.NewReader(token))
}

func TestFavoriteList(t *testing.T) {
	token := jwt.GenerateToken("qly", 1000)
	url := "http://127.0.0.1:8080/douyin/favorite/list/?user_id=1000&token=" + token
	method := "GET"
	SendRequest(method, url, nil)
}
