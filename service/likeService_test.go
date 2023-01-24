package service

import (
	"fmt"
	"testing"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

func TestLike(t *testing.T) {
	dao.Init()
	err := Like(2, 223)
	fmt.Printf("%v", err)
}

func TestUnlike(t *testing.T) {
	dao.Init()
	Like(2, 223)
	Like(3, 323)

	err := Unlike(2, 223)
	fmt.Printf("%v", err)
}
