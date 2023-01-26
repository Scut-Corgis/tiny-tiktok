package service

import (
	"fmt"
	"testing"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

// func TestLike(t *testing.T) {
// 	dao.Init()
// 	err := Like(2, 223)
// 	fmt.Printf("%v", err)
// }

// func TestUnlike(t *testing.T) {
// 	dao.Init()
// 	Like(2, 223)
// 	Like(3, 323)

// 	err := Unlike(2, 223)
// 	fmt.Printf("%v", err)
// }

// func TestGetVideo(t *testing.T) {
// 	dao.Init()

// 	res := GetVideo(1205, 1002)

// 	fmt.Println(res)
// }

func TestGetLikeLists(t *testing.T) {
	dao.Init()
	res, err := GetLikeLists(1012)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("失败")
	}
	fmt.Println(res)
}
