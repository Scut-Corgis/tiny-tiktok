package dao

import (
	"fmt"
	"testing"
	"time"
)

func TestCountComments(t *testing.T) {
	Init()
	cnt, err := CountComments(1000)
	fmt.Println(cnt)
	fmt.Println(err)
}

func TestQueryCommentsByVideoId(t *testing.T) {
	Init()
	comments, err := QueryCommentsByVideoId(1000)
	fmt.Println(comments)
	fmt.Println(err)
}

func TestInsertComment(t *testing.T) {
	Init()
	comment := Comment{
		UserId:      1000,
		VideoId:     1000,
		CommentText: "你好",
		CreateDate:  time.Now(),
	}
	newComment, err := InsertComment(comment)
	fmt.Println(newComment)
	fmt.Println(err)
}

func TestDeleteComment(t *testing.T) {
	Init()
	flag := DeleteComment(1000)
	fmt.Println(flag)
}
