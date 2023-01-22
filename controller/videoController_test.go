package controller

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/Scut-Corgis/tiny-tiktok/middleware/jwt"
)

// 该测试函数没用，gin端口为8080，而http包只能往80发
func TestPublish(t *testing.T) {
	//注册用户Corgis
	url := "http://127.0.0.1:8080/douyin/user/register/?username=Corgis&password=123456"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	token := jwt.GenerateToken("Corgis")
	return
	//投稿测试
	url = "http://127.0.0.1:8080//douyin/publish/action/"
	method = "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	videoPath := "/home/hjg/go/src/tiny-tiktok/data/heibao.mp4"
	file, errFile1 := os.Open(videoPath)
	if errFile1 != nil {
		log.Println("测试无法打开视频文件")
	}
	defer file.Close()
	part1,
		errFile1 := writer.CreateFormFile("data", filepath.Base(videoPath))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return
	}
	_ = writer.WriteField("token", token)
	_ = writer.WriteField("title", "HeiBao")
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client = &http.Client{}
	req, err = http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}
