package controller

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Scut-Corgis/tiny-tiktok/middleware/jwt"
)

func TestPublish(t *testing.T) {
	token := jwt.GenerateToken("qly", 1000)
	url := "http://127.0.0.1:8080/douyin/publish/action/"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	videoPath := "/Users/zaizai/Projects/GolandProjects/tiny-tiktok/data/bear.mp4"
	file, errFile1 := os.Open(videoPath)
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile("data", filepath.Base(videoPath))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return
	}
	_ = writer.WriteField("token", token)
	_ = writer.WriteField("title", "HeiBao")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")
	req.Header.Set("Content-Type", writer.FormDataContentType())
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

func TestFeed(t *testing.T) {
	url := "http://127.0.0.1:8080/douyin/feed/"
	method := "GET"
	SendRequest(method, url, nil)
}

func TestPublishList(t *testing.T) {
	token := jwt.GenerateToken("qly", 1000)
	url := "http://127.0.0.1:8080/douyin/publish/list/?user_id=1000&token=" + token
	method := "GET"
	SendRequest(method, url, strings.NewReader(token))
}
