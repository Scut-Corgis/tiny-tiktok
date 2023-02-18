package dao

import (
	"fmt"
	"testing"
)

func TestInsertLike(t *testing.T) {
	Init()
	for i := 0; i < 1; i++ {
		err := InsertLike(&Like{
			UserId:  int64(1002),
			VideoId: int64(1205 + i),
		})
		fmt.Printf("%v", err)
	}

}

func TestDeleteLike(t *testing.T) {
	Init()
	for i := 4; i < 5; i++ {
		err := DeleteLike(int64(1001), int64(1200+i))
		fmt.Printf("%v", err)
	}
}

func TestGetLikeVideoIdList(t *testing.T) {
	Init()
	list, err := GetLikeVideoIdList(1000)
	fmt.Println(list)
	fmt.Println(err)
}

func TestGetLikeCountByVideoId(t *testing.T) {
	Init()
	cnt, err := GetLikeCountByVideoId(1000)
	fmt.Println(cnt)
	fmt.Println(err)
}

//func TestGetLikInfo(t *testing.T) {
//	Init()
//	likeInfo, err := GetLikInfo(1000, 1000)
//	fmt.Println(likeInfo)
//	fmt.Println(err)
//}
