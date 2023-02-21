package dao

import (
	"fmt"
	"testing"
	"time"
)

func TestInsertMessage(t *testing.T) {
	Init()
	id, err := InsertMessage(1000, 1001, "你好", time.Now())
	fmt.Println(id)
	fmt.Println(err)
}

func TestQueryMessagesByMsgKey(t *testing.T) {
	Init()
	messages, err := QueryMessagesByMsgKey(1000, 1001)
	fmt.Println(messages)
	fmt.Println(err)
}

func TestQueryLatestMessageByUserId(t *testing.T) {
	Init()
	message, err := QueryLatestMessageByUserId(1000, 1001)
	fmt.Println(message)
	fmt.Println(err)
}
