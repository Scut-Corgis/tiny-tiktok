package controller

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/Scut-Corgis/tiny-tiktok/middleware/jwt"
)

func TestFavoriteAction(t *testing.T) {
	//注册用户Corgis
	url1 := "http://127.0.0.1:8080/douyin/user/register/?username=Corgis&password=123456"
	method1 := "POST"

	client1 := &http.Client{}
	req1, err := http.NewRequest(method1, url1, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	res1, err := client1.Do(req1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res1.Body.Close()

	body1, err := io.ReadAll(res1.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body1))

	//登录用户
	url2 := "http://127.0.0.1:8080/douyin/user/login/?username=Corgis&password=123456"
	method2 := "POST"

	client2 := &http.Client{}
	req2, err := http.NewRequest(method2, url2, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req2.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")

	res2, err := client2.Do(req2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res2.Body.Close()

	body2, err := ioutil.ReadAll(res2.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body2))

	token := "token=" + jwt.GenerateToken("Corgis")
	url := "http://127.0.0.1:8080/douyin/favorite/action/?video_id=1007&action_type=1"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(token))
	log.Println(&req)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
