package dao

import (
	"fmt"
	"testing"
)

func TestInsertLike(t *testing.T) {
	Init()
	err := InsertLike(Like{
		UserId:  1,
		VideoId: 123,
	})
	fmt.Printf("%v", err)
}

func TestDeleteLike(t *testing.T) {
	Init()
	InsertLike(Like{
		UserId:  1,
		VideoId: 123,
	})
	err := DeleteLike(1, 123)
	fmt.Printf("%v", err)
}
