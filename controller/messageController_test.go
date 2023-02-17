package controller

import (
	"github.com/Scut-Corgis/tiny-tiktok/middleware/jwt"
	"strings"
	"testing"
)

func TestMessageAction(t *testing.T) {
	token := "token=" + jwt.GenerateToken("qly", 1000)
	url := "http://127.0.0.1:8080/douyin/message/action/?to_user_id=1001&action_type=1&content=test"
	method := "POST"
	SendRequest(method, url, strings.NewReader(token))
}

func TestChatRecord(t *testing.T) {
	token := jwt.GenerateToken("qly", 1000)
	url := "http://127.0.0.1:8080/douyin/message/chat/?to_user_id=1001&token=" + token
	method := "GET"
	SendRequest(method, url, nil)
}
