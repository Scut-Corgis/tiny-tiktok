package dao

import (
	"bytes"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func RebuildTable() bool {
	cmd := exec.Command("sh", "/Users/zaizai/Projects/GolandProjects/tiny-tiktok/config/rebuildTable.sh")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err, stderr.String())
		return false
	}
	return true
}

func FakeUsers(num int) {
	gofakeit.Seed(time.Now().Unix())
	for i := 0; i < num; i++ {
		user := User{}
		user.Name = gofakeit.Username()
		user.Password = gofakeit.Password(false, false, true, false, false, 8)
		_, err := InsertUser(user)
		if err != nil {
			return
		}
	}
}

func FakeFollows(num int) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < num; i++ {
		var count int64
		Db.Model(&User{}).Count(&count)
		var a, b int64
		a = rand.Int63n(count)
		b = rand.Int63n(count)
		for a == b {
			b = rand.Int63n(count)
		}
		err := InsertFollow(a+1000, b+1000)
		if err != nil {
			return
		}
	}
}

func FakeVideos(num int) {
	gofakeit.Seed(time.Now().Unix())
	rand.Seed(time.Now().Unix())
	for i := 0; i < num; i++ {
		var count int64
		Db.Model(&User{}).Count(&count)
		video := Video{}
		video.AuthorId = rand.Int63n(count) + 1000
		video.PlayUrl = gofakeit.URL()
		video.CoverUrl = gofakeit.URL()
		video.PublishTime = gofakeit.Date()
		video.Title = gofakeit.Noun()
		video, err := InsertVideosTable(video)
		if err != nil {
			return
		}
	}
}

func FakeLikes(num int) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < num; i++ {
		var countUser int64
		var countVideo int64
		Db.Model(&User{}).Count(&countUser)
		Db.Model(&Video{}).Count(&countVideo)
		var a, b int64
		a = rand.Int63n(countUser)
		b = rand.Int63n(countVideo)
		like := Like{UserId: a + 1001, VideoId: b + 1001}
		err := InsertLike(&like)
		if err != nil {
			return
		}
	}
}
func FakeComments(num int) {
	gofakeit.Seed(time.Now().Unix())
	rand.Seed(time.Now().Unix())
	for i := 0; i < num; i++ {
		var userCount, videoCount int64
		Db.Model(&User{}).Count(&userCount)
		Db.Model(&Video{}).Count(&videoCount)
		comment := Comment{}
		comment.UserId = rand.Int63n(userCount) + 1000
		comment.VideoId = rand.Int63n(videoCount) + 1000
		comment.CommentText = gofakeit.Sentence(20)
		comment.CreateDate = gofakeit.Date()
		_, err := InsertComment(comment)
		if err != nil {
			return
		}
	}
}
