package controller

import (
	"strings"
	"testing"

	"github.com/Scut-Corgis/tiny-tiktok/middleware/jwt"
)

func TestFavoriteAction(t *testing.T) {
	//点赞1
	token := "token=" + jwt.GenerateToken("weipengyan2", 1107)
	url := "http://47.108.112.214:8080/douyin/favorite/action/?video_id=1069&action_type=1"
	method := "POST"
	SendRequest(method, url, strings.NewReader(token))
	//点赞2
	token = "token=" + jwt.GenerateToken("weipengyan2", 1107)
	url = "http://47.108.112.214:8080/douyin/favorite/action/?video_id=1070&action_type=1"
	method = "POST"
	SendRequest(method, url, strings.NewReader(token))
	//取消点赞
	token = "token=" + jwt.GenerateToken("weipengyan2", 1107)
	url = "http://47.108.112.214:8080/douyin/favorite/action/?video_id=1069&action_type=2"
	method = "POST"
	SendRequest(method, url, strings.NewReader(token))
}

func TestFavoriteList(t *testing.T) {
	token := jwt.GenerateToken("weipengyan2", 1107)
	url := "http://47.108.112.214:8080/douyin/favorite/list/?user_id=1107&token=" + token
	method := "GET"
	SendRequest(method, url, nil)
}
