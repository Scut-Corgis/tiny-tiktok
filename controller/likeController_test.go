package controller

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
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

	token := jwt.GenerateToken("Corgis")
	url := "http://127.0.0.1:8080/douyin/favorite/action/?token=" + token + "video_id=1007&action_type=1"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("token", token)
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	//payload := bytes.NewReader([]byte(token))

	// post := "{\"token\":\"" + token + "\"}"
	// var payload = []byte(post)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")

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
