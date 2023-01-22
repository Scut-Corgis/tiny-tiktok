package ftp

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func TestFfmpeg(t *testing.T) {
	Init()
	for i := 1; i <= 1; i++ {
		name := "wuming.mp4"
		file, err := os.Open("/home/hjg/go/src/tiny-tiktok/data/" + name)
		if err != nil {
			t.Fatal("打开不了文件")
		}
		go SendVideoFile(strconv.Itoa(i)+name, file)
	}
	time.Sleep(2 * time.Second)
}
