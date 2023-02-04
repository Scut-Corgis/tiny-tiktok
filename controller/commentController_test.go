package controller

import (
	"github.com/Scut-Corgis/tiny-tiktok/middleware/jwt"
	"strings"
	"testing"
)

func TestCommentAction(t *testing.T) {
	token := "token=" + jwt.GenerateToken("qly")

	// 评论操作——添加评论
	url1 := "http://127.0.0.1:8080/douyin/comment/action/?video_id=1002&action_type=1&comment_text=fuck"
	method1 := "POST"
	SendRequest(method1, url1, strings.NewReader(token))

	// 评论操作——删除评论
	url2 := "http://127.0.0.1:8080/douyin/comment/action/?video_id=1002&action_type=2&comment_id=1000"
	method2 := "POST"
	SendRequest(method2, url2, strings.NewReader(token))
}

func TestCommentList(t *testing.T) {
	token := jwt.GenerateToken("qly")
	url1 := "http://127.0.0.1:8080/douyin/comment/list/?video_id=1002&token=" + token
	method1 := "GET"
	SendRequest(method1, url1, strings.NewReader(token))
}
